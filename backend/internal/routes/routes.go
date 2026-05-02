package routes

import (
	"golang-crud/internal/database"
	"golang-crud/internal/handlers"
	"golang-crud/internal/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Cria o router com middleware e de log e recovery
	route := gin.Default()

	// Cria ou atualiza a tabela automaticamente
	database.DB.AutoMigrate(&models.User{})

	// Define o prefixo para as rotas
	api := route.Group("/api")
	{
		api.POST("/users/paginate", handlers.PaginateUser)
		
		// Cria a rota que lista os usuarios
		api.GET("/users", handlers.GetUsers)

		// Pega o user especifico, pelo id
		api.GET("/users/:id", handlers.GetUser)

		// Cria a rota que armazena os usuarios
		api.POST("/users", handlers.CreateUser)

		// Atualiza os dados de um usuario
		api.PUT("/users/:id", handlers.UpdateUser)

		// Delete o usuario
		api.DELETE("/users/:id", handlers.DeleteUser)
	}

	return route
}