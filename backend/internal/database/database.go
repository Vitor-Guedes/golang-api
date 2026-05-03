// Responsavel pelo conexao com o database
package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// variavel global compartilhada com todo a aplicacao
var DB *gorm.DB

// Cria a conexão usando as variáveis de ambiente do Docker.
func Connect() {

    // Cria a stirng de conexao com banco
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    // Inicializa o GORM com PostgreSQL
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    
    if err != nil {
        log.Fatal("Erro ao conectar ao banco:", err)
    }

    DB = db
}