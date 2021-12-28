package main

import (
	service "chittychat/service"
	"chittychat/utils"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

type Connection struct {
	clientConn *grpc.ClientConn
	client     service.ChatClientServiceClient
	context    context.Context
	port       string
}

type Client struct {
	service.UnimplementedBroadcastServiceServer
	clientListningPort string
	name               string
	lamport            utils.Lamport
}

var connections []Connection

var serverPorts = []string{":9000", ":9001", ":9002"}

func main() {

	UserInput := os.Args[1:]
	ReadPort := UserInput[0]
	ReadName := UserInput[1]

	c := Client{clientListningPort: ReadPort, name: ReadName}

	// var conn *grpc.ClientConn
	// conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("could not connect: %s", err)
	// }
	// c := chat.NewChatClientServiceClient(conn)
	// defer conn.Close()

	// UserInput := os.Args[1:]
	// ReadPort := UserInput[0]
	// ReadName := UserInput[1]

	// client := Client{clientListningPort: ReadPort, name: ReadName}

	// JoinClientPart(&client, c)

	// go BroadcastListeningPart(&client)
	// ChooseCommand(&client, c)


	for i := range serverPorts {
		ctx, conn, c := setupConnection(i, &c)
		
		newConn := Connection{
			context:    ctx,
			clientConn: conn,
			client:     c,
			port:       serverPorts[i],
		}

		connections = append(connections, newConn)

		defer newConn.clientConn.Close()
	}

	for{

	}

}

func setupConnection(index int, c *Client) (context.Context, *grpc.ClientConn, service.ChatClientServiceClient) {
	
	context := context.Background()

	conn, err := grpc.Dial(serverPorts[index], grpc.WithInsecure())

	if err != nil {
		log.Printf("Error: %v", err)
	}

	client := service.NewChatClientServiceClient(conn)

	JoinClientToServer(c, client) 

	fmt.Printf("Connecting to: %v \n", c.clientListningPort)
	return context, conn, client
}

func JoinClientToServer(client *Client, serviceClient service.ChatClientServiceClient) {

	client.lamport = *utils.NewLamport()

	addClientMessage := service.AddMessage{
		ClientName:       client.name,
		Port:             client.clientListningPort,
		LamportTimestamp: client.lamport.Value(),
	}
	
	response, err := serviceClient.AddClient(context.Background(), &addClientMessage)
	if err != nil {
		log.Fatalf("The client could not join: %v", err)
	}

	client.lamport.Increment()
	log.Printf("%v. The Lamport Timestamp is: %v", response.Succes, response.Lamport)
}