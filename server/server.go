package main

import (
	service "chittychat/service"
	"chittychat/utils"
	"context"
	//"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type Server struct {
	service.UnimplementedChatClientServiceServer
	clientPorts []string
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
		log.Fatalf("Failed to serve at: %v", err)
	}
}

func (s *Server) Publish(ctx context.Context, msg *service.Message) (*service.Message, error) {
	msg.LamportTimestamp++
	log.Printf("%v said: %v, Lamport Timestamp is: %v", msg.ClientName, msg.Body, msg.GetLamportTimestamp())
	
	//only 9000
	
	SendMessageToClients(msg, s)
	return &service.Message{Body: "The server received the message", LamportTimestamp: msg.GetLamportTimestamp()}, nil
}

func SendMessageToClients(broadcastMessage *service.Message, s *Server) {
	s.lamport.Increment()
	// for h:= 0; h < len(s.clientPorts); h++{
	// 	log.Println(s.clientPorts[h])
	// } 

	log.Printf("Broadcasting is starting now. LamportTimestamp is: %d", s.lamport.Value())
	for i := 0; i < len(s.clientPorts); i++ {
		
		//log.Println(s.clientPorts[i])
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":"+s.clientPorts[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %s", err)
		}

		defer conn.Close()
		c := service.NewBroadcastServiceClient(conn)

		message := service.Message{
			Body:             broadcastMessage.Body,
			LamportTimestamp: s.lamport.Value(),
			ClientName:       broadcastMessage.ClientName,
		}

		_, err = c.Broadcast(context.Background(), &message)
		if err != nil {
			log.Fatalf("The message could not be broadcasted: %s", err)
		}
	}
}

func (s *Server) AddClient(ctx context.Context, msg *service.AddMessage) (*service.Acknowledgment, error) {
	
	s.clientPorts = append(s.clientPorts, msg.GetPort())
	
	log.Printf("Server registred client: %v.\n", msg.GetPort())
	log.Printf("Client: %v has Lamport %v.\n", msg.GetPort(), s.lamport.Value())
	
	return &service.Acknowledgment{Succes: true, Lamport: s.lamport.Value()}, nil
}
