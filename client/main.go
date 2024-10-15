package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "github.com/Horizon-School-of-Digital-Technologies/library/api" // Replace with the actual path where the generated proto files are
	"google.golang.org/grpc"
)

func RandomBook() *pb.Book {
	return &pb.Book{
		Id:              int32(rand.Intn(10)),
		Title:           randString(20),
		Author:          randString(15),
		Isbn:            randString(15),
		PublicationYear: int32(rand.Intn(2023-1900) + 1900), // Random year between 1900 and 2023
		Genre:           randString(10),
	}
}

// randString generates a random string of a given length
func randString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// loadTest performs random queries on the library service
func loadTest(client pb.LibraryServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10; i++ { // Number of operations per goroutine
		operation := rand.Intn(4) // Randomly choose operation
		book := RandomBook()

		switch operation {
		case 0: // CreateBook
			_, err := client.CreateBook(context.Background(), &pb.CreateBookRequest{Book: book})
			if err != nil {
				log.Printf("Failed to create book: %v", err)
			} else {
				log.Printf("Created book: %v", book)
			}

		case 1: // GetBook
			_, err := client.GetBook(context.Background(), &pb.GetBookRequest{Id: book.Id})
			if err != nil {
				log.Printf("Failed to get book: %v", err)
			} else {
				log.Printf("Retrieved book: %v", book)
			}

		case 2: // UpdateBook
			_, err := client.UpdateBook(context.Background(), &pb.UpdateBookRequest{Book: book})
			if err != nil {
				log.Printf("Failed to update book: %v", err)
			} else {
				log.Printf("Updated book: %v", book)
			}

		case 3: // DeleteBook
			_, err := client.DeleteBook(context.Background(), &pb.DeleteBookRequest{Id: book.Id})
			if err != nil {
				log.Printf("Failed to delete book: %v", err)
			} else {
				log.Printf("Deleted book: %v", book.Id)
			}
		}

		// Sleep for a short duration before the next operation to mimic user delay
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Establish connection to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewLibraryServiceClient(conn)
	//
	//b1 := &pb.Book{
	//	Id:              1,
	//	Title:           "Test book",
	//	Author:          "Mohamed",
	//	Isbn:            "someISBN",
	//	PublicationYear: 2024, // Random year between 1900 and 2023
	//	Genre:           "Fiction",
	//}
	//_ , err = client.CreateBook(context.Background(), &pb.CreateBookRequest{Book: b1})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ { // Number of concurrent clients
		wg.Add(1)
		go loadTest(client, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Load testing completed.")
}
