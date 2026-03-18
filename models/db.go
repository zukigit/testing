package models

type DBType string

const (
	DBTypePsql  DBType = "psql"
	DBTypeMysql DBType = "mysql"
	DBTypeMaria DBType = "maria"
)
