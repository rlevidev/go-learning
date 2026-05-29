package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"url-shortening/internal/db"
	"url-shortening/internal/handlers"
	"url-shortening/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:root@localhost:5432/postgres?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.Ping(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v\n", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	repo := db.New(conn)
	svc := service.New(repo)
	handler := handlers.New(svc)

	http.HandleFunc("/shorten", handler.PostShorten)
	http.HandleFunc("/{shortCode}", handler.GetShorten)
	http.HandleFunc("/{shortCode}/stats", handler.GetStats)

	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
