CREATE TABLE book_category (
  id int NOT NULL AUTO_INCREMENT,
  book_id int DEFAULT NULL,
  category_id int DEFAULT NULL,
  PRIMARY KEY (id),
  KEY book_id (book_id),
  KEY category_id (category_id),
  CONSTRAINT intermediary_ibfk_1 FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
  CONSTRAINT intermediary_ibfk_2 FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);