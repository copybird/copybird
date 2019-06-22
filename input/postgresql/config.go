package postgres

var defaultSchemaName = "public"

// config stores configuration for PostgreSQL backups
type (
	config struct {
		DSN string
	}

	tableColumn struct {
		columnName, columnDefault, dataType, characterMaximumLength, isNullable, constraintName, constraintType, sequence string
	}
)
