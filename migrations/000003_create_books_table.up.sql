CREATE TABLE books (
     id int NOT NULL AUTO_INCREMENT,
     title varchar(100) NOT NULL,
     description text NOT NULL,
     content text NOT NULL,
     created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     author_id int DEFAULT NULL,
     PRIMARY KEY (id),
     KEY author_id (author_id),
     CONSTRAINT books_ibfk_1 FOREIGN KEY (author_id) REFERENCES authors (id)
);
