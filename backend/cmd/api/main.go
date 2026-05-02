// Ponte de entrada da aplicacao

package main

import (
    "golang-crud/internal/database"
    "golang-crud/internal/routes"
)

func main() {
    // Abre a conexao com o banco (postgres)
    database.Connect()

    // Configura as rotas da api
    router := routes.SetupRouter()

    // Inicia o servidor, na porta 8080
    router.Run(":8080")
}