package mysqldump

// MySQLDumpConfig stores configuration for Mysqldump utility
type MySQLDumpConfig struct {
	Host              string
	Port              string
	Username          string
	Password          string
	Database          string
	Routines          bool
	Events            bool
	Triggers          bool
	SingleTransaction bool
	ColumnStatistics  bool
}
