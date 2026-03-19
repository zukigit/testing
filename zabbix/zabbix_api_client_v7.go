package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zukigit/testing/models"
)

type ZabbixAPIClientV7 struct {
	URL        string
	Username   string
	Password   string
	Auth       string
	HTTPClient *http.Client
}

func newZabbixClientV7(config models.ZabbixAPIConfig) (*ZabbixAPIClientV7, error) {
	username := config.Username
	if username == "" {
		username = "Admin"
	}

	password := config.Password
	if password == "" {
		password = "zabbix"
	}

	client := &ZabbixAPIClientV7{
		URL:      fmt.Sprintf("http://%s:%s/api_jsonrpc.php", config.WebHost, config.WebPort),
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	if err := client.login(); err != nil {
		return nil, fmt.Errorf("failed to login to zabbix api: %w", err)
	}

	return client, nil
}

// ---------- Generic API Call ----------
func (c *ZabbixAPIClientV7) call(method models.ZabbixAPIMethod, params any, result any) error {
	// Ensure auth
	if c.Auth == "" && method != models.UserLogin {
		if err := c.login(); err != nil {
			return err
		}
	}

	reqBody := models.ApiRequest{
		Jsonrpc: "2.0",
		Method:  string(method),
		Params:  params,
		ID:      1,
		Auth:    c.Auth,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp models.ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	if apiResp.Error != nil {
		return fmt.Errorf("zabbix API error [method=%s]: %v", string(method), apiResp.Error)
	}

	if result != nil {
		return json.Unmarshal(apiResp.Result, result)
	}

	return nil
}

// ---------- login ----------
func (c *ZabbixAPIClientV7) login() error {
	var result string

	err := c.call(models.UserLogin, map[string]any{
		"username": c.Username,
		"password": c.Password,
	}, &result)

	if err != nil {
		return fmt.Errorf("failed to login with username (%s) and password (%s): %w", c.Username, c.Password, err)
	}

	c.Auth = result
	return nil
}

func (c *ZabbixAPIClientV7) ensureHostGroup(name string) (string, error) {
	// 1. Try to get existing group
	var resp []struct {
		GroupID string `json:"groupid"`
		Name    string `json:"name"`
	}

	err := c.call(models.HostGroupGet, map[string]any{
		"filter": map[string]any{
			"name": []string{name},
		},
	}, &resp)
	if err != nil {
		return "", err
	}

	if len(resp) > 0 {
		return resp[0].GroupID, nil
	}

	// 2. Create if not exists
	var createResp struct {
		GroupIDs []string `json:"groupids"`
	}

	err = c.call(models.HostGroupCreate, map[string]any{
		"name": name,
	}, &createResp)
	if err != nil {
		return "", err
	}

	if len(createResp.GroupIDs) == 0 {
		return "", fmt.Errorf("hostgroup created but no groupid returned")
	}

	return createResp.GroupIDs[0], nil
}

// ---------- Main Function ----------
func (c *ZabbixAPIClientV7) CreateHost(host models.ZbxHost) error {

	// 1. Ensure hostgroup
	groupID, err := c.ensureHostGroup("Job Arranger")
	if err != nil {
		return err
	}

	// 2. Check if host exists
	var hostResp []map[string]string
	err = c.call(models.HostGet, map[string]any{
		"filter": map[string]any{
			"host": []string{host.Hostname},
		},
	}, &hostResp)
	if err != nil {
		return fmt.Errorf("failed to check host '%s' existence", host.Hostname)
	}

	// 3. If exists → delete
	if len(hostResp) > 0 {
		hostID := hostResp[0]["hostid"]

		err = c.call(models.HostDelete, []string{hostID}, nil)
		if err != nil {
			return fmt.Errorf("failed to delete host '%s'", host.Hostname)
		}

		time.Sleep(500 * time.Millisecond)
	}

	// 4. Create host with agent interface
	useIPInt := 0
	if host.UseIP {
		useIPInt = 1
	}

	intf := models.Interface{
		Type:  1,
		Main:  1,
		UseIP: useIPInt,
		IP:    host.IP,
		DNS:   host.DNS,
		Port:  "10050",
	}

	params := models.CreateHostParams{
		Host: host.Hostname,
		Groups: []models.Group{
			{GroupID: groupID},
		},
		Interfaces: []models.Interface{
			intf,
		},
	}

	var createResp map[string][]string
	err = c.call(models.HostCreate, params, &createResp)
	if err != nil {
		return fmt.Errorf("failed to create host '%s'", host.Hostname)
	}

	return nil
}
