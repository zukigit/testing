package jaz_server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zukigit/testing/lib"
	"github.com/zukigit/testing/zabbix"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
)

type jaz1Psql struct {
	envs   map[string]string
	zabbix zabbix.Zabbix

	serverDnsName, serverPort, serverHost, serverMappedPort string

	db *sql.DB

	serverContainer testcontainers.Container
}

func (j *jaz1Psql) GetServerContainer() testcontainers.Container {
	return j.serverContainer
}

func (j *jaz1Psql) GetDB() *sql.DB {
	return j.db
}

func (j *jaz1Psql) GetEnvs() map[string]string {
	return j.envs
}

func (j *jaz1Psql) GetZabbix() zabbix.Zabbix {
	return j.zabbix
}

func (j *jaz1Psql) GetServerDnsName() string {
	return j.serverDnsName
}

func (j *jaz1Psql) GetServerPort() string {
	return j.serverPort
}

func (j *jaz1Psql) GetServerHost() string {
	return j.serverHost
}

func (j *jaz1Psql) GetServerMappedPort() string {
	return j.serverMappedPort
}

func (j *jaz1Psql) newServer(ctx context.Context) (testcontainers.Container, error) {
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
		Image:        j.envs["JAZ_SERVER_IMAGE"],
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

	j.serverHost, err = container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get jaz server host, err: %s", err.Error())
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", j.serverPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get jaz server mapped port, err: %s", err.Error())
	}
	j.serverMappedPort = mappedPort.Port()

	return container, nil
}

func (j *jaz1Psql) connectDB() error {
	// we will use zabbix db because jaz1 use zabbix database
	host := j.zabbix.GetDBHost()
	port := j.zabbix.GetDBMappedPort()
	user := j.zabbix.GetDBUsername()
	password := j.zabbix.GetDBPassword()
	dbname := j.zabbix.GetDBName()
	sslMode := j.zabbix.GetDBSslMode()

	// Build connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)

	// Open a database connection
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database connection, err: %s", err.Error())
	}
	j.db = db

	return db.Ping()
}

func newJaz1Psql(ctx context.Context, envs map[string]string, zabbix zabbix.Zabbix) (JazServer, error) {
	jaz1Psql := &jaz1Psql{
		envs:   envs,
		zabbix: zabbix,
	}

	// server
	container, err := jaz1Psql.newServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create jaz server container, err: %s", err.Error())
	}
	jaz1Psql.serverContainer = container

	// db
	err = jaz1Psql.connectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect jaz db, err: %s", err.Error())
	}

	return jaz1Psql, nil
}
