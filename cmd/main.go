package main

import (
	"log"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/database"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/routes"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o env")
	}
}

func main() {
	db := database.Connection()

	r := routes.Routers(db)

	log.Println("Server on :8080")
	r.Run(":8080")
}
