package routes

import (
	"database/sql"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	BookHandler     *books.BookHandler
	AuthorHandler   *authors.AuthorHandler
	CategoryHandler *categories.CategoryHandler
	UserHandler     *users.UserHandler
}

func NewApp(db *sql.DB, logApp *zap.Logger) *App {
	bookRepo := books.NewBookRepository(db)
	bookSvc := books.NewBookService(bookRepo)
	bookHandler := books.NewBookHandler(bookSvc, logApp)

	authRepo := authors.NewAuthorsRepository(db)
	authSvc := authors.NewAuthorsService(authRepo)
	authHandler := authors.NewAuthorsHandler(authSvc, logApp)

	catRepo := categories.NewCategoryRepository(db)
	catSvc := categories.NewCategoryService(catRepo)
	catHanlder := categories.NewCategoryHandler(catSvc, logApp)

	userRepo := users.NewUsersRepository(db)
	userSvc := users.NewUsersService(userRepo)
	userHandler := users.NewUsersHandler(userSvc, logApp)

	return &App{
		BookHandler:     bookHandler,
		AuthorHandler:   authHandler,
		CategoryHandler: catHanlder,
		UserHandler:     userHandler,
	}
}

func Routers(db *sql.DB, logApp *zap.Logger) *gin.Engine {
	r := gin.Default()

	app := NewApp(db, logApp)

	r.Use(middleware.ErrorHandler())

	public := r.Group("/public")

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	routersBook(protected, public, app.BookHandler)
	routesUsers(protected, public, app.UserHandler)
	routersAuthors(protected, public, app.AuthorHandler)
	routersCategories(protected, public, app.CategoryHandler)

	return r
}

func routersBook(pr *gin.RouterGroup, pl *gin.RouterGroup, h *books.BookHandler) {
	booksPr := pr.Group("/api/books")
	booksPl := pl.Group("/api/books")

	booksPr.POST("/", middleware.RequireRole("admin"), h.CreateBook)
	booksPr.POST("/relation", middleware.RequireRole("admin"), h.RelationBookCategory)
	booksPr.PUT("/:id", middleware.RequireRole("admin"), h.UpdateBook)
	booksPr.DELETE("/:id", middleware.RequireRole("admin"), h.DeleteBook)

	booksPl.GET("/", h.ReadAllBooks)
	booksPl.GET("/:id", h.ReadBook)
}

func routersAuthors(pr *gin.RouterGroup, pl *gin.RouterGroup, h *authors.AuthorHandler) {
	authorsPr := pr.Group("/api/authors")
	authorsPl := pl.Group("/api/authors")

	authorsPr.POST("/", middleware.RequireRole("admin"), h.CreateAuthor)
	authorsPr.PUT("/:id", middleware.RequireRole("admin"), h.UpdateAuthor)
	authorsPr.DELETE("/:id", middleware.RequireRole("admin"), h.DeleteAuthor)

	authorsPl.GET("/", h.ReadAuthors)
	authorsPl.GET("/:id", h.ReadAuthor)
}

func routersCategories(pr *gin.RouterGroup, pl *gin.RouterGroup, h *categories.CategoryHandler) {
	categoriesPr := pr.Group("/api/categories")
	categoriesPl := pl.Group("/api/categories")

	categoriesPr.POST("/", middleware.RequireRole("admin"), h.CreateCategory)
	categoriesPr.PUT("/:id", middleware.RequireRole("admin"), h.UpdateCategory)
	categoriesPr.DELETE("/:id", middleware.RequireRole("admin"), h.DeleteCategory)

	categoriesPl.GET("/", h.ReadCategories)
	categoriesPl.GET("/:id", h.ReadCategory)
}

func routesUsers(pr *gin.RouterGroup, pl *gin.RouterGroup, h *users.UserHandler) {
	usersPr := pr.Group("/api/users")
	usersPl := pl.Group("/api/users")

	usersPr.PUT("/:id", middleware.RequireRole("admin"), h.UpdateUser)
	usersPr.DELETE("/:id", middleware.RequireRole("admin"), h.DeleteUser)

	usersPr.GET("/", h.ReadAllUsers)
	usersPl.POST("/login", h.LoginUser)
	usersPl.POST("/", h.CreateUser)
	usersPl.GET("/:id", h.ReadUser)
}
