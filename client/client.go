package main

import (
	"context"
	"fmt"
	pb "github.com/zhiyi57/PingCapHW/proto/kvrawpb"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:8080"
)

func main() {

	fmt.Println("Client started")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTikvClient(conn)

	fmt.Println("Proxy Server Connected")

	// Contact the server and put value to kv.
	key := []byte("Company")
	value := []byte("PingCap")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.RawPut(ctx, &pb.RawPutRequest{Key: key, Value: value})
	if err != nil {
		log.Fatalf("could not send put request: %v", err)
	}
	log.Println("sent put request to kv")

	out, err := c.RawGet(ctx, &pb.RawGetRequest{Key: key})
	if err !=nil {
		log.Fatal("could not send get request: %v", err)
	}

}
