CREATE TABLE IF NOT EXISTS authors (
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK(name <> ''),
    description TEXT NOT NULL 
);
		
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK(name <> ''),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
    id INTEGER NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL CHECK(title <> ''),
    description TEXT NOT NULL CHECK(description <> ''),
    content TEXT NOT NULL CHECK(content <> ''),
    author_id INTEGER,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    foreign key(author_id) references authors(id)
);

CREATE TABLE IF NOT EXISTS book_category (
    id INTEGER NOT NULL PRIMARY KEY,
    book_id INTEGER,
    category_id INTEGER,
    foreign key(book_id) references books(id),
    foreign key(category_id) references categories(id)
);