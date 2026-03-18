package zabbix

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
)

type Zabbix struct {
	dbUsername, dbPassword, dbName, DBHost, MappedPort string
}

func NewZabbix(ctx context.Context) (*Zabbix, error) {
	zabbix := &Zabbix{
		dbUsername: lib.Getenv("ZABBIX_DB_USER", "zabbix"),
		dbPassword: lib.Getenv("ZABBIX_DB_PASSWORD", "zabbix"),
		dbName:     lib.Getenv("ZABBIX_DB_NAME", "zabbix"),
	}

	container, err := zabbix.getZabbixDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix db container, err: %s", err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to mapped port, err: %s", err.Error())
	}
	zabbix.MappedPort = mappedPort.Port()

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host, err: %s", err.Error())
	}
	zabbix.DBHost = host

	return zabbix, nil
}

func (z *Zabbix) getZabbixDBContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_DB_IMAGE", "postgres:14-alpine3.22"),
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     z.dbUsername,
			"POSTGRES_PASSWORD": z.dbPassword,
			"POSTGRES_DB":       z.dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generic container, err: %s", err.Error())
	}

	return container, nil
}
