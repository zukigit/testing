package zabbix

import "github.com/zukigit/testing/zabbix"

type ZabbixV1 struct {
	zabbix zabbix.Zabbix
}

func NewZabbixV1(zabbix zabbix.Zabbix) Zabbix {
	return &ZabbixV1{
		zabbix: zabbix,
	}
}
