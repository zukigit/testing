package zabbix

import (
	"fmt"

	"github.com/zukigit/testing/models"
)

type ZabbixAPIClient interface {
	// Host operations
	CreateHost(host models.ZbxHost) error
}

func NewZabbixAPIClient(config models.ZabbixAPIConfig) (ZabbixAPIClient, error) {
	switch config.Version {
	case models.ZabbixVersion6:
		return nil, nil // TODO: to be implemented
	case models.ZabbixVersion7:
		return newZabbixClientV7(config)
	default:
		return nil, fmt.Errorf("Unsupported zabbix api version: %s", config.Version)
	}
}
