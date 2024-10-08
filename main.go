package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "github.com/Horizon-School-of-Digital-Technologies/library/api" // Replace with the actual import path of the generated protobuf files
	"google.golang.org/grpc"
)

// BookStore struct to hold the in-memory storage
type BookStore struct {
	books map[int32]*pb.Book
	mu    sync.Mutex // Mutex to handle concurrent access
}

// LibraryServer is used to implement the LibraryService
type LibraryServer struct {
	pb.UnimplementedLibraryServiceServer
	store *BookStore
}

// CreateBook implementation
func (s *LibraryServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	if _, exists := s.store.books[req.Book.Id]; exists {
		return nil, errors.New("book with the given ID already exists")
	}

	// Add the book to the store
	s.store.books[req.Book.Id] = req.Book
	log.Printf("Book added: %v", req.Book)

	return &pb.CreateBookResponse{Book: req.Book}, nil
}

// GetBook implementation
func (s *LibraryServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	book, exists := s.store.books[req.Id]
	if !exists {
		return nil, errors.New("book not found")
	}

	return &pb.GetBookResponse{Book: book}, nil
}

// UpdateBook implementation
func (s *LibraryServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	if _, exists := s.store.books[req.Book.Id]; !exists {
		return nil, errors.New("book not found")
	}

	// Update the book
	s.store.books[req.Book.Id] = req.Book
	log.Printf("Book updated: %v", req.Book)

	return &pb.UpdateBookResponse{Book: req.Book}, nil
}

// DeleteBook implementation
func (s *LibraryServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	if _, exists := s.store.books[req.Id]; !exists {
		return nil, errors.New("book not found")
	}

	// Delete the book
	delete(s.store.books, req.Id)
	log.Printf("Book deleted: %v", req.Id)

	return &pb.DeleteBookResponse{Success: true}, nil
}

// ListBooks implementation
func (s *LibraryServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	var books []*pb.Book
	for _, book := range s.store.books {
		books = append(books, book)
	}

	return &pb.ListBooksResponse{Books: books}, nil
}

// Main function to start the server
func main() {
	// Create a BookStore
	bookStore := &BookStore{
		books: make(map[int32]*pb.Book),
	}

	// Create a new LibraryServer with the BookStore
	server := &LibraryServer{
		store: bookStore,
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the LibraryServer with the gRPC server
	pb.RegisterLibraryServiceServer(grpcServer, server)

	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server is listening on port :50051")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
