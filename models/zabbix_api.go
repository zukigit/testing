package models

import "encoding/json"

const (
	ZabbixVersion6 = "6"
	ZabbixVersion7 = "7"
)

type ZabbixAPIMethod string

const (
	UserLogin       ZabbixAPIMethod = "user.login"
	HostCreate      ZabbixAPIMethod = "host.create"
	HostGet         ZabbixAPIMethod = "host.get"
	HostDelete      ZabbixAPIMethod = "host.delete"
	HostGroupGet    ZabbixAPIMethod = "hostgroup.get"
	HostGroupCreate ZabbixAPIMethod = "hostgroup.create"
)

type ApiRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
	Auth    string `json:"auth,omitempty"`
}

type ApiResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   any             `json:"error"`
	ID      int             `json:"id"`
}

type ZabbixAPIConfig struct {
	Version  string
	WebHost  string // IP or DNS
	WebPort  string
	Username string // Default: Admin
	Password string // Default: zabbix
}

type CreateHostParams struct {
	Host       string      `json:"host"`
	Groups     []Group     `json:"groups"`
	Interfaces []Interface `json:"interfaces"`
}

type Group struct {
	GroupID string `json:"groupid"`
}

type Interface struct {
	Type  int    `json:"type"`  // 1 = agent
	Main  int    `json:"main"`  // 1 = default interface
	UseIP int    `json:"useip"` // 1 = use IP, 0 = use DNS
	IP    string `json:"ip"`
	DNS   string `json:"dns"`
	Port  string `json:"port"`
}

type ZbxHost struct {
	Hostname string
	IP       string
	DNS      string
	UseIP    bool
}
