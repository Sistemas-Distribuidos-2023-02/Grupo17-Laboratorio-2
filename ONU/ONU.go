package main

import (
	"fmt"
	"log"
	"google.golang.org/grpc"
	pb "main/proto"
	"context"

)

func MandarDataOMS(Persona string, conn *grpc.ClientConn) {
	client := pb.NewOMSClient(conn)
    stream, err := client.NotifyBidirectional(context.Background())
	fmt.Println(Persona)
    if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
    }
	mensaje := &pb.Request{Message: Persona}
    if err := stream.Send(mensaje); err != nil {
        log.Fatalf("Error al enviar mensaje: %v", err)
    }

}

func main() {
	var  PersonaBuscada string
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
    }
    defer conn.Close()
	for {
		fmt.Print("Datos de Muertos/Infectados: ")
		_ , err := fmt.Scanln(&PersonaBuscada)
		if err != nil {
	        fmt.Println("Error al leer la entrada:", err)
        	return
		}
		MandarDataOMS(PersonaBuscada, conn)
	}
	
}