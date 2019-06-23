package postgres

var defaultSchemaName = "public"

// Config stores configuration for PostgreSQL backups
type (
	Config struct {
		DSN string
	}

	tableScheme struct {
		columnName, columnDefault, dataType, characterMaximumLength, isNullable, constraintName, constraintType, sequence string
	}
	sequenceScheme struct {
		name string
	}
)
