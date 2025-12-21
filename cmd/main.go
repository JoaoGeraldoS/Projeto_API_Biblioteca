package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/docs"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/database"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/logger"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/routes"
	"github.com/joho/godotenv"

	_ "github.com/JoaoGeraldoS/Projeto_API_Biblioteca/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	_ = godotenv.Load()
}

// @title API da Biblioteca
// @version 1.0
// @description Esta é a API de gerenciamento de livros da Biblioteca.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:8080
func main() {
	host := os.Getenv("RENDER_EXTERNAL_HOSTNAME")
	if host == "" {
		docs.SwaggerInfo.Host = "localhost:8080"
	} else {
		docs.SwaggerInfo.Host = host
	}

	loggerEnv := os.Getenv("LOGGER_APP")

	loggerApp := logger.NewLogger(loggerEnv)
	defer func() {
		_ = loggerApp.Sync()
	}()

	if loggerEnv != "development" {
		loggerApp.Info("Log em modo de PRODUÇÃO")
	} else {
		loggerApp.Debug("Log em mode de DESENVOLVIMENTO")
	}

	db := database.Connection()

	r := routes.Routers(db, loggerApp)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("O servidor está rodando na porta 8080")
}
