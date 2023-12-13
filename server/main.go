package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/lucas-code42/grpc-stream/contracts"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PersonService struct {
	pb.UnimplementedPersonServiceServer
}

func (p *PersonService) CreatePerson(stream pb.PersonService_CreatePersonServer) error {
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println("do something with person...", person)
		if err = stream.Send(&pb.PersonResponse{
			Person: &pb.Person{
				Id:   "TEST_ID",
				Name: person.Name,
				Age:  person.Age,
			},
		}); err != nil {
			return err
		}
	}
}

func NewPersonService() *PersonService {
	return &PersonService{}
}

func main() {
	p := NewPersonService()

	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, p)
	reflection.Register(grpcServer)

	lisSrv, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	fmt.Println("start gRPC server")
	if err := grpcServer.Serve(lisSrv); err != nil {
		panic(err)
	}
}
