package jaz_server

import (
	"github.com/zukigit/testing/zabbix"
)

type JazServerApiClientV1 struct {
	jazServer JazServer
	zabbix    zabbix.Zabbix
}

func (j *JazServerApiClientV1) GetJazServer() JazServer { return j.jazServer }
func (j *JazServerApiClientV1) GetZabbix() zabbix.Zabbix           { return j.zabbix }

func newJazServerApiClientV1(zabbix zabbix.Zabbix, jazServer JazServer) JazServerApiClient {
	return &JazServerApiClientV1{
		jazServer: jazServer,
		zabbix:    zabbix,
	}
}