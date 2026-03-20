package jaz_server_api

import (
	"fmt"

	"github.com/zukigit/testing/jaz_server"
	"github.com/zukigit/testing/zabbix"
)

type JazServerApi interface {
	GetJazServer() jaz_server.JazServer
	GetZabbix() zabbix.Zabbix
}

func NewJazServerApi(zabbix zabbix.Zabbix, jazServer jaz_server.JazServer) (JazServerApi, error) {
	switch v := jazServer.(type) {
	case *jaz_server.JazServerPsqlV1:
		return newJazServerApiV1(zabbix, jazServer), nil
	default:
		return nil, fmt.Errorf("unsupported JazServer type: %T", v)
	}
}
