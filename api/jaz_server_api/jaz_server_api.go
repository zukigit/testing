package jaz_server_api

import (
	"github.com/zukigit/testing/jaz_server"
	"github.com/zukigit/testing/zabbix"
)

type JazServerApi interface {
	GetJazServer() jaz_server.JazServer
	GetZabbix() zabbix.Zabbix
}
