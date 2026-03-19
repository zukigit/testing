package zabbix

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
)

type ZabbixPsql struct {
	DBUsername, DBPassword, DBName, DBDnsName, DBPort       string
	ServerHost, ServerMappedPort, ServerDnsName, ServerPort string
	WebHost, WebMappedPort, WebDnsName, WebPort             string
	NetworkName                                             string
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

func (z *ZabbixPsql) getZabbixWebContainer(ctx context.Context) (testcontainers.Container, error) {
	z.WebDnsName = "zabbix-web"
	z.WebPort = "8080"

	req := testcontainers.ContainerRequest{
		Image:        lib.Getenv("ZABBIX_WEB_IMAGE", "zabbix/zabbix-web-nginx-pgsql:7.0-alpine-latest"),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", z.WebPort)},
		Networks:     []string{z.NetworkName},
		NetworkAliases: map[string][]string{
			z.NetworkName: {z.WebDnsName},
		},
		Env: map[string]string{
			"DB_SERVER_HOST":    z.DBDnsName,
			"DB_SERVER_PORT":    z.DBPort,
			"POSTGRES_USER":     z.DBUsername,
			"POSTGRES_PASSWORD": z.DBPassword,
			"ZBX_SERVER_HOST":   z.ServerDnsName,
			"PHP_TZ":            "Asia/Yangon",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("8080/tcp").WithStartupTimeout(2*time.Minute),
			wait.ForHTTP("/").
				WithPort("8080/tcp").
				WithStatusCodeMatcher(func(status int) bool {
					return status == http.StatusOK || status == http.StatusFound
				}).
				WithStartupTimeout(5*time.Minute),
		).WithDeadline(10 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start zabbix-web container: %s", err.Error())
	}

	return container, nil
}

// getters
func (z *ZabbixPsql) GetDBUsername() string { return z.DBUsername }
func (z *ZabbixPsql) GetDBPassword() string { return z.DBPassword }
func (z *ZabbixPsql) GetDBName() string     { return z.DBName }
func (z *ZabbixPsql) GetDBDnsName() string  { return z.DBDnsName }
func (z *ZabbixPsql) GetDBPort() string     { return z.DBPort }

func (z *ZabbixPsql) GetServerHost() string       { return z.ServerHost }
func (z *ZabbixPsql) GetServerMappedPort() string { return z.ServerMappedPort }
func (z *ZabbixPsql) GetServerDnsName() string    { return z.ServerDnsName }
func (z *ZabbixPsql) GetServerPort() string       { return z.ServerPort }

func (z *ZabbixPsql) GetWebHost() string       { return z.WebHost }
func (z *ZabbixPsql) GetWebMappedPort() string { return z.WebMappedPort }
func (z *ZabbixPsql) GetWebDnsName() string    { return z.WebDnsName }
func (z *ZabbixPsql) GetWebPort() string       { return z.WebPort }

func (z *ZabbixPsql) GetNetworkName() string { return z.NetworkName }

// zabbix represents active running container
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

	// zabbix db
	_, err = zabbix.getZabbixDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix DB container, err: %s", err.Error())
	}

	zbxServerContainer, err := zabbix.getZabbixServerContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server container, err: %s", err.Error())
	}

	// zabbix server
	zabbix.ServerHost, err = zbxServerContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server host, err: %s", err.Error())
	}

	mappedPort, err := zbxServerContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", zabbix.ServerPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server mapped port, err: %s", err.Error())
	}
	zabbix.ServerMappedPort = mappedPort.Port()

	// zabbix web
	zbxWebContainer, err := zabbix.getZabbixWebContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web container, err: %s", err.Error())
	}

	zabbix.WebHost, err = zbxWebContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web host, err: %s", err.Error())
	}

	mappedPort, err = zbxWebContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", zabbix.WebPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web mapped port, err: %s", err.Error())
	}
	zabbix.WebMappedPort = mappedPort.Port()

	return zabbix, nil
}
