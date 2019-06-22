package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const postsSchema = "CREATE TABLE `authors` (\n" +
	"  `id` int(11) NOT NULL AUTO_INCREMENT,\n" +
	"  `first_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `last_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `email` varchar(100) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `birthdate` date NOT NULL,\n" +
	"  `added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
	"  PRIMARY KEY (`id`),\n" +
	"  UNIQUE KEY `email` (`email`)\n" +
	") ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"

func TestGetTables(t *testing.T) {
	d, err := New("root:root@tcp(localhost:3306)/test")
	require.NoError(t, err)
	tables, err := d.getTables()
	require.NoError(t, err)
	assert.Equal(t, 2, len(tables))
	assert.Equal(t, "authors", tables[0])
	assert.Equal(t, "posts", tables[1])
	tableSchema, err := d.getTableSchema(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, postsSchema, tableSchema)
}
