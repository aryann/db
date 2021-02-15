package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/aryann/db/coordinator"
	"github.com/aryann/db/proto"
	"github.com/aryann/db/storage/json"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	address    = flag.String("address", ":8080", "The address used by the server.")
	storageDir = flag.String("storage-dir", "", "A directory used for storing the database data.")
)

type dbServer struct {
	proto.UnimplementedDBServer
	coordinator *coordinator.Coordinator
}

func (s *dbServer) Insert(ctx context.Context, request *proto.InsertRequest) (*proto.InsertResponse, error) {
	err := s.coordinator.Insert(request.GetDocument().GetKey(), request.GetDocument().GetPayload())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &proto.InsertResponse{}, nil
}

func (s *dbServer) Update(context.Context, *proto.UpdateRequest) (*proto.UpdateRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *dbServer) Delete(context.Context, *proto.DeleteRequest) (*proto.EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func (s *dbServer) Lookup(context.Context, *proto.LookupRequest) (*proto.LookupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}

func (s *dbServer) Scan(*proto.ScanRequest, proto.DB_ScanServer) error {
	return status.Errorf(codes.Unimplemented, "method Scan not implemented")
}

func main() {
	flag.Parse()

	if *storageDir == "" {
		log.Fatal("-storage-dir must not be empty")
	}
	storage, err := json.NewJSONStorage(*storageDir)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterDBServer(server, &dbServer{
		coordinator: coordinator.NewCoordinator(storage),
	})
	log.Printf("Starting server at %s...", *address)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
