package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func stringMatch(items []string, wanted string) bool {
	for _, item := range items {
		if item == wanted {
			return true
		}
	}

	return false
}
func connectionString() string {
	if os.Getenv("CI") == "yes" {
		return "travis:@/db_test?charset=utf8"
	} else {
		return "root:toor@/db_test?charset=utf8"
	}
}
func setupDefaultConn() *Connection {
	db, err := sql.Open("mysql", connectionString())
	if err != nil {
		panic(err)
	}
	for _, line := range createScript {
		_, err = db.Exec(line)
		if err != nil {
			panic(err)
		}
	}

	conn, err := NewConnection("mysql", "db_test", connectionString())
	if err != nil {
		panic(err)
	}
	conn.Config = NewSimpleConfig()
	conn.CreateMapper("Post", &post{})

	return conn
}

type post struct {
	ID        int
	Title     string
	Permalink string
	Body      string
	Views     int
}

var createScript = []string{
	"DROP TABLE IF EXISTS `post` CASCADE;",
	"CREATE TABLE `post` ( \n" +
		"	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, \n" +
		"	`title` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`permalink` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`views` Int( 255 ) UNSIGNED NOT NULL DEFAULT '1', \n" +
		"	`author` Int( 255 ) UNSIGNED NOT NULL,\n" +
		"	 PRIMARY KEY ( `id` )\n" +
		" )\n" +
		"CHARACTER SET = utf8\n" +
		"COLLATE = utf8_general_ci\n" +
		"ENGINE = InnoDB\n" +
		"AUTO_INCREMENT = 3;\n",
	"INSERT INTO `post`(`id`,`title`,`permalink`,`body`,`views`,`author`) VALUES ( '1', 'First Post', 'first_post', 'This is the first post', '1', '0' );\n",
	"INSERT INTO `post`(`id`,`title`,`permalink`,`body`,`views`,`author`) VALUES ( '2', 'Second Post', 'second_post', 'Hey must be committed to this, I wrote a second post', '1', '0' );;\n",
	"CREATE UNIQUE INDEX `unique_id` USING BTREE ON `post`( `id` );\n",
	"CREATE UNIQUE INDEX `unique_permalink` USING BTREE ON `post`( `permalink` );\n",
	"DROP TABLE IF EXISTS `user` CASCADE;\n",
	"CREATE TABLE `user` ( \n" +
		"	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, \n" +
		"	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'nobody', \n" +
		"	`email` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`password` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`story` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`image` Blob NULL,\n" +
		"	 PRIMARY KEY ( `id` )\n" +
		" )\n" +
		"CHARACTER SET = utf8\n" +
		"COLLATE = utf8_general_ci\n" +
		"ENGINE = InnoDB\n" +
		"AUTO_INCREMENT = 1;\n",
	"CREATE UNIQUE INDEX `unique_email` USING BTREE ON `user`( `email` );\n",
	"CREATE UNIQUE INDEX `unique_id` USING BTREE ON `user`( `id` );\n",
	"CREATE UNIQUE INDEX `unique_name` USING BTREE ON `user`( `name` );",
}
