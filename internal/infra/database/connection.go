package database

import (
	"database/sql"
	_ "embed"
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

	return db
}

//go:embed schema.sql
var schema string

func SetupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("Erro ao ativar foreign keys: %v", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}

	return db
}
