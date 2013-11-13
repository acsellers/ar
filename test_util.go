package db

import (
	"database/sql"
	"fmt"
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
	conn.Config = NewRailsConfig()

	p := conn.MustCreateMapper("Post", &post{})
	u := conn.MustCreateMapper("User", &user{})
	createDefaultPosts(p)
	createDefaultUsers(u)

	return conn
}

func createDefaultPosts(Posts Mapper) {
	Posts.SaveAll([]post{
		post{
			Title:     "First Post",
			Permalink: "first_post",
			Body:      "This is the first-est post",
			Views:     1,
			UserId:    1,
		},
		post{
			Title:     "Second Post",
			Permalink: "second_post",
			Body:      "This is the first post",
			Views:     1,
		},
	})
}

func createDefaultUsers(Users Mapper) {
	return
}

func setupPostgresTestConn() *Connection {
	db, err := sql.Open("postgres", postgresConnectionString())
	if err != nil {
		panic(err)
	}
	for _, line := range postgresCreateScript {
		_, err = db.Exec(line)
		if err != nil {
			panic(fmt.Sprint(line, err))
		}
	}

	conn, err := NewConnection("postgres", "db_test", postgresConnectionString())
	if err != nil {
		panic(err)
	}
	conn.Config = NewRailsConfig()
	p := conn.MustCreateMapper("Post", &post{})
	u := conn.MustCreateMapper("User", &user{})
	createDefaultPosts(p)
	createDefaultUsers(u)

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

	conn.Config = NewRailsConfig()
	p := conn.MustCreateMapper("Post", &post{})
	u := conn.MustCreateMapper("User", &user{})
	createDefaultPosts(p)
	createDefaultUsers(u)

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
	UserId    int
	User      user
	*Mixin
}
type user struct {
	Id       int
	Name     string
	Email    string
	Password string
	Story    string
	Post     []post
	Posts    Scope
	*Mixin
}

var mysqlCreateScript = []string{
	"DROP TABLE IF EXISTS `posts` CASCADE;",
	"CREATE TABLE `posts` ( \n" +
		"	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, \n" +
		"	`title` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`permalink` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`views` Int( 255 ) UNSIGNED NOT NULL DEFAULT '1', \n" +
		"	`user_id` Int( 255 ) UNSIGNED,\n" +
		"	 PRIMARY KEY ( `id` )\n" +
		" )\n" +
		"CHARACTER SET = utf8\n" +
		"COLLATE = utf8_general_ci\n" +
		"ENGINE = InnoDB\n" +
		"AUTO_INCREMENT = 1;\n",
	"CREATE UNIQUE INDEX `unique_id` USING BTREE ON `posts`( `id` );\n",
	"CREATE UNIQUE INDEX `unique_permalink` USING BTREE ON `posts`( `permalink` );\n",
	"DROP TABLE IF EXISTS `users` CASCADE;\n",
	"CREATE TABLE `users` ( \n" +
		"	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, \n" +
		"	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'nobody', \n" +
		"	`email` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`password` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, \n" +
		"	`story` Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL, \n" +
		"	`image` Blob NULL,\n" +
		"	 PRIMARY KEY ( `id` )\n" +
		" )\n" +
		"CHARACTER SET = utf8\n" +
		"COLLATE = utf8_general_ci\n" +
		"ENGINE = InnoDB\n" +
		"AUTO_INCREMENT = 1;\n",
	"CREATE UNIQUE INDEX `unique_email` USING BTREE ON `users`( `email` );\n",
	"CREATE UNIQUE INDEX `unique_id` USING BTREE ON `users`( `id` );\n",
	"CREATE UNIQUE INDEX `unique_name` USING BTREE ON `users`( `name` );\n",
	"INSERT INTO `users` (`id`,`email`,`password`,`name`) VALUES ('1','user@example.com', 'id10t', 'wat');",
}

var sqliteCreateScript = []string{
	`DROP TABLE IF EXISTS "posts";`,

	`CREATE TABLE "posts"(
    "id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "title" Text NOT NULL,
    "permalink" Text NOT NULL,
    "body" Text NOT NULL,
    "views" Integer NOT NULL DEFAULT 1,
    "user_id" Integer,
  CONSTRAINT "unique_id" UNIQUE ( "id" ),
  CONSTRAINT "unique_permalink" UNIQUE ( "permalink" ) );`,

	`CREATE INDEX "index_id" ON "posts"( "id" );`,

	`DROP TABLE IF EXISTS "users";`,

	`CREATE TABLE "users"(
    "id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" Text NOT NULL,
    "email" Text NOT NULL,
    "password" Text NOT NULL,
    "story" Text NULL,
    "image" BLOB,
  CONSTRAINT "unique_email" UNIQUE ( "email" ),
  CONSTRAINT "unique_id" UNIQUE ( "id" ),
  CONSTRAINT "unique_name" UNIQUE ( "name" ) );`,

	`CREATE INDEX "index_id1" ON "users"( "id" );`,

	`INSERT INTO "users" ("id", "email","password","name") VALUES (1, 'user@example.com', 'id10t', 'wat');`,
}

var postgresCreateScript = []string{
	`DROP TABLE IF EXISTS "posts";`,
	`CREATE TABLE posts(
    id bigserial NOT NULL,
    title character varying(255) NOT NULL,
    permalink character varying(255) NOT NULL,
    body text,
    views integer NOT NULL DEFAULT 1,
    user_id integer,
    CONSTRAINT pk_posts PRIMARY KEY (id)
  ) WITH (
    OIDS=FALSE
  );`,
	`ALTER TABLE posts OWNER TO postgres;`,

	`DROP TABLE IF EXISTS "users";`,
	`CREATE TABLE "users"(
    id bigserial NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    story text,
    image bytea,
    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_name UNIQUE (name)
  ) WITH (
    OIDS=FALSE
  );`,
	`ALTER TABLE "users" OWNER TO postgres;`,
	`INSERT INTO "users" (email,password,name) VALUES ('user@example.com', 'id10t', 'wat');`,
}
