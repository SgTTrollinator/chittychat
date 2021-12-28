package main

import (
	service "chittychat/service"
	"chittychat/utils"
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type Server struct {
	service.UnimplementedChatClientServiceServer
	clients []string
	lamport utils.Lamport
}

type Client struct {
	service.UnimplementedBroadcastServiceServer
	clientListningPort string
	name               string
	lamport            int32
}


func main() {

	args := os.Args

	port := ":" + args[1]

	//Create a Server with the port parameter
	setupServer(port)
}

func setupServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	server := grpc.NewServer()
	s := Server{}

	service.RegisterChatClientServiceServer(server, &s)

	log.Printf("Server Port: %v \n", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *Server) Publish(ctx context.Context, msg *service.Message) (*service.Message, error) {
	return nil, nil
}

func (s *Server) AddClient(ctx context.Context, msg *service.AddMessage) (*service.Acknowledgment, error) {
	
	s.clients = append(s.clients, msg.GetPort())
	
	log.Printf("Server registred client: %v.\n", msg.GetPort())
	log.Printf("Client: %v has Lamport %v.\n", msg.GetPort(), s.lamport.Value())
	
	return &service.Acknowledgment{Succes: true, Lamport: s.lamport.Value()}, nil
}
