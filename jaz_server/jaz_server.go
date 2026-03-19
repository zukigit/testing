package jaz

import (
	"context"
	"fmt"

	"github.com/zukigit/testing/models"
	"github.com/zukigit/testing/zabbix"
)

type JazServer interface {
	GetEnvs() map[string]string
	GetZabbix() zabbix.Zabbix
	GetServerDnsName() string
	GetServerPort() string
	GetServerHost() string
	GetServerMappedPort() string
}

// mandatory envs:
// JAZ_SERVER_VERSION - version of JAZ_SERVER to test (e.g. "1")
// JAZ_DB_TYPE - type of database to use (e.g. "psql")
func NewJaz(ctx context.Context, envs map[string]string, zabbix zabbix.Zabbix) (JazServer, error) {
	switch envs["JAZ_SERVER_VERSION"] {
	case "1":
		switch models.DBType(envs["JAZ_DB_TYPE"]) {
		case models.DBTypePsql:
			return NewJaz1Psql(ctx, envs, zabbix)
		default:
			return nil, fmt.Errorf("unsupported database: %s", envs["JAZ_DB"])
		}
	default:
		return nil, fmt.Errorf("unsupported JAZ version: %s", envs["JAZ_VERSION"])
	}
}
