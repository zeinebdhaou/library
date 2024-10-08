package server

import (
	"context"
	pb "github.com/Horizon-School-of-Digital-Technologies/library/api"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

// BookStore struct to hold the in-memory storage
type BookStore struct {
	Books map[int32]*pb.Book
	mu    sync.Mutex // Mutex to handle concurrent access
}

// LibraryServer is used to implement the LibraryService
type LibraryServer struct {
	pb.UnimplementedLibraryServiceServer
	Store *BookStore
}

// CreateBook implementation
func (s *LibraryServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()

	if _, exists := s.Store.Books[req.Book.Id]; exists {
		return nil, status.Error(400, "book with the given ID already exists")
	}

	// Add the book to the store
	s.Store.Books[req.Book.Id] = req.Book
	log.Printf("Book added: %v", req.Book)

	return &pb.CreateBookResponse{Book: req.Book}, nil
}

// GetBook implementation
func (s *LibraryServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()

	book, exists := s.Store.Books[req.Id]
	if !exists {
		return nil, status.Error(404, "book not found")
	}

	return &pb.GetBookResponse{Book: book}, nil
}

// UpdateBook implementation
func (s *LibraryServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()

	if _, exists := s.Store.Books[req.Book.Id]; !exists {
		return nil, status.Error(404, "book not found")
	}

	// Update the book
	s.Store.Books[req.Book.Id] = req.Book
	log.Printf("Book updated: %v", req.Book)

	return &pb.UpdateBookResponse{Book: req.Book}, nil
}

// DeleteBook implementation
func (s *LibraryServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()

	if _, exists := s.Store.Books[req.Id]; !exists {
		return nil, status.Error(404, "book not found")
	}

	// Delete the book
	delete(s.Store.Books, req.Id)
	log.Printf("Book deleted: %v", req.Id)

	return &pb.DeleteBookResponse{Success: true}, nil
}

// ListBooks implementation
func (s *LibraryServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	s.Store.mu.Lock()
	defer s.Store.mu.Unlock()

	var books []*pb.Book
	for _, book := range s.Store.Books {
		books = append(books, book)
	}

	return &pb.ListBooksResponse{Books: books}, nil
}
