package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	pb "github.com/lucas-code42/grpc-stream/contracts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func generateRandomPersons(count int) []*pb.PersonRequest {
	var persons []*pb.PersonRequest

	for i := 0; i < count; i++ {
		person := &pb.PersonRequest{
			Name: gofakeit.Name(),
			Age:  fmt.Sprintf("%d", rand.Intn(50)+18),
		}
		persons = append(persons, person)
	}

	return persons
}

func main() {
	persons := generateRandomPersons(100)

	// Conecta ao servidor gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	// Cria um cliente gRPC
	client := pb.NewPersonServiceClient(conn)

	// Chama a função CreatePerson no servidor usando um stream bidirecional
	stream, err := client.CreatePerson(context.Background())
	if err != nil {
		log.Fatalf("Erro ao chamar CreatePerson: %v", err)
	}

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Erro ao receber resposta do servidor: %v", err)
			}

			fmt.Printf("Resposta do servidor: %v\n", response.GetPerson())
		}
	}()

	// Envia alguns dados ao servidor usando o stream
	for i := 0; i < len(persons); i++ {
		person := persons[i]

		if err := stream.Send(person); err != nil {
			log.Fatalf("Erro ao enviar dados: %v", err)
		}

		time.Sleep(time.Millisecond)
	}

	// Fecha o stream de envio
	// if err := stream.CloseSend(); err != nil {
	// 	log.Fatalf("Erro ao fechar o stream de envio: %v", err)
	// }
}
