package zabbix

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
)

type Zabbix struct {
	DBUsername, DBPassword, DBName, DBHost, MappedPort string
	ServerHost, ServerPort                             string
}

func NewZabbix(ctx context.Context) (*Zabbix, error) {
	zabbix := &Zabbix{
		DBUsername: lib.Getenv("ZABBIX_DB_USER", "zabbix"),
		DBPassword: lib.Getenv("ZABBIX_DB_PASSWORD", "zabbix"),
		DBName:     lib.Getenv("ZABBIX_DB_NAME", "zabbix"),
	}

	dbContainer, err := zabbix.getZabbixDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix DB container, err: %s", err.Error())
	}

	mappedPort, err := dbContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to mapped port, err: %s", err.Error())
	}
	zabbix.MappedPort = mappedPort.Port()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host, err: %s", err.Error())
	}
	zabbix.DBHost = host

	serverContainer, err := zabbix.getZabbixServerContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server container, err: %s", err.Error())
	}

	mappedPort, err = serverContainer.MappedPort(ctx, "10051/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to mapped port, err: %s", err.Error())
	}
	zabbix.ServerPort = mappedPort.Port()

	host, err = serverContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host, err: %s", err.Error())
	}
	zabbix.ServerHost = host

	return zabbix, nil
}

func (z *Zabbix) getZabbixDBContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_DB_IMAGE", "postgres:14-alpine3.22"),
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     z.DBUsername,
			"POSTGRES_PASSWORD": z.DBPassword,
			"POSTGRES_DB":       z.DBName,
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

func (z *Zabbix) getZabbixServerContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_SERVER_IMAGE", "zabbix/zabbix-server-pgsql:7.0-alpine-latest"),
		ExposedPorts: []string{"10051/tcp"},
		Env: map[string]string{
			"DB_SERVER_HOST":   z.DBHost,
			"DB_SERVER_PORT":   z.MappedPort,
			"DB_SERVER_DBNAME": z.DBName,
			"DB_SERVER_USER":   z.DBUsername,
			"DB_SERVER_PASSWD": z.DBPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("server #0 started [main process]"),
			wait.ForListeningPort("10051/tcp"),
		).WithDeadline(2 * time.Minute),
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
