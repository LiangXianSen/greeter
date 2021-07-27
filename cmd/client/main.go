package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	api "github.com/LiangXianSen/greeter/api"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := api.NewGreeterClient(conn)
	reply, err := client.Hello(context.Background(), &api.Request{Name: "Kevin"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.Msg)
}
