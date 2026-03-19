package jaz_server_api

import (
	"github.com/zukigit/testing/jaz_server"
	"github.com/zukigit/testing/zabbix"
)

type JazServerApiV1 struct {
	jazServer jaz_server.JazServer
	zabbix    zabbix.Zabbix
}

func (j *JazServerApiV1) GetJazServer() jaz_server.JazServer { return j.jazServer }
func (j *JazServerApiV1) GetZabbix() zabbix.Zabbix           { return j.zabbix }
