package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Horizon-School-of-Digital-Technologies/library/api" // Replace with the actual path where the generated proto files are
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a client for the LibraryService
	client := pb.NewLibraryServiceClient(conn)

	// Define a timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example: Create a new book
	createBookResponse, err := client.CreateBook(ctx, &pb.CreateBookRequest{
		Book: &pb.Book{
			Id:              2,
			Title:           "gRPC and Proto",
			Author:          "Mohamed and Others",
			Isbn:            "978-0134190441",
			PublicationYear: 2024,
			Genre:           "Fake Books",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create book: %v", err)
	}
	log.Printf("Book created: %v", createBookResponse.Book)

	// Example: Get a book by ID
	getBookResponse, err := client.GetBook(ctx, &pb.GetBookRequest{Id: 1})
	if err != nil {
		log.Fatalf("Failed to get book: %v", err)
	}
	log.Printf("Book retrieved: %v", getBookResponse.Book)

	//// Example: Update a book
	//updatedBook := &pb.Book{
	//	Id:              1,
	//	Title:           "The Go Programming Language (Updated)",
	//	Author:          "Alan Donovan",
	//	Isbn:            "978-0134190440",
	//	PublicationYear: 2016,
	//	Genre:           "Programming",
	//}
	//updateBookResponse, err := client.UpdateBook(ctx, &pb.UpdateBookRequest{Book: updatedBook})
	//if err != nil {
	//	log.Fatalf("Failed to update book: %v", err)
	//}
	//log.Printf("Book updated: %v", updateBookResponse.Book)

	// Example: List all books
	listBooksResponse, err := client.ListBooks(ctx, &pb.ListBooksRequest{})
	if err != nil {
		log.Fatalf("Failed to list books: %v", err)
	}
	for _, book := range listBooksResponse.Books {
		log.Printf("Book: %v", book)
	}

	//// Example: Delete a book by ID
	//deleteBookResponse, err := client.DeleteBook(ctx, &pb.DeleteBookRequest{Id: 1})
	//if err != nil {
	//	log.Fatalf("Failed to delete book: %v", err)
	//}
	//log.Printf("Book deleted successfully: %v", deleteBookResponse.Success)
}
