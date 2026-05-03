// Define o nome do pacote
package handlers

// Faz o import dos modulos necessarios
import (
	// Import para servidor e rotas
	"net/http"

	// Imports internaos do banco e da model
	"golang-crud/internal/database"
	"golang-crud/internal/models"

	"golang-crud/pkg/paginate"

	"github.com/gin-gonic/gin"
)

func GetUsers(context *gin.Context) {
	// Cria variavel
	var users []models.User

	// Busca todos os registros
	database.DB.Find(&users)

	// Retorna o json com status 200
	context.JSON(http.StatusOK, users)
}

func GetUser(context *gin.Context) {
	var user models.User

	id := context.Param("id")

	if err := database.DB.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	context.JSON(http.StatusOK, user)
}

func CreateUser(context *gin.Context) {
	// Cria a variavel
	var user models.User

	// Recebe o json da requisicao. E converte o json em struct GO
	if err := context.ShouldBindJSON(&user); err != nil {
		
		// Se houver erro de validaca/sintaxe retorna erro
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		
		return
	}

	// Cria o registro no banco
	database.DB.Create(&user)

	// Retorna na api com status 201
	context.JSON(http.StatusCreated, user)
}

func UpdateUser(context *gin.Context) {
	var user models.User

	id := context.Param("id")

	if err := database.DB.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User

	// Recebe o json da requisicao. E converte o json em struct GO
	if err := context.ShouldBindJSON(&input); err != nil {
		
		// Se houver erro de validaca/sintaxe retorna erro
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		
		return
	}

	database.DB.Model(&user).Updates(input)

	context.JSON(http.StatusOK, user)
}

func DeleteUser(context *gin.Context) {
	var user models.User

	id := context.Param("id")

	if err := database.DB.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Delete(&user)

	context.JSON(http.StatusOK, gin.H{"message": "User Deleted!"})
}

func PaginateUser(context *gin.Context) {
	var users         []models.User
	var searchRequest paginate.SearchRequest

	// Recebe o json da requisicao. E converte o json em struct GO
	if err := context.ShouldBindJSON(&searchRequest); err != nil {
		
		// Se houver erro de validaca/sintaxe retorna erro
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		
		return
	}

	// Aplica o filtro de busca
	result, err := searchRequest.Apply(&users)

	// Valida se tem algum erro
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, result)
}
