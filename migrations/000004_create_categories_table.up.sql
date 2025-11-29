 CREATE TABLE categories (
   id int NOT NULL AUTO_INCREMENT,
   name varchar(100) NOT NULL,
   created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (id)
);