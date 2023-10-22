package main

import (
	"net"
	"log"
	pb "main/proto"
	"google.golang.org/grpc"
	"strings"
	"fmt"
	"strconv"
	"os"
	"context"
	"sync"
)

var mu sync.Mutex

var personas map[int]string
var archivo *os.File

type server struct {
	pb.UnimplementedOMSServer
}

func (s *server) NotifyBidirectional(steam pb.OMS_NotifyBidirectionalServer) error {
		mu.Lock()
		request, _ := steam.Recv()
		id, nombre := ObtenerIDNombre(request.Message)
		if nombre == "P" {
			fmt.Println("La ONU a preguntado por los datos de: " + strconv.Itoa(id))
			MandarDataOMS(id)
			
		}else {
			
			personas[id] = nombre
			log.Printf("Mensaje recibido: %s", request.Message)
			
			fmt.Println(personas[id])
			archivo.WriteString( strconv.Itoa(id) + " " + nombre + "\n")
		}
		mu.Unlock()
		return nil
}	

func MandarDataOMS(id int) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("No se pudo conectar al servidor: %v", err)
			}
	client := pb.NewOMSClient(conn)
    stream, err := client.NotifyBidirectional(context.Background())
    if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
    }
	mensaje := &pb.Request{Message: (personas[id] + "\nData")}
	print("\n Se ha enviado el mensaje: ",mensaje.Message , "\n")
    if err := stream.Send(mensaje); err != nil {
        log.Fatalf("Error al enviar mensaje: %v", err)
    }

}
//Sacar id y nombre de la persona del mensaje
func ObtenerIDNombre(mensaje string) (int, string){
    // Dividir el mensaje en espacios en blanco
    partes := strings.Fields(mensaje)

    if len(partes) >= 2 {
        // El primer elemento es el ID, y el resto es el nombre
        id,_ := strconv.Atoi(partes[0])
        nombre := strings.Join(partes[1:], " ")
        return id, nombre
    }
	return 0, ""
}

func main(){

	//Coneccion con el servidor.
	//SE DEBE DIFERENCIAR LOS DATANODES POR EL PUERTO 50052 SERA NODO 1 Y 50053 NODO 2
	listener, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	archivo ,_  = os.Create("DATA.txt")
	serv := grpc.NewServer()
	personas = make(map[int]string)
	pb.RegisterOMSServer(serv, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
