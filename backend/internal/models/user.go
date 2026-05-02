package models

import "gorm.io/gorm"

type User struct {
	// estende as variaveis do framework 
	gorm.Model

	// tag json:"..." controla o nome dos campos na API
	Name   string `json:"name"`

	// gorm:unique garante que o campo seja unico
	Email  string `json:"email" gorm:"unique"`
	Status bool   `json:"status"`
}