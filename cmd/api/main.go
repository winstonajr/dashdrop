package main

import (
	"database/sql"
	"log"

	"github.com/winstonajr/dashdrop/internal/handler"
	"github.com/winstonajr/dashdrop/internal/repository"
	"github.com/winstonajr/dashdrop/internal/router"
	"github.com/winstonajr/dashdrop/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./dashdrop.db")
	if err != nil {
		log.Fatalf("falha ao conectar ao banco de dados: %v", err)
	}

	defer db.Close()

	repo := repository.New(db)
	fileService := service.NewFileService(repo, "./uploads")
	fileHandler := handler.NewFileHandler(fileService)

	r := router.SetupRouter(fileHandler)

	log.Println("Servidor iniciado em http://localhost:8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("erro ao executar o servidor: %v", err)
	}
}
