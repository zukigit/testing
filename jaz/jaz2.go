package jaz

import "github.com/zukigit/testing/lib"

type Jaz2 struct {
	DBUsername, DBPassword, DBName, DBDnsName, DBPort string
}

func NewJaz2() Jaz {
	return &Jaz2{
		DBUsername: lib.Getenv("JAZ_DB_USERNAME", "zabbix"),
		DBPassword: lib.Getenv("JAZ_DB_PASSWORD", "zabbix"),
		DBName:     lib.Getenv("JAZ_DB_NAME", "jobarranager"),
	}
}
