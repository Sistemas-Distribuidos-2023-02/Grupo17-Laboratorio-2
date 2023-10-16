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
)

var personas map[int]string
var archivo *os.File

type server struct {
	pb.UnimplementedOMSServer
}

func (s *server) NotifyBidirectional(steam pb.OMS_NotifyBidirectionalServer) error {
	for {
		request, _ := steam.Recv()
		log.Printf("Mensaje recibido: %s", request.Message)
		id, nombre := ObtenerIDNombre(request.Message)
		personas[id] = nombre
		fmt.Println(personas[id])
		archivo.WriteString( fmt.Sprint(id) + " " + nombre + "\n")
		fmt.Println(personas)
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
