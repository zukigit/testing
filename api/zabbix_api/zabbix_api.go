package zabbix_api

import "github.com/zukigit/testing/zabbix"

type ZabbixApi interface {
	GetZabbix() zabbix.Zabbix
}

func NewZabbixApi(zabbix zabbix.Zabbix) ZabbixApi {
	return newZabbixApiV1(zabbix)
}
