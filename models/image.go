package models

type Image string

const (
	JazServer609ImagePsql  Image = "jobarg-server-postgres:6.0.9-1"
	ZabbixServer6ImagePsql Image = "zabbix/zabbix-server-pgsql:alpine-6.0-latest"
	ZabbixWeb6ImagePsql    Image = "zabbix/zabbix-web-nginx-pgsql:alpine-6.0-latest"
	Postgres14Image        Image = "postgres:14-alpine3.22"
)
