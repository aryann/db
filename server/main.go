package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/aryann/db/proto"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	address = flag.String("address", ":8080", "The address used by the server.")
)

type dbServer struct {
	proto.UnimplementedDBServer
}

func (s *dbServer) Insert(context.Context, *proto.InsertRequest) (*proto.InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
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
	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterDBServer(server, &dbServer{})
	log.Printf("Starting server at %s...", *address)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
