package jaz_server

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/zukigit/testing/lib"
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

	GetServerContainer() testcontainers.Container
}

// NewJazServer creates and returns a new JazServer instance.
//
// Required environment variables (via envs):
//   - JAZ_SERVER_VERSION: version of the JAZ server to test (e.g. "1")
//   - JAZ_DB_TYPE:        database type to use (e.g. "psql")
//   - JAZ_SERVER_IMAGE:   Docker image for the server (e.g. "jobarg-server-postgres:6.0.9-1")
func NewJazServer(ctx context.Context, envs map[string]string, zabbix zabbix.Zabbix) (JazServer, error) {
	if zabbix == nil {
		return nil, fmt.Errorf("zabbix is nil")
	}

	err := lib.CheckEmptyValues(envs, []string{
		"JAZ_SERVER_VERSION",
		"JAZ_DB_TYPE",
		"JAZ_SERVER_IMAGE",
	})
	if err != nil {
		return nil, err
	}

	switch envs["JAZ_SERVER_VERSION"] {
	case "1":
		switch models.DBType(envs["JAZ_DB_TYPE"]) {
		case models.DBTypePsql:
			return newJazServerPsqlV1(ctx, envs, zabbix)
		default:
			return nil, fmt.Errorf("unsupported database: %s", envs["JAZ_DB"])
		}
	default:
		return nil, fmt.Errorf("unsupported JAZ version: %s", envs["JAZ_VERSION"])
	}
}
