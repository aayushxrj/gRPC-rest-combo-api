package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	mainpb "github.com/aayushxrj/gRPC-rest-combo-api/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"

	// gw "github.com/aayushxrj/gRPC-rest-combo-api/proto/gen" // alias for gateway registration
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type server struct {
	mainpb.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, req *mainpb.AddRequest) (*mainpb.AddResponse, error) {
	// Validate the request (protoc-gen-validate will generate a Validate() method)
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	a := req.GetA()
	b := req.GetB()

	sum := a + b
	log.Printf("Add called with a=%d, b=%d, sum=%d", a, b, sum)

	return &mainpb.AddResponse{
		Sum: sum,
	}, nil
}

func (s *server) GenerateFibonacci(req *mainpb.FibonacciRequest, stream mainpb.Calculator_GenerateFibonacciServer) error {
	// Validate request
	if err := req.Validate(); err != nil {
		return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata recieved")
	}
	fmt.Println("Metadata recieved:", md)
	val, ok := md["authorization"]
	if !ok {
		log.Println("No metadata recieved")
	}
	log.Println("Authorization:", val)

	// Response headers to client
	responseHeaders := metadata.Pairs("test", "testing1", "test2", "testing2")
	if err := stream.SendHeader(responseHeaders); err != nil {
		return err
	}

	n := req.GetN()
	a, b := 0, 1

	for i := 0; i < int(n); i++ {
		err := stream.Send(&mainpb.FibonacciResponse{
			Number: int32(a),
		})
		log.Println("Sent number:", a)
		if err != nil {
			return err
		}
		a, b = b, a+b
		time.Sleep(time.Second)
	}

	trailer := metadata.New(map[string]string{
		"end-status":   "completed",
		"processed-by": "fibonacci-service",
	})
	stream.SetTrailer(trailer)

	return nil
}

func (s *server) SendNumbers(stream mainpb.Calculator_SendNumbersServer) error {
	var sum int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mainpb.NumberResponse{Sum: sum})
		}
		if err != nil {
			return err
		}

		// Validate each incoming request
		if err := req.Validate(); err != nil {
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}

		log.Println(req.GetNumber())
		sum += req.GetNumber()
	}
}

func (s *server) Chat(stream mainpb.Calculator_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Validate incoming chat messages
		if err := req.Validate(); err != nil {
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}

		log.Println("Received Message:", req.GetMessage())

		err = stream.Send(&mainpb.ChatMessage{
			Message: req.GetMessage(),
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Returning control")
	return nil
}

func main() {
	port := ":50051"
	// cert := "cert.pem"
	// key := "key.pem"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	// creds, err := credentials.NewServerTLSFromFile(cert, key)
	// if err != nil {
	// 	log.Fatal("Failed to load credentials:", err)
	// }
	// grpcServer := grpc.NewServer(grpc.Creds(creds))
	grpcServer := grpc.NewServer()

	mainpb.RegisterCalculatorServer(grpcServer, &server{})

	// enable reflection
	reflection.Register(grpcServer)

	// log.Printf("Server is running on the port %s", port)
	// err = grpcServer.Serve(lis)
	// if err != nil {
	// 	log.Fatal("Failed to serve:", err)
	// }

	// Start gRPC server in goroutine
	go func() {
		log.Printf("gRPC server running on %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC:", err)
		}
	}()

	// Start REST gateway
	restPort := ":8080"
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		// opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))}
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Register gateway
		if err := mainpb.RegisterCalculatorHandlerFromEndpoint(ctx, mux, "localhost"+port, opts); err != nil {
			log.Fatal("Failed to register gateway:", err)
		}

		log.Printf("REST gateway running on %s", restPort)
		if err := http.ListenAndServe(restPort, mux); err != nil {
			log.Fatal("Failed to serve REST gateway:", err)
		}
	}()

	// Block forever
	select {}
}
