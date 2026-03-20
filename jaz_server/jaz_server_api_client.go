package jaz_server

import (
	"fmt"

	"github.com/zukigit/testing/zabbix"
)

type JazServerApiClient interface {
	GetJazServer() JazServer
	GetZabbix() zabbix.Zabbix
}

func NewJazServerApi(zabbix zabbix.Zabbix, jazServer JazServer) (JazServerApiClient, error) {
	switch v := jazServer.(type) {
	case *jazServerPsqlV1:
		return newJazServerApiClientV1(zabbix, jazServer), nil
	default:
		return nil, fmt.Errorf("unsupported JazServer type: %T", v)
	}
}
