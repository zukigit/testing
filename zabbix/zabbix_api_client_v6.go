package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zukigit/testing/models"
)

type ZabbixAPIClientV6 struct {
	URL        string
	Username   string
	Password   string
	Auth       string
	HTTPClient *http.Client
}

func newZabbixClientV6(config models.ZabbixAPIConfig) (*ZabbixAPIClientV6, error) {
	username := config.Username
	if username == "" {
		username = "Admin"
	}

	password := config.Password
	if password == "" {
		password = "zabbix"
	}

	client := &ZabbixAPIClientV6{
		URL:      fmt.Sprintf("http://%s:%s/api_jsonrpc.php", config.WebHost, config.WebPort),
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: models.ZabbixAPITimeout,
		},
	}

	if err := client.login(); err != nil {
		return nil, fmt.Errorf("failed to login to zabbix api: %w", err)
	}

	return client, nil
}

// ---------- Generic API Call ----------
func (c *ZabbixAPIClientV6) call(method models.ZabbixAPIMethod, params any, result any) error {
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
func (c *ZabbixAPIClientV6) login() error {
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

func (c *ZabbixAPIClientV6) ensureHostGroup(name string) (string, error) {
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

func (c *ZabbixAPIClientV6) getTemplateByName(name string) (string, error) {
	var resp []struct {
		TemplateID string `json:"templateid"`
		Name       string `json:"name"`
	}

	err := c.call(models.TemplateGet, map[string]any{
		"filter": map[string]any{
			"host": []string{name},
		},
	}, &resp)
	if err != nil {
		return "", err
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("template '%s' not found", name)
	}

	return resp[0].TemplateID, nil
}

func (c *ZabbixAPIClientV6) getHostID(name string) (string, error) {
	var resp []struct {
		HostID string `json:"hostid"`
	}

	err := c.call(models.HostGet, map[string]any{
		"filter": map[string]any{
			"host": []string{name},
		},
	}, &resp)
	if err != nil {
		return "", err
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("host '%s' not found", name)
	}

	return resp[0].HostID, nil
}

// ---------- Exposed Functions ----------
func (c *ZabbixAPIClientV6) CreateHost(host models.ZbxHost) error {

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

	if len(host.Templates) > 0 {
		if err := c.AttachHostTemplates(host.Hostname, host.Templates); err != nil {
			return fmt.Errorf("failed to attach templates: %v", err)
		}
	}

	if len(host.HostMacros) > 0 {
		if err := c.AddHostMacros(host.Hostname, host.HostMacros); err != nil {
			return fmt.Errorf("failed to add host macros: %v", err)
		}
	}

	return nil
}

func (c *ZabbixAPIClientV6) ImportTemplate(tmplFilePath string) error {
	sourceData, err := os.ReadFile(tmplFilePath)
	if err != nil {
		return fmt.Errorf("failed to read import file '%s': %v", tmplFilePath, err)
	}

	ext := filepath.Ext(tmplFilePath)

	// Remove the dot
	if len(ext) > 0 {
		ext = ext[1:]
	}

	params := map[string]any{
		"format": ext,
		"source": string(sourceData),
		"rules": map[string]any{
			"templates": map[string]any{
				"createMissing":  true,
				"updateExisting": true,
			},
			"items": map[string]any{
				"createMissing":  true,
				"updateExisting": true,
			},
			"triggers": map[string]any{
				"createMissing":  true,
				"updateExisting": true,
			},
			"graphs": map[string]any{
				"createMissing":  true,
				"updateExisting": true,
			},
		},
	}

	err = c.call(models.ConfigurationImport, params, nil)
	if err != nil {
		return fmt.Errorf("failed to import template: %w", err)
	}

	return nil
}

func (c *ZabbixAPIClientV6) AddHostMacros(hostName string, newMacros []models.HostMacro) error {
	if len(newMacros) == 0 {
		return fmt.Errorf("no host macro is provided")
	}

	// 1. Get host ID
	hostID, err := c.getHostID(hostName)
	if err != nil {
		return err
	}

	// 2. Get current macros
	var hostInfo []map[string]any
	getParams := map[string]any{
		"hostids":      hostID,
		"selectMacros": "extend",
	}
	if err := c.call(models.HostGet, getParams, &hostInfo); err != nil {
		return fmt.Errorf("failed to get current macros for host '%s': %w", hostName, err)
	}

	// 3. Merge new macros (avoid duplicates)
	if len(hostInfo) <= 0 {
		return fmt.Errorf("failed to get host info for '%s': %w", hostName, err)
	}

	var currentMacros []models.HostMacro
	macros := hostInfo[0]["macros"]
	// Convert interface{} → JSON → struct slice
	b, err := json.Marshal(macros)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &currentMacros); err != nil {
		return err
	}

	macroMap := make(map[string]models.HostMacro)

	// Add current macros
	for _, m := range currentMacros {
		if m.Macro != "" { // ignore empty macros
			macroMap[m.Macro] = m
		}
	}

	// Add new macros (overwrite if same macro name)
	for _, m := range newMacros {
		if m.Macro != "" {
			macroMap[m.Macro] = m
		}
	}

	// Convert back to slice
	mergedMacros := make([]models.HostMacro, 0, len(macroMap))
	for _, m := range macroMap {
		mergedMacros = append(mergedMacros, m)
	}

	// 4. Update host with combined macros
	params := map[string]any{
		"hostid": hostID,
		"macros": mergedMacros,
	}
	if err := c.call(models.HostUpdate, params, nil); err != nil {
		return fmt.Errorf("failed to add macros to host '%s': %w", hostName, err)
	}

	return nil
}

func (c *ZabbixAPIClientV6) AttachHostTemplates(hostName string, templateNames []string) error {

	if len(templateNames) == 0 {
		return nil
	}

	// 1. Get host ID + existing templates
	var resp []struct {
		HostID          string `json:"hostid"`
		ParentTemplates []struct {
			TemplateID string `json:"templateid"`
			Name       string `json:"name"`
		} `json:"parentTemplates"`
	}

	err := c.call(models.HostGet, map[string]any{
		"filter": map[string]any{
			"host": []string{hostName},
		},
		"selectParentTemplates": []string{"templateid", "name"},
	}, &resp)
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return fmt.Errorf("host '%s' not found", hostName)
	}

	hostID := resp[0].HostID

	// 2. Build existing template map
	templateMap := make(map[string]string) // name -> id
	for _, t := range resp[0].ParentTemplates {
		templateMap[t.Name] = t.TemplateID
	}

	// 3. Resolve new templates
	for _, name := range templateNames {
		if _, exists := templateMap[name]; exists {
			continue
		}

		tid, err := c.getTemplateByName(name)
		if err != nil {
			return err
		}

		templateMap[name] = tid
	}

	// 4. Convert to payload
	var templates []map[string]string
	for _, tid := range templateMap {
		templates = append(templates, map[string]string{
			"templateid": tid,
		})
	}

	// 5. Update host
	params := map[string]any{
		"hostid":    hostID,
		"templates": templates,
	}

	if err := c.call(models.HostUpdate, params, nil); err != nil {
		return fmt.Errorf("failed to attach templates to host '%s': %w", hostName, err)
	}

	return nil
}
