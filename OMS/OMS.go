package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	pb "main/proto"
	"fmt"
	"strings"
	"os"
	"context"
)

var archivo *os.File
var id int


type server struct {
	pb.UnimplementedOMSServer
}

func (s *server) NotifyBidirectional(steam pb.OMS_NotifyBidirectionalServer) error {
	for {
		request, err := steam.Recv()
		if err != nil {
			break
		}
		log.Printf("Mensaje recibido: %s", request.Message)
		//Casi seguro que es un mensaje de la ONU
		if len(request.Message) <= 10 {
			fmt.Println("La ONU a preguntado por los datos de: " + request.Message)
			continue
		}
		inicialApellido := RevisarInicial(request.Message)
		estado,nombre := ObtenerEstado(request.Message)

		// Se decide a que nodo se debe enviar el mensaje.
		//Se guarda en el archivo DATA.txt el ID  NODO ESTADO 
		if inicialApellido <= 77{
			archivo.WriteString( fmt.Sprint(id) + " Nodo 1 " +  estado + "\n")
			MandarDataDatanodes(id, 1, nombre)
		}else {
			archivo.WriteString( fmt.Sprint(id) + " Nodo 2  " + estado + "\n")
		}
		id+=1
		
	}
	
	return nil
}

func MandarDataDatanodes(id , nodo int, nombre string ){
	msg := fmt.Sprint(id) + " " + nombre
	nodo+=1
	conn, err := grpc.Dial("localhost:5005"+fmt.Sprint(nodo), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
    }
	client := pb.NewOMSClient(conn)
    stream, err := client.NotifyBidirectional(context.Background())
	if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
    }
	mensaje := &pb.Request{Message: msg}
    if err := stream.Send(mensaje); err != nil {
        log.Fatalf("Error al enviar mensaje: %v", err)
    }
}

func RevisarInicial(Persona string) byte{
	lineas := strings.Split(Persona, " ")
	apellido := lineas[1]
	inicialApellido := apellido[0]
	return inicialApellido
}

func ObtenerEstado(Persona string)(string , string){
	lineas := strings.Split(Persona, "\n")
	return lineas[1], lineas[0]
}

func main(){
	//Coneccion con el servidor.
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := grpc.NewServer()
	id = 1
	pb.RegisterOMSServer(serv, &server{})
    archivo ,_  = os.Create("DATA.txt")
	log.Printf("server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}