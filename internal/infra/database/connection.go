package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func Connection() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o env")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal("Erro ao verificar conexão")
	}

	log.Print("Conexão realizada")

	return db
}

func SetupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	// db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("Erro ao ativar foreign keys: %v", err)
	}
	_, err = db.Exec(`
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
			title VARCHAR(100) NOT NULL CHECK(name <> ''),
			description TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			author_id INTEGER,
			foreign key(author_id) references authors(id)
		);

		CREATE TABLE IF NOT EXISTS book_category (
			id INTEGER NOT NULL PRIMARY KEY,
			book_id INTEGER DEFAULT NULL,
			category_id INTEGER DEFAULT NULL,
			foreign key(book_id) references books(id),
			foreign key(category_id) references categories(id)
		);
	`)

	if err != nil {
		db.Close()
		log.Fatalf("Erro ao criar tabelas no MySQL: %v", err)
	}

	return db
}
