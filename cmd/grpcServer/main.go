package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/Soter-Tec/grpc/internal/database"
	"github.com/Soter-Tec/grpc/internal/pb"
	"github.com/Soter-Tec/grpc/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Conexão com o banco de dados SQLite
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializa o serviço de categorias
	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

	// Configura e registra o servidor gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	// Inicia o servidor gRPC na porta 50051
	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
