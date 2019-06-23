package mysql

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authorsSchema = "CREATE TABLE `authors` (\n" +
	"  `id` int(11) NOT NULL AUTO_INCREMENT,\n" +
	"  `first_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `last_name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `email` varchar(100) COLLATE utf8_unicode_ci NOT NULL,\n" +
	"  `birthdate` date NOT NULL,\n" +
	"  `added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
	"  PRIMARY KEY (`id`),\n" +
	"  UNIQUE KEY `email` (`email`)\n" +
	") ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci"
const authorsData = "(1,'Bradly','McLaughlin','purdy.richard@example.net','2001-07-21','1978-06-04 08:02:42'),(2,'Camryn','Osinski','jayce.haag@example.org','2002-04-04','1974-12-07 10:31:11'),(3,'Daryl','Haag','vleffler@example.com','1992-10-26','1978-12-05 01:11:43'),(4,'Ozella','Grant','bradtke.chaya@example.com','2013-11-20','1993-02-13 16:10:24'),(5,'Daryl','Ritchie','nader.laura@example.com','1980-02-03','1972-10-18 21:30:41'),(6,'Ellis','Donnelly','sanford.gaylord@example.org','1986-04-21','1972-02-16 16:45:19'),(7,'Clinton','Rau','shemar56@example.org','2012-09-28','1977-05-03 11:06:45'),(8,'Alberta','Denesik','xrempel@example.net','2015-04-16','1986-11-13 00:53:44'),(9,'Johanna','Nolan','anabelle20@example.com','1998-03-28','1998-04-14 18:20:08'),(10,'Monique','Pollich','gorczany.gilbert@example.net','1972-06-10','2010-02-20 03:39:31'),(11,'Tess','Rodriguez','jaren.pacocha@example.com','2001-08-01','1989-06-26 02:35:32'),(12,'Jamir','Bartell','arnoldo.smith@example.org','1970-11-23','1978-05-12 10:04:37'),(13,'Yadira','Hermiston','qdibbert@example.org','1996-07-18','1991-02-09 23:08:26'),(14,'Seth','Harris','hermiston.meta@example.com','2016-11-19','1975-01-28 13:35:52'),(15,'Judge','Fahey','trevor83@example.net','2003-01-11','1983-03-07 10:04:10'),(16,'Jarrett','Stoltenberg','kenneth.zieme@example.net','2000-07-17','2008-04-13 13:03:41'),(17,'Summer','O\\'Hara','niko85@example.com','1988-12-21','2018-03-05 18:07:26'),(18,'Erik','Green','hkulas@example.net','1979-08-06','1977-10-12 07:48:33'),(19,'Julian','Marquardt','paolo50@example.com','2017-11-23','2016-11-14 00:12:28'),(20,'Tiana','Jerde','haley.marjory@example.org','1984-09-27','2001-04-25 15:01:20'),(21,'Kaia','Connelly','dach.dorothea@example.com','1974-02-18','1997-02-25 11:30:57'),(22,'Yolanda','Lockman','julianne.johns@example.net','1998-07-02','1977-05-20 01:38:16'),(23,'Brian','Keeling','laila17@example.org','1997-03-02','1985-02-27 22:14:42'),(24,'Milford','Cartwright','bridie.pollich@example.com','2008-05-23','2013-02-02 19:54:54'),(25,'Santa','Jacobi','block.kitty@example.net','1984-10-17','2008-07-30 11:06:35'),(26,'Declan','Cartwright','rossie.hartmann@example.org','1977-12-10','2015-11-26 07:53:10'),(27,'Ressie','Gerlach','wilber.armstrong@example.org','2009-03-09','1972-04-28 16:42:28'),(28,'Rosamond','Nikolaus','bednar.price@example.com','2018-04-16','2000-11-03 00:50:30'),(29,'Tatum','Abbott','willis09@example.net','2003-02-06','2010-01-19 04:09:58'),(30,'Ismael','Skiles','marley03@example.net','2018-07-22','2009-11-24 14:08:01'),(31,'Loy','Glover','lyric10@example.com','1987-12-05','1991-05-28 19:22:25'),(32,'Margaretta','Gleichner','lelah.ziemann@example.org','1981-04-16','1974-11-29 14:16:41'),(33,'Toney','Senger','colten10@example.net','1984-12-04','1985-08-20 12:54:23'),(34,'Holden','Wilderman','thalia19@example.org','2009-08-26','1994-03-29 08:47:31'),(35,'Ottis','Muller','julie45@example.com','2006-10-19','2012-01-11 13:07:48'),(36,'Destinee','Abbott','xborer@example.org','2000-02-07','2005-01-26 18:48:05'),(37,'Blanca','Grimes','kelly.yundt@example.net','1997-11-13','1989-10-21 18:50:33'),(38,'Rachel','Herzog','kaylin53@example.com','2003-12-24','1971-09-19 10:47:01'),(39,'Kallie','Ebert','greg07@example.com','2002-02-15','1989-09-28 12:58:00'),(40,'Shaina','McCullough','tyrel57@example.net','1971-10-05','1984-04-24 14:50:15'),(41,'Madilyn','Greenholt','pierre99@example.net','1993-06-08','2013-08-06 11:51:39'),(42,'Augustine','Gutkowski','dustin.schuppe@example.org','1971-12-19','1978-08-05 12:57:02'),(43,'Magdalen','Adams','estevan.bruen@example.com','1970-08-15','2002-10-07 05:15:43'),(44,'Phoebe','Bashirian','mohammed.stokes@example.org','2001-06-14','1970-12-23 03:09:45'),(45,'Kenyon','Wilderman','cassandra12@example.com','2012-08-19','1998-10-23 01:43:03'),(46,'Salvador','O\\'Conner','jamie61@example.net','2003-03-09','1978-05-12 06:13:39'),(47,'Electa','Hoppe','qledner@example.com','1988-04-29','2008-05-01 17:24:21'),(48,'Marley','Fay','horacio.ankunding@example.org','1998-02-03','1999-05-13 04:50:26'),(49,'Gerard','Weber','ona.mills@example.org','2017-01-26','2005-11-20 22:17:49'),(50,'Zelda','Effertz','zrolfson@example.net','1995-08-06','1986-08-21 18:44:46'),(51,'Nikita','Crona','dietrich.michele@example.net','2009-06-30','2006-02-06 01:10:51'),(52,'Sanford','Little','bbradtke@example.org','1978-06-16','1975-02-07 06:21:51'),(53,'Dulce','Parisian','iwill@example.org','1991-10-28','2019-05-15 10:04:40'),(54,'Merl','Braun','qerdman@example.com','1978-07-29','2004-02-04 01:04:21'),(55,'Melvina','Bednar','beier.toy@example.com','1996-01-24','1981-04-25 03:13:24'),(56,'Erica','Swift','roxanne.dibbert@example.org','2011-09-10','2001-12-20 06:48:28'),(57,'Aniyah','Wilderman','kris.lela@example.com','2002-06-19','2002-05-25 01:04:56'),(58,'Kaya','Casper','murazik.keeley@example.com','1977-12-22','2004-02-06 11:03:08'),(59,'Jayda','O\\'Keefe','tillman.kade@example.com','1992-02-21','1976-09-03 12:36:30'),(60,'Shanna','Hammes','damon.abernathy@example.com','1970-06-22','1981-07-06 10:21:00'),(61,'Loy','Goodwin','smarquardt@example.org','2011-05-11','1983-12-02 08:57:46'),(62,'Paige','Stracke','jdaniel@example.com','1970-10-22','1986-09-15 16:09:08'),(63,'Thad','Prosacco','avis11@example.net','2018-12-17','1978-07-17 05:03:19'),(64,'Kyla','Schaden','tabitha.schinner@example.org','1985-04-27','1993-07-09 03:12:21'),(65,'Jamey','Rath','kaylin.lesch@example.org','2007-07-31','1980-01-14 19:50:45'),(66,'Jackie','McCullough','brianne74@example.net','1995-12-21','2004-01-15 10:43:08'),(67,'Norma','Dare','cecelia.mann@example.org','2003-06-27','2019-05-31 10:08:00'),(68,'Hunter','Witting','sasha71@example.net','1971-02-10','2007-10-04 06:34:12'),(69,'Izaiah','Hilpert','kunze.patricia@example.org','2015-01-19','1986-12-16 15:38:47'),(70,'Kelsie','Herman','pollich.jammie@example.org','2015-01-02','2010-03-15 20:48:55'),(71,'Maia','Haag','roxane.kohler@example.net','2011-07-28','2001-07-15 23:06:58'),(72,'Grayson','Grimes','leora.kozey@example.net','2002-10-26','2016-05-01 21:39:01'),(73,'Loma','Morissette','mariana31@example.org','1995-01-31','1972-09-30 08:02:09'),(74,'Logan','Von','russel.adriana@example.com','2015-05-29','1984-01-22 22:05:42'),(75,'Clinton','Beer','eusebio34@example.org','1975-04-08','1990-10-01 16:11:41'),(76,'Clotilde','Johnson','nquitzon@example.net','2009-12-20','2007-11-29 08:03:49'),(77,'Tevin','Friesen','marvin.rosamond@example.net','2014-04-03','1982-09-27 21:35:21'),(78,'Davon','Tromp','lia.botsford@example.net','1977-06-03','2008-01-28 21:52:30'),(79,'Robyn','Walsh','hernser@example.org','1994-08-15','1994-04-10 20:14:53'),(80,'Alyson','Yost','lreichel@example.org','1971-09-15','1995-07-29 21:06:36'),(81,'Chaz','Abbott','evan.moore@example.com','1997-05-13','1996-01-20 13:05:10'),(82,'Rosendo','Howe','gerlach.cicero@example.net','2001-03-19','2016-06-15 20:14:35'),(83,'Turner','Mertz','rippin.chloe@example.com','2000-12-31','1970-06-10 23:51:52'),(84,'Kitty','Witting','abigail09@example.net','1984-05-26','2008-09-04 15:08:06'),(85,'Rusty','Schroeder','ukiehn@example.org','2000-03-15','2010-08-11 16:34:45'),(86,'Lauren','Stamm','lowe.clementine@example.com','1995-09-15','1997-11-02 22:09:36'),(87,'Adolf','Rippin','joany.huel@example.org','1995-08-11','2012-11-21 20:08:59'),(88,'Robin','Block','lowe.maida@example.org','1985-03-27','1992-08-03 18:08:44'),(89,'Delta','Eichmann','tjohnson@example.org','2004-03-28','1982-05-13 22:12:14'),(90,'Lucy','Franecki','jgorczany@example.net','1974-10-09','2003-08-14 08:18:12'),(91,'Lucio','Sauer','helga22@example.com','2019-06-22','1997-10-30 08:43:29'),(92,'Jennifer','Moen','vfranecki@example.org','1982-12-28','2009-10-03 23:21:33'),(93,'Agnes','Stoltenberg','abshire.lura@example.com','1970-12-17','2003-04-08 01:28:05'),(94,'Adrian','Abernathy','iwaelchi@example.com','1986-08-07','2000-09-14 08:21:33'),(95,'Jalen','Wehner','hertha11@example.com','1999-01-17','1993-12-05 16:44:28'),(96,'Federico','DuBuque','ryan.freeda@example.net','2006-10-26','1990-08-03 20:44:58'),(97,'Kenna','Goodwin','jacobi.mervin@example.net','1984-10-18','1972-05-23 13:42:59'),(98,'Declan','Gislason','arch73@example.net','2019-05-31','2005-04-21 08:04:44'),(99,'Gaetano','Davis','bins.adele@example.net','1993-10-05','1984-01-13 06:24:51'),(100,'Tobin','Prohaska','kemard@example.net','2008-07-09','1993-01-09 08:10:38')"

func TestGetTables(t *testing.T) {
	d := &BackupInputMysql{}
	c := d.GetConfig().(*MySQLConfig)
	c.DSN = "root:root@tcp(localhost:3306)/test"

	require.NoError(t, d.InitModule(c))
	f, err := os.Create("dump.sql")
	assert.NoError(t, err)
	assert.NoError(t, d.InitPipe(f, nil))
	tables, err := d.getTables()
	require.NoError(t, err)
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "authors", tables[0])
	tableSchema, err := d.getTableSchema(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, authorsSchema, tableSchema)
	data, err := d.getTableData(tables[0])
	assert.NoError(t, err)
	assert.Equal(t, authorsData, data)
	err = d.Run()
	assert.NoError(t, err)
}

func TestMysqlDump(t *testing.T) {
	m := &BackupInputMysql{}

	c := m.GetConfig().(*MySQLConfig)
	c.DSN = "root:root@tcp(localhost:3306)/test"

	m.InitModule(c)
	buf := bytes.Buffer{}
	m.InitPipe(&buf, nil)
	assert.NoError(t, m.Run())
	t.Log(buf.String())
}
