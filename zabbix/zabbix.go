package zabbix

import (
	"context"
	"fmt"

	"github.com/zukigit/testing/lib"
	"github.com/zukigit/testing/models"
)

type Zabbix interface {
	GetServerHost() string
	GetServerPort() string
}

func NewZabbix(ctx context.Context) (Zabbix, error) {
	dbType := lib.Getenv("ZABBIX_DB_TYPE", "")

	switch models.DBType(dbType) {
	case models.DBTypeMysql:
		return NewZabbixMysql(ctx)
	case models.DBTypePsql:
		return NewZabbixPsql(ctx)
	default:
		return nil, fmt.Errorf("unknown db type: %s", dbType)
	}
}
