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
	"github.com/zukigit/testing/models"
)

type ZabbixPsql struct {
	dbUsername, dbPassword, dbName, dbDnsName, dbPort, dbSslMode, dbHost, dbMappedPort string
	serverHost, serverMappedPort, serverDnsName, serverPort                            string
	webHost, webMappedPort, webDnsName, webPort                                        string
	networkName                                                                        string
	envs                                                                               map[string]string
	dbContainer, serverContainer, webContainer                                         testcontainers.Container
}

// getters
func (z *ZabbixPsql) GetDBUsername() string   { return z.dbUsername }
func (z *ZabbixPsql) GetDBPassword() string   { return z.dbPassword }
func (z *ZabbixPsql) GetDBName() string       { return z.dbName }
func (z *ZabbixPsql) GetDBDnsName() string    { return z.dbDnsName }
func (z *ZabbixPsql) GetDBPort() string       { return z.dbPort }
func (z *ZabbixPsql) GetDBSslMode() string    { return z.dbSslMode }
func (z *ZabbixPsql) GetDBHost() string       { return z.dbHost }
func (z *ZabbixPsql) GetDBMappedPort() string { return z.dbMappedPort }

func (z *ZabbixPsql) GetServerHost() string       { return z.serverHost }
func (z *ZabbixPsql) GetServerMappedPort() string { return z.serverMappedPort }
func (z *ZabbixPsql) GetServerDnsName() string    { return z.serverDnsName }
func (z *ZabbixPsql) GetServerPort() string       { return z.serverPort }

func (z *ZabbixPsql) GetWebHost() string       { return z.webHost }
func (z *ZabbixPsql) GetWebMappedPort() string { return z.webMappedPort }
func (z *ZabbixPsql) GetWebDnsName() string    { return z.webDnsName }
func (z *ZabbixPsql) GetWebPort() string       { return z.webPort }

func (z *ZabbixPsql) GetNetworkName() string { return z.networkName }

func (z *ZabbixPsql) GetDBContainer() testcontainers.Container     { return z.dbContainer }
func (z *ZabbixPsql) GetServerContainer() testcontainers.Container { return z.serverContainer }
func (z *ZabbixPsql) GetWebContainer() testcontainers.Container    { return z.webContainer }

func (z *ZabbixPsql) getZabbixDBContainer(ctx context.Context) (testcontainers.Container, error) {
	z.dbPort = lib.GetEnv(z.envs, "ZABBIX_DB_PORT", "5432")
	z.dbDnsName = lib.GetEnv(z.envs, "ZABBIX_DB_DNS_NAME", "zabbix-db")
	z.dbSslMode = lib.GetEnv(z.envs, "ZABBIX_DB_SSL_MODE", string(models.SslModeDisable))
	portUrl := fmt.Sprintf("%s/tcp", z.dbPort)

	req := testcontainers.ContainerRequest{
		Image:        lib.GetEnv(z.envs, "ZABBIX_DB_IMAGE", "postgres:14-alpine3.22"),
		ExposedPorts: []string{portUrl},
		Networks:     []string{z.networkName},
		NetworkAliases: map[string][]string{
			z.networkName: {z.dbDnsName},
		},
		Env: map[string]string{
			"POSTGRES_USER":             z.dbUsername,
			"POSTGRES_PASSWORD":         z.dbPassword,
			"POSTGRES_DB":               z.dbName,
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

	z.dbHost, err = container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server host, err: %s", err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(portUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server mapped port, err: %s", err.Error())
	}
	z.dbMappedPort = mappedPort.Port()

	return container, nil
}

func (z *ZabbixPsql) getZabbixServerContainer(ctx context.Context) (testcontainers.Container, error) {
	z.serverDnsName = lib.GetEnv(z.envs, "ZABBIX_SERVER_DNS_NAME", "zabbix-server")
	z.serverPort = lib.GetEnv(z.envs, "ZABBIX_SERVER_PORT", "10051")

	req := testcontainers.ContainerRequest{
		Image:        lib.GetEnv(z.envs, "ZABBIX_SERVER_IMAGE", "zabbix/zabbix-server-pgsql:7.0-alpine-latest"),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", z.serverPort)},
		Networks:     []string{z.networkName},
		NetworkAliases: map[string][]string{
			z.networkName: {z.serverDnsName},
		},
		Env: map[string]string{
			"DB_SERVER_HOST":    z.dbDnsName,
			"DB_SERVER_PORT":    z.dbPort,
			"POSTGRES_USER":     z.dbUsername,
			"POSTGRES_PASSWORD": z.dbPassword,
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

	z.serverHost, err = container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server host, err: %s", err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", z.serverPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server mapped port, err: %s", err.Error())
	}
	z.serverMappedPort = mappedPort.Port()

	return container, nil
}

func (z *ZabbixPsql) getZabbixWebContainer(ctx context.Context) (testcontainers.Container, error) {
	z.webDnsName = lib.GetEnv(z.envs, "ZABBIX_WEB_DNS_NAME", "zabbix-web")
	z.webPort = lib.GetEnv(z.envs, "ZABBIX_WEB_PORT", "8080")
	portWithTcp := fmt.Sprintf("%s/tcp", z.webPort)

	req := testcontainers.ContainerRequest{
		Image:        lib.GetEnv(z.envs, "ZABBIX_WEB_IMAGE", "zabbix/zabbix-web-nginx-pgsql:7.0-alpine-latest"),
		ExposedPorts: []string{portWithTcp},
		Networks:     []string{z.networkName},
		NetworkAliases: map[string][]string{
			z.networkName: {z.webDnsName},
		},
		Env: map[string]string{
			"DB_SERVER_HOST":    z.dbDnsName,
			"DB_SERVER_PORT":    z.dbPort,
			"POSTGRES_USER":     z.dbUsername,
			"POSTGRES_PASSWORD": z.dbPassword,
			"ZBX_SERVER_HOST":   z.serverDnsName,
			"PHP_TZ":            lib.GetEnv(z.envs, "ZABBIX_WEB_PHP_TZ", "Asia/Yangon"),
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(portWithTcp)).WithStartupTimeout(2*time.Minute),
			wait.ForHTTP("/").
				WithPort(nat.Port(portWithTcp)).
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

	z.webHost, err = container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web host, err: %s", err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", z.webPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web mapped port, err: %s", err.Error())
	}
	z.webMappedPort = mappedPort.Port()

	return container, nil
}

// zabbix represents active running container
func newZabbixPsql(ctx context.Context, envs map[string]string) (Zabbix, error) {
	zabbix := &ZabbixPsql{
		dbUsername: lib.GetEnv(envs, "ZABBIX_DB_USER", "zabbix"),
		dbPassword: lib.GetEnv(envs, "ZABBIX_DB_PASSWORD", "zabbix"),
		dbName:     lib.GetEnv(envs, "ZABBIX_DB_NAME", "zabbix"),
		envs:       envs,
	}

	net, err := network.New(ctx, network.WithDriver("bridge"))
	if err != nil {
		return nil, fmt.Errorf("failed to create network, err: %s", err.Error())
	}
	zabbix.networkName = net.Name

	// zabbix db
	zabbix.dbContainer, err = zabbix.getZabbixDBContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix DB container, err: %s", err.Error())
	}

	// zabbix server
	zabbix.serverContainer, err = zabbix.getZabbixServerContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix server container, err: %s", err.Error())
	}

	// zabbix web
	zabbix.webContainer, err = zabbix.getZabbixWebContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zabbix web container, err: %s", err.Error())
	}

	return zabbix, nil
}
