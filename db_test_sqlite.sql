------------- SQLite3 Dump File -------------

-- ------------------------------------------
-- Dump of "post"
-- ------------------------------------------

DROP TABLE IF EXISTS "post";

CREATE TABLE "post"(
	"id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"title" Text NOT NULL,
	"permalink" Text NOT NULL,
	"body" Text NOT NULL,
	"views" Integer NOT NULL DEFAULT 1,
	"author" Integer NOT NULL,
CONSTRAINT "unique_id" UNIQUE ( "id" ),
CONSTRAINT "unique_permalink" UNIQUE ( "permalink" ) );

CREATE INDEX "index_id" ON "post"( "id" );

INSERT INTO "post"("id","title","permalink","body","views","author") VALUES ( 1,'First Post','first_post','This is the first post',1,0 );
INSERT INTO "post"("id","title","permalink","body","views","author") VALUES ( 2,'Second Post
','second_post','Hey must be committed to this, I wrote a second post',1,0 );

-- ------------------------------------------
-- Dump of "user"
-- ------------------------------------------

DROP TABLE IF EXISTS "user";

CREATE TABLE "user"(
	"id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" Text NOT NULL,
	"email" Text NOT NULL,
	"password" Text NOT NULL,
	"story" Text NOT NULL,
	"image" BLOB,
CONSTRAINT "unique_email" UNIQUE ( "email" ),
CONSTRAINT "unique_id" UNIQUE ( "id" ),
CONSTRAINT "unique_name" UNIQUE ( "name" ) );

CREATE INDEX "index_id1" ON "user"( "id" );

