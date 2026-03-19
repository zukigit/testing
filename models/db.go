package models

type DBType string

const (
	DBTypePsql  DBType = "psql"
	DBTypeMysql DBType = "mysql"
	DBTypeMaria DBType = "maria"
)

type SslMode string

const (
	SslModeDisable    SslMode = "disable"
	SslModeAllow      SslMode = "allow"
	SslModePrefer     SslMode = "prefer"
	SslModeRequire    SslMode = "require"
	SslModeVerifyCA   SslMode = "verify-ca"
	SslModeVerifyFull SslMode = "verify-full"
)
