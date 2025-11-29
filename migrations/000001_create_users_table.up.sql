  CREATE TABLE users (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  email varchar(100) NOT NULL,
  username varchar(100) NOT NULL,
  password varchar(255) NOT NULL,
  bio text,
  created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  role enum('user','admin') DEFAULT 'user',
  PRIMARY KEY (id),
  UNIQUE KEY email (email)
);