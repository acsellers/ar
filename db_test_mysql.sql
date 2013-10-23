-- Valentina Studio --
-- MySQL dump --
-- ---------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
-- ---------------------------------------------------------


-- CREATE TABLE "post" -------------------------------------
DROP TABLE IF EXISTS `post` CASCADE;

CREATE TABLE `post` ( 
	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, 
	`title` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`permalink` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`views` Int( 255 ) UNSIGNED NOT NULL DEFAULT '1', 
	`author` Int( 255 ) UNSIGNED NOT NULL,
	 PRIMARY KEY ( `id` )
 )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 3;
-- ---------------------------------------------------------


-- Dump data of "post" -------------------------------------
INSERT INTO `post`(`id`,`title`,`permalink`,`body`,`views`,`author`) VALUES ( '1', 'First Post', 'first_post', 'This is the first post', '1', '0' );
INSERT INTO `post`(`id`,`title`,`permalink`,`body`,`views`,`author`) VALUES ( '2', 'Second Post', 'second_post', 'Hey must be committed to this, I wrote a second post', '1', '0' );;
-- ---------------------------------------------------------


-- CREATE INDEX "unique_id" --------------------------------
CREATE UNIQUE INDEX `unique_id` USING BTREE ON `post`( `id` );
-- ---------------------------------------------------------


-- CREATE INDEX "unique_permalink" -------------------------
CREATE UNIQUE INDEX `unique_permalink` USING BTREE ON `post`( `permalink` );
-- ---------------------------------------------------------


/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
-- ---------------------------------------------------------


-- CREATE TABLE "user" -------------------------------------
DROP TABLE IF EXISTS `user` CASCADE;

CREATE TABLE `user` ( 
	`id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, 
	`name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'nobody', 
	`email` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`password` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`story` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
	`image` Blob NULL,
	 PRIMARY KEY ( `id` )
 )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 1;
-- ---------------------------------------------------------


-- Dump data of "user" -------------------------------------
-- ---------------------------------------------------------


-- CREATE INDEX "unique_email" -----------------------------
CREATE UNIQUE INDEX `unique_email` USING BTREE ON `user`( `email` );
-- ---------------------------------------------------------


-- CREATE INDEX "unique_id" --------------------------------
CREATE UNIQUE INDEX `unique_id` USING BTREE ON `user`( `id` );
-- ---------------------------------------------------------


-- CREATE INDEX "unique_name" ------------------------------
CREATE UNIQUE INDEX `unique_name` USING BTREE ON `user`( `name` );
-- ---------------------------------------------------------


/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
-- ---------------------------------------------------------


