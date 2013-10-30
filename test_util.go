package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var cachedConnection []*Connection

func stringMatch(items []string, wanted string) bool {
	for _, item := range items {
		if item == wanted {
			return true
		}
	}

	return false
}
func mysqlConnectionString() string {
	if os.Getenv("TRAVIS") != "" {
		return "travis:@/db_test?charset=utf8"
	} else {
		return "root:toor@/db_test?charset=utf8"
	}
}
func postgresConnectionString() string {
	if os.Getenv("TRAVIS") != "" {
		return "user=postgres dbname=db_test"
	} else {
		return "user=postgres password=postgres dbname=db_test"
	}
}
func setupMysqlTestConn() *Connection {
	db, err := sql.Open("mysql", mysqlConnectionString())
	if err != nil {
		panic(err)
	}
	for _, line := range mysqlCreateScript {
		_, err = db.Exec(line)
		if err != nil {
			panic(err)
		}
	}

	conn, err := NewConnection("mysql", "db_test", mysqlConnectionString())
	if err != nil {
		panic(err)
	}
	conn.Config = NewSimpleConfig()
	conn.CreateMapper("Post", &post{})

	return conn
}

func setupPostgresTestConn() *Connection {
	db, err := sql.Open("postgres", postgresConnectionString())
	if err != nil {
		panic(err)
	}
	for _, line := range postgresCreateScript {
		_, err = db.Exec(line)
		if err != nil {
			panic(err)
		}
	}

	conn, err := NewConnection("postgres", "db_test", postgresConnectionString())
	if err != nil {
		panic(err)
	}
	conn.Config = NewSimpleConfig()
	conn.CreateMapper("Post", &post{})

	return conn
}

func setupSqliteTestConn() *Connection {
	conn, err := NewConnection("sqlite3", "main", ":memory:")
	if err != nil {
		panic(err)
	}
	for _, line := range sqliteCreateScript {
		_, err = conn.DB.Exec(line)
		if err != nil {
			panic(err)
		}
	}

	conn.Config = NewSimpleConfig()
	conn.CreateMapper("Post", &post{})

	return conn
}

func availableTestConns() []*Connection {
	if len(cachedConnection) == 0 {
		cachedConnection = []*Connection{
			setupMysqlTestConn(),
			setupSqliteTestConn(),
			setupPostgresTestConn(),
		}
		return cachedConnection
	}
	return cachedConnection
}

type post struct {
	Id        int
	Title     string
	Permalink string
	Body      string
	Views     int
}
type user struct {
	Id       int
	Name     string
	Email    string
	Password string
	Story    string
}

var mysqlCreateScript = []string{
	"DROP TABLE IF EXISTS `post` CASCADE;",
	"CREATE TABLE `post` ( \n" +
		"	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, \n" +
		"	`title` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`permalink` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`views` Int( 255 ) UNSIGNED NOT NULL DEFAULT '1', \n" +
		"	`author` Int( 255 ) UNSIGNED,\n" +
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

var sqliteCreateScript = []string{
	`DROP TABLE IF EXISTS "post";`,

	`CREATE TABLE "post"(
    "id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "title" Text NOT NULL,
    "permalink" Text NOT NULL,
    "body" Text NOT NULL,
    "views" Integer NOT NULL DEFAULT 1,
    "author" Integer,
  CONSTRAINT "unique_id" UNIQUE ( "id" ),
  CONSTRAINT "unique_permalink" UNIQUE ( "permalink" ) );`,

	`CREATE INDEX "index_id" ON "post"( "id" );`,

	`INSERT INTO "post"("id","title","permalink","body","views","author") VALUES ( 1,'First Post','first_post','This is the first post',1,0 );`,
	`INSERT INTO "post"("id","title","permalink","body","views","author") VALUES ( 2,'Second Post','second_post','Hey must be committed to this, I wrote a second post',1,0 );`,

	`DROP TABLE IF EXISTS "user";`,

	`CREATE TABLE "user"(
    "id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" Text NOT NULL,
    "email" Text NOT NULL,
    "password" Text NOT NULL,
    "story" Text NOT NULL,
    "image" BLOB,
  CONSTRAINT "unique_email" UNIQUE ( "email" ),
  CONSTRAINT "unique_id" UNIQUE ( "id" ),
  CONSTRAINT "unique_name" UNIQUE ( "name" ) );`,

	`CREATE INDEX "index_id1" ON "user"( "id" );`,
}

var postgresCreateScript = []string{
	`DROP TABLE IF EXISTS "post"`,
	`CREATE TABLE post(
    id bigserial NOT NULL,
    title character varying(255) NOT NULL,
    permalink character varying(255) NOT NULL,
    body text,
    views integer NOT NULL DEFAULT 1,
    CONSTRAINT pk_post PRIMARY KEY (id)
  ) WITH (
    OIDS=FALSE
  );`,
	`ALTER TABLE post OWNER TO postgres;`,
	"INSERT INTO post (title,permalink,body,views) VALUES ('First Post', 'first_post', 'This is the first post', 1);",
	"INSERT INTO post (title,permalink,body,views) VALUES ('Second Post', 'second_post', 'Hey must be committed to this, I wrote a second post', 1);",
	`DROP TABLE IF EXISTS "user"`,
	`CREATE TABLE "user" (
    id bigserial NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    story text,
    image bytea,
    name character varying(255) NOT NULL,
    CONSTRAINT pk_user PRIMARY KEY (id),
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_name UNIQUE (name)
  ) WITH (
    OIDS=FALSE
  );`,
	`ALTER TABLE "user" OWNER TO postgres;`,
}
