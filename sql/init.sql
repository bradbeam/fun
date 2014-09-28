CREATE DATABASE IF NOT EXISTS fun;

CREATE TABLE IF NOT EXISTS fun.account (
  id INT NOT NULL AUTO_INCREMENT,
  username VARCHAR(64),
  password VARCHAR(256),
  PRIMARY KEY(id)
) ENGINE = INNODB;

CREATE TABLE IF NOT EXISTS fun.characters (
  id INT NOT NULL AUTO_INCREMENT,
  account_id INT,
  name VARCHAR(64) NOT NULL,
  level INT UNSIGNED NOT NULL,
  health INT UNSIGNED NOT NULL,
  mana INT UNSIGNED NOT NULL,
  attackpower INT UNSIGNED NOT NULL,
  defensepower INT UNSIGNED NOT NULL,
  status VARCHAR(64),
  kills INT UNSIGNED NOT NULL,
  experience INT UNSIGNED NOT NULL,
  played INT UNSIGNED NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (account_id) REFERENCES fun.account(id) ON UPDATE CASCADE ON DELETE CASCADE,
  UNIQUE  KEY (account_id, name )
) ENGINE = INNODB;

CREATE TABLE IF NOT EXISTS fun.monsters (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  level INT UNSIGNED NOT NULL,
  health INT UNSIGNED NOT NULL,
  mana INT UNSIGNED NOT NULL,
  attackpower INT UNSIGNED NOT NULL,
  defensepower INT UNSIGNED NOT NULL,
  status VARCHAR(64) NOT NULL,
  experience INT UNSIGNED NOT NULL,
  PRIMARY KEY (id),
  UNIQUE  KEY (name)
) ENGINE = INNODB;

REPLACE INTO fun.account values
  ( 1, "bob", "bunnies"),
  ( 2, "joe", "unicorns");

REPLACE INTO fun.characters values
  ( 1, 1, "ohai", 1, 100, 100, 5, 5, "", 0, 0, 0),
  ( 2, 1, "nohai", 1, 100, 100, 5, 5, "", 0, 0, 0),
  ( 3, 2, "fluffy", 1, 100, 100, 5, 5, "", 0, 0, 0),
  ( 4, 1, "skittles", 1, 100, 100, 5, 5, "", 0, 0, 0),
  ( 5, 2, "tootes", 1, 100, 100, 5, 5, "", 0, 0, 0);

REPLACE INTO fun.monsters values
  ( 1,'slime',1,10,0,5,10,"",0 ),
  ( 2,'imp',1,10,0,7,8,"",0 ),
  ( 3,'goblin',1,10,0,10,5,"",0 ),
  ( 4, 'dragon', 10, 1000, 1000, 100, 100, "", 0 );

GRANT ALL ON fun.* to 'fun'@'%' IDENTIFIED BY 'Passw0rd5Suck!';
FLUSH PRIVILEGES;
