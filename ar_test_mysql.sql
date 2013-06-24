CREATE TABLE `user` ( 
  `id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, 
  `name` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'nobody', 
  `email` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `password` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `story` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `image` Blob NULL,
   PRIMARY KEY ( `id` )
, 
  CONSTRAINT `unique_email` UNIQUE( `email` ), 
  CONSTRAINT `unique_id` UNIQUE( `id` ), 
  CONSTRAINT `unique_name` UNIQUE( `name` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 1;

CREATE TABLE `post` ( 
  `id` Int( 255 ) UNSIGNED AUTO_INCREMENT NOT NULL, 
  `title` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `permalink` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `body` Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL, 
  `views` Int( 255 ) UNSIGNED NOT NULL DEFAULT '1', 
  `author` Int( 255 ) UNSIGNED NOT NULL,
   PRIMARY KEY ( `id` )
, 
  CONSTRAINT `unique_id` UNIQUE( `id` ), 
  CONSTRAINT `unique_permalink` UNIQUE( `permalink` ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB
AUTO_INCREMENT = 1;
