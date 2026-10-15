package client

import (
	"TaskManagerApi/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewTaskServiceClient(conn)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	response, err := client.GetTask(ctx,
		&proto.GetTaskRequest{
			Id: 1,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Task ID:", response.GetId())
	fmt.Println("Title", response.GetTitle())
	fmt.Println("Description", response.GetDescription())
	fmt.Println("Completed", response.GetCompleted())

}
