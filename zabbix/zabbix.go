package zabbix

import (
	"context"
	"fmt"

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

	// server
	GetServerHost() string       // it can be used in testcases
	GetServerMappedPort() string // it can be used in testcases
	GetServerDnsName() string    // only for docker internal communication
	GetServerPort() string       // only for docker internal communication

	// network
	GetNetworkName() string
}

func NewZabbix(ctx context.Context) (Zabbix, error) {
	dbType := lib.Getenv("ZABBIX_DB_TYPE", "")

	switch models.DBType(dbType) {
	case models.DBTypeMysql:
		return nil, nil
	case models.DBTypePsql:
		return NewZabbixPsql(ctx)
	default:
		return nil, fmt.Errorf("unknown db type: %s", dbType)
	}
}
