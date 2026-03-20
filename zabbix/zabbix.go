package zabbix

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/zukigit/testing/lib"
	"github.com/zukigit/testing/models"
)

type Zabbix interface {
	// DB
	GetDBUsername() string
	GetDBPassword() string
	GetDBName() string
	GetDBDnsName() string // only for docker internal communication
	GetDBPort() string    // only for docker internal communication
	GetDBContainer() testcontainers.Container
	GetDBSslMode() string
	GetDBHost() string
	GetDBMappedPort() string

	// server
	GetServerHost() string       // it can be used in testcases
	GetServerMappedPort() string // it can be used in testcases
	GetServerDnsName() string    // only for docker internal communication
	GetServerPort() string       // only for docker internal communication
	GetServerContainer() testcontainers.Container

	// web
	GetWebHost() string       // it can be used in testcases
	GetWebMappedPort() string // it can be used in testcases
	GetWebDnsName() string    // only for docker internal communication
	GetWebPort() string       // only for docker internal communication
	GetWebContainer() testcontainers.Container

	// network
	GetNetworkName() string
}

func NewZabbix(ctx context.Context, envs map[string]string) (Zabbix, error) {
	dbType := lib.GetEnv(envs, "ZABBIX_DB_TYPE", "")

	switch models.DBType(dbType) {
	case models.DBTypeMysql:
		return nil, nil
	case models.DBTypePsql:
		return newZabbixPsql(ctx, envs)
	default:
		return nil, fmt.Errorf("unknown db type: %s", dbType)
	}
}
