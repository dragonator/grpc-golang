package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/dragonator/grpc-golang/internal/pb"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile("certs/localhost.pem", "certs/localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(s, new(employeeService))

	log.Println("Starting server on port " + port)
	s.Serve(lis)
}

type employeeService struct{}

func (*employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata received: %v\n", md)
	}

	for _, e := range employees {
		if e.BadgeNumber == req.BadgeNumber {
			return &pb.EmployeeResponse{Employee: e}, nil
		}
	}
	return nil, status.Error(codes.NotFound, "Employee not found")
}

func (*employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: e})
	}

	return nil
}

func (*employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	employees = append(employees, req.Employee)
	for _, e := range employees {
		fmt.Println(e)
	}
	return &pb.EmployeeResponse{Employee: req.Employee}, nil
}

func (*employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}
	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}

func (*employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Receiving photo for badge number %v\n", md["badgenumber"][0])
	}

	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received %v bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}
