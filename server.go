package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	mainpb "github.com/aayushxrj/gRPC-rest-combo-api/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

type server struct {
	mainpb.UnimplementedCalculatorServer
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
	cert := "cert.pem"
	key := "key.pem"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatal("Failed to load credentials:", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	mainpb.RegisterCalculatorServer(grpcServer, &server{})

	// enable reflection
	reflection.Register(grpcServer)

	log.Printf("Server is running on the port %s", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
