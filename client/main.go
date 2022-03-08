package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dragonator/grpc-golang/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "Command to run")
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	switch *option {
	case 1:
		SendMetadata(client)
	case 2:
		GetByBadgeNumber(client)
	case 3:
		GetAll(client)
	case 4:
		AddPhoto(client)
	case 5:
		SaveAllEmployees(client)
	}
}

// SendMetadata -
func SendMetadata(client pb.EmployeeServiceClient) {
	md := metadata.MD{}
	md["user"] = []string{"tdraganov"}
	md["password"] = []string{"my-secure-password"}
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{BadgeNumber: 1234})
}

// GetByBadgeNumber -
func GetByBadgeNumber(client pb.EmployeeServiceClient) {
	res, err := client.GetByBadgeNumber(context.Background(),
		&pb.GetByBadgeNumberRequest{BadgeNumber: 1234})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Employee)
}

// GetAll -
func GetAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(emp)
	}
}

// AddPhoto -
func AddPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("red-panda.jpg")
	if err != nil {
		log.Fatal(err)
	}

	md := metadata.New(map[string]string{"badgenumber": "1234"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal()
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		if err := stream.Send(&pb.AddPhotoRequest{Data: chunk}); err != nil {
			log.Fatal(err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// SaveAllEmployees -
func SaveAllEmployees(client pb.EmployeeServiceClient) {
	employees := []*pb.Employee{
		{
			Id:          6,
			BadgeNumber: 6789,
			FirstName:   "Alberto",
			LastName:    "West",
		},
		{
			Id:          7,
			BadgeNumber: 7899,
			FirstName:   "Indigo",
			LastName:    "Hays",
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go func() {
		for {
			emp, err := stream.Recv()
			if err == io.EOF {
				close(doneCh)
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(emp)
		}
	}()

	for _, e := range employees {
		if err := stream.Send(&pb.EmployeeRequest{Employee: e}); err != nil {
			log.Fatal(err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	<-doneCh
}
