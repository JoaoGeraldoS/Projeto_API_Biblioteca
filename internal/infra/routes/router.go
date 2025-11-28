package routes

import (
	"database/sql"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	BookHandler *books.BookHandler
}

func NewApp(db *sql.DB) *App {
	bookRepo := books.NewBookRepository(db)
	bookSvc := books.NewBookService(bookRepo)
	bookHandler := books.NewBookHandler(bookSvc)

	return &App{
		BookHandler: bookHandler,
	}
}

func Routers(db *sql.DB) *gin.Engine {
	r := gin.Default()

	app := NewApp(db)

	r.Use(middleware.ErrorHandler())

	routersBook(r.Group("/"), app.BookHandler)

	return r
}

func routersBook(r *gin.RouterGroup, h *books.BookHandler) {
	books := r.Group("/api/books")

	books.POST("/", h.CreateBook)
}
