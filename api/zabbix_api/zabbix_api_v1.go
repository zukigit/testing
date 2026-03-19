package zabbix_api

import "github.com/zukigit/testing/zabbix"

type ZabbixApiV1 struct {
	zabbix zabbix.Zabbix
}

func (z *ZabbixApiV1) GetZabbix() zabbix.Zabbix { return z.zabbix }
