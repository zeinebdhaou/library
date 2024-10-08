package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Horizon-School-of-Digital-Technologies/library/api" // Replace with the actual import path of the generated protobuf files
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement the LibraryService
type server struct {
	pb.UnimplementedLibraryServiceServer
}

// Implement the methods for LibraryService (CRUD Operations)

// CreateBook implements the CreateBook method of LibraryService
func (s *server) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	// Business logic to create a book in the system
	// For now, return the same book back as a placeholder
	return &pb.CreateBookResponse{
		Book: req.Book,
	}, nil
}

// GetBook implements the GetBook method of LibraryService
func (s *server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	// Logic to retrieve a book by ID
	// Here, we just return a dummy book for demonstration purposes
	book := &pb.Book{
		Id:              req.GetId(),
		Title:           "Sample Book",
		Author:          "John Doe",
		Isbn:            "123456789",
		PublicationYear: 2021,
		Genre:           "Fiction",
	}
	return &pb.GetBookResponse{Book: book}, nil
}

// UpdateBook implements the UpdateBook method of LibraryService
func (s *server) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	// Logic to update a book
	return &pb.UpdateBookResponse{
		Book: req.Book,
	}, nil
}

// DeleteBook implements the DeleteBook method of LibraryService
func (s *server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	// Logic to delete a book
	return &pb.DeleteBookResponse{Success: true}, nil
}

// ListBooks implements the ListBooks method of LibraryService
func (s *server) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// Return a list of dummy books for demonstration
	books := []*pb.Book{
		{Id: 1, Title: "Book One", Author: "Author A", Isbn: "1111", PublicationYear: 2020, Genre: "Sci-Fi"},
		{Id: 2, Title: "Book Two", Author: "Author B", Isbn: "2222", PublicationYear: 2019, Genre: "Fantasy"},
	}
	return &pb.ListBooksResponse{Books: books}, nil
}

func main() {
	// Listen on a TCP port for gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC libraryServer
	grpcServer := grpc.NewServer()

	// Register the LibraryService libraryServer with gRPC
	pb.RegisterLibraryServiceServer(grpcServer, &server{})

	// Enable reflection for gRPC CLI and debugging tools
	reflection.Register(grpcServer)

	// Start the libraryServer
	log.Println("gRPC libraryServer listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
