package jaz

import "github.com/zukigit/testing/zabbix"

type Jaz interface {
	GetEnvs() map[string]string
	GetZabbix() zabbix.Zabbix
	GetServerDnsName() string
	GetServerPort() string
	GetServerHost() string
	GetServerMappedPort() string
}
