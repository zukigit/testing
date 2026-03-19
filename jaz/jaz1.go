package jaz

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
	"github.com/zukigit/testing/zabbix"
)

type Jaz1 struct {
	envs   map[string]string
	zabbix zabbix.Zabbix

	serverDnsName, serverPort, serverHost, serverMappedPort string
}

func (j *Jaz1) GetEnvs() map[string]string {
	return j.envs
}

func (j *Jaz1) GetZabbix() zabbix.Zabbix {
	return j.zabbix
}

func (j *Jaz1) GetServerDnsName() string {
	return j.serverDnsName
}

func (j *Jaz1) GetServerPort() string {
	return j.serverPort
}

func (j *Jaz1) GetServerHost() string {
	return j.serverHost
}

func (j *Jaz1) GetServerMappedPort() string {
	return j.serverMappedPort
}

func NewJaz1(ctx context.Context, envs map[string]string, zabbix zabbix.Zabbix) (Jaz, error) {
	jaz1 := &Jaz1{
		envs:   envs,
		zabbix: zabbix,
	}

	serverContainer, err := jaz1.newServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create jaz server container, err: %s", err.Error())
	}

	jaz1.serverHost, err = serverContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get jaz server host, err: %s", err.Error())
	}

	mappedPort, err := serverContainer.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", jaz1.serverPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get jaz server mapped port, err: %s", err.Error())
	}
	jaz1.serverMappedPort = mappedPort.Port()

	return jaz1, nil
}

func (j *Jaz1) newServer(ctx context.Context) (testcontainers.Container, error) {
	j.serverDnsName = lib.GetEnv(j.envs, "JAZ_SERVER_DNS_NAME", "jaz-server")
	j.serverPort = lib.GetEnv(j.envs, "JAZ_SERVER_PORT", "10061")
	portWithTcp := fmt.Sprintf("%s/tcp", j.serverPort)

	zabbixUrl := fmt.Sprintf("http://%s:%s", j.zabbix.GetServerDnsName(), j.zabbix.GetServerPort())

	j.envs["DB_SERVER_HOST"] = lib.GetEnv(j.envs, j.envs["DB_SERVER_HOST"], j.zabbix.GetDBDnsName())
	j.envs["POSTGRES_DATABASE"] = lib.GetEnv(j.envs, j.envs["POSTGRES_DATABASE"], j.zabbix.GetDBName())
	j.envs["POSTGRES_USER"] = lib.GetEnv(j.envs, j.envs["POSTGRES_USER"], j.zabbix.GetDBUsername())
	j.envs["POSTGRES_PASSWORD"] = lib.GetEnv(j.envs, j.envs["POSTGRES_PASSWORD"], j.zabbix.GetDBPassword())
	j.envs["JAZABBIXURL"] = lib.GetEnv(j.envs, j.envs["JAZABBIXURL"], zabbixUrl)
	j.envs["LOGTYPE"] = lib.GetEnv(j.envs, j.envs["LOGTYPE"], "file")
	j.envs["DEBUGLEVEL"] = lib.GetEnv(j.envs, j.envs["DEBUGLEVEL"], "3")

	req := testcontainers.ContainerRequest{
		Image:        lib.GetEnv(j.envs, "JAZ_SERVER_IMAGE", "jobarg-server-postgres:6.0.9-1"),
		ExposedPorts: []string{portWithTcp},
		Networks:     []string{j.zabbix.GetNetworkName()},
		NetworkAliases: map[string][]string{
			j.zabbix.GetNetworkName(): {j.serverDnsName},
		},
		Env:        j.envs,
		WaitingFor: wait.ForListeningPort(nat.Port(portWithTcp)).WithStartupTimeout(5 * time.Minute),
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
