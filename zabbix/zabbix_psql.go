package zabbix

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
)

type ZabbixPsql struct {
	DBUsername, DBPassword, DBName, DBDnsName, DBPort, DBContainerIp string
	ServerHost, ServerMappedPort, ServerDnsName, ServerPort          string
	NetworkName                                                      string
}

func NewZabbixPsql(ctx context.Context) (Zabbix, error) {
	zabbix := &ZabbixPsql{
		DBUsername: lib.Getenv("ZABBIX_DB_USER", "zabbix"),
		DBPassword: lib.Getenv("ZABBIX_DB_PASSWORD", "zabbix"),
		DBName:     lib.Getenv("ZABBIX_DB_NAME", "zabbix"),
	}

	net, err := network.New(ctx, network.WithDriver("bridge"))
	if err != nil {
		return nil, fmt.Errorf("failed to create network, err: %s", err.Error())
	}
	zabbix.NetworkName = net.Name

	_, err = zabbix.getZabbixDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix DB container, err: %s", err.Error())
	}

	zbxServerContainer, err := zabbix.getZabbixServerContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server container, err: %s", err.Error())
	}

	zabbix.ServerHost, err = zbxServerContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server host, err: %s", err.Error())
	}

	mappedPort, err := zbxServerContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", zabbix.ServerPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server mapped port, err: %s", err.Error())
	}
	zabbix.ServerMappedPort = mappedPort.Port()

	return zabbix, nil
}

func (z *ZabbixPsql) getZabbixDBContainer(ctx context.Context) (testcontainers.Container, error) {
	z.DBPort = "5432"
	z.DBDnsName = "zabbix-postgres"
	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_DB_IMAGE", "postgres:14-alpine3.22"),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", z.DBPort)},
		Networks:     []string{z.NetworkName},
		NetworkAliases: map[string][]string{
			z.NetworkName: {z.DBDnsName},
		},
		Env: map[string]string{
			"POSTGRES_USER":             z.DBUsername,
			"POSTGRES_PASSWORD":         z.DBPassword,
			"POSTGRES_DB":               z.DBName,
			"POSTGRES_HOST_AUTH_METHOD": "trust",
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

func (z *ZabbixPsql) getZabbixServerContainer(ctx context.Context) (testcontainers.Container, error) {
	z.ServerDnsName = "zabbix-server"
	z.ServerPort = "10051"

	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_SERVER_IMAGE", "zabbix/zabbix-server-pgsql:7.0-alpine-latest"),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", z.ServerPort)},
		Networks:     []string{z.NetworkName},
		NetworkAliases: map[string][]string{
			z.NetworkName: {z.ServerDnsName},
		},
		Env: map[string]string{
			"DB_SERVER_HOST":    z.DBDnsName,
			"DB_SERVER_PORT":    z.DBPort,
			"POSTGRES_USER":     z.DBUsername,
			"POSTGRES_PASSWORD": z.DBPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("server #0 started [main process]"),
			wait.ForListeningPort("10051/tcp"),
		).WithDeadline(30 * time.Minute),
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

func (z *ZabbixPsql) GetServerHost() string {
	return z.ServerHost
}

func (z *ZabbixPsql) GetServerPort() string {
	return z.ServerPort
}
