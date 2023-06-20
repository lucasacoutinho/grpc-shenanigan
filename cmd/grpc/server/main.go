package main

import (
	"database/sql"
	"net"

	_ "github.com/lib/pq"

	"github.com/lucasacoutinho/go-grpc/internal/database"
	"github.com/lucasacoutinho/go-grpc/internal/pb"
	"github.com/lucasacoutinho/go-grpc/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@db/postgres?sslmode=disable")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(categoryDB)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
