package main

import (
	"bufio"
	service "chittychat/service"
	"chittychat/utils"
	"context"
	"fmt"
	"log"
	"net"
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

	go BroadcastListeningPart(&c)

	for{
		Prompt(&c)	
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

func Prompt(c *Client){
		context := context.Background()

		//var scannedLine string
		//fmt.Scanf("%s", &scannedLine)

		log.Printf("Please input a message:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		scannedLine := scanner.Text()

		message := service.Message{
			Body: scannedLine,
			LamportTimestamp: c.lamport.Value(),
			ClientName: c.name,
			Counter: 0,
		}

		//TODO skal finpudses så det scaler ift at slette døde connections, det virker pt kun med 3 servere
		empty := service.Empty{}
		counter := 0
		for i := range connections {
			_, err :=  connections[i].client.Heartbeat(context, &empty)
			if err == nil{
				message.Counter++
			}
			_, error := connections[i].client.Publish(context, &message)
			if error != nil {
				connections[i] = connections[len(connections)-1]
				counter++
			}
		}
			
		connections = connections[:len(connections)-counter]
		
}


func (c *Client) Broadcast(ctx context.Context, msg *service.Message) (*service.Empty, error) {
	msg.LamportTimestamp++
	log.Printf("%s said: %s \n", msg.ClientName, msg.Body)

	return &service.Empty{}, nil
}

func BroadcastListeningPart(c *Client) {
	lis, err := net.Listen("tcp", ":"+c.clientListningPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	service.RegisterBroadcastServiceServer(grpcServer, c)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}