package postgres

var defaultSchemaName = "public"

// Config stores configuration for PostgreSQL backups
type (
	Config struct {
		DSN string
	}

	tableColumn struct {
		columnName, columnDefault, dataType, characterMaximumLength, isNullable, constraintName, constraintType, sequence string
	}
)
