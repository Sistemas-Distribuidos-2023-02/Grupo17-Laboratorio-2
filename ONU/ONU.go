package main

import (
	"fmt"
	"log"
	"google.golang.org/grpc"
	pb "main/proto"
	"context"
	"net"
	"time"

)

type server struct {
	pb.UnimplementedOMSServer
}

func MandarDataOMS(Persona string, conn *grpc.ClientConn) {
	client := pb.NewOMSClient(conn)
    stream, err := client.NotifyBidirectional(context.Background())
    if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
    }
	mensaje := &pb.Request{Message: Persona}
    if err := stream.Send(mensaje); err != nil {
        log.Fatalf("Error al enviar mensaje: %v", err)
    }
}


func (s *server) NotifyBidirectional(steam pb.OMS_NotifyBidirectionalServer) error {
		request, _ := steam.Recv()
		fmt.Println( request.Message)
	
	return nil
}

func Preguntar (){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
    }
    defer conn.Close()
	var  PersonaBuscada string
	for {
		time.Sleep(3*time.Second)
		fmt.Print("Datos de Muertos/Infectados: ")
		_ , err := fmt.Scanln(&PersonaBuscada)
		if err != nil {
	        fmt.Println("Error al leer la entrada:", err)
        	return
		}
		if PersonaBuscada == "Muertos" {
			PersonaBuscada = "1"
		} else {
			PersonaBuscada = "2"
		}
		MandarDataOMS(PersonaBuscada, conn)
	}
}
func main() {
	

	go Preguntar()
	listener, err := net.Listen("tcp", "localhost:50070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := grpc.NewServer()
	pb.RegisterOMSServer(serv, &server{})
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)	
	}
	}

	