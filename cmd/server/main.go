package main

import (
	"log"
	"net"

	"github.com/dylanbernhardt/drynklab-recipe-service/internal/recipe"
	pb "github.com/dylanbernhardt/drynklab-recipe-service/proto/recipe"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRecipeServiceServer(s, recipe.NewService())

	log.Println("Starting DrynkLab Recipe Service gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
