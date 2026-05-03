package models

import (
	"gorm.io/gorm"
	"golang-crud/pkg/types"
)

type User struct {
	// estende as variaveis do framework 
	// gorm.Model
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt types.JSONTime `json:"created_at"`
	UpdatedAt types.JSONTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// tag json:"..." controla o nome dos campos na API
	Name      string `json:"name"`

	// gorm:unique garante que o campo seja unico
	Email     string `json:"email" gorm:"unique"`
	Status    bool   `json:"status"`
}

func (User) FilterableFields() map[string]string {
    return map[string]string{
        "id":         "id",
        "name":       "name",
        "email":      "email",
        "status":     "status",
        "created_at": "created_at",
    }
}