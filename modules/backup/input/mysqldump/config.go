package mysqldump

// MySQLDumpConfig stores configuration for Mysqldump utility
type MySQLDumpConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
