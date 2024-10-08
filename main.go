package main

import (
	"context"
	pb "github.com/Horizon-School-of-Digital-Technologies/library/api"
	sv "github.com/Horizon-School-of-Digital-Technologies/library/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

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

	// Create a new LibraryServer
	server := sv.NewLibraryServer()

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
