package main

import (
	"context"
	pb "github.com/Horizon-School-of-Digital-Technologies/library/api" // Replace with the actual import path of the generated protobuf files
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
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

// Prometheus metrics
var (
	grpcRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	grpcRequestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_request_errors_total",
			Help: "Total number of gRPC errors",
		},
		[]string{"method"},
	)
)

// CreateBook implementation
func (s *LibraryServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	if _, exists := s.store.books[req.Book.Id]; exists {
		return nil, status.Error(400, "book with the given ID already exists")
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
		return nil, status.Error(404, "book not found")
	}

	return &pb.GetBookResponse{Book: book}, nil
}

// UpdateBook implementation
func (s *LibraryServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	if _, exists := s.store.books[req.Book.Id]; !exists {
		return nil, status.Error(404, "book not found")
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
		return nil, status.Error(404, "book not found")
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

// Unary interceptor to collect Prometheus metrics
func prometheusUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Start time for measuring latency
	startTime := time.Now()

	// Handle the request
	resp, err := handler(ctx, req)

	// Measure request duration
	duration := time.Since(startTime).Seconds()

	// Extract method name
	method := info.FullMethod

	// Collect metrics
	grpcRequests.WithLabelValues(method, grpc.Code(err).String()).Inc()
	grpcRequestDuration.WithLabelValues(method).Observe(duration)

	if err != nil {
		grpcRequestErrors.WithLabelValues(method).Inc()
	}

	return resp, err
}

// Function to expose Prometheus metrics
func exposePrometheusMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Prometheus metrics exposed on :2112/metrics")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Failed to expose Prometheus metrics: %v", err)
		}
	}()
}

func main() {
	// Register Prometheus metrics
	prometheus.MustRegister(grpcRequests)
	prometheus.MustRegister(grpcRequestDuration)
	prometheus.MustRegister(grpcRequestErrors)

	// Expose Prometheus metrics
	exposePrometheusMetrics()

	// Create a BookStore
	bookStore := &BookStore{
		books: make(map[int32]*pb.Book),
	}

	// Create a new LibraryServer with the BookStore
	server := &LibraryServer{
		store: bookStore,
	}

	// Create a new gRPC server with the Prometheus interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(prometheusUnaryInterceptor),
	)

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
