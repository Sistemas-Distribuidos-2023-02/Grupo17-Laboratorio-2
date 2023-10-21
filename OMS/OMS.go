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
	"bufio"
	"strconv"
	"sync"
)

var mu sync.Mutex
var archivo *os.File
var id int
var Total int

type server struct {
	pb.UnimplementedOMSServer
}
var DataEntregarONU []string

func (s *server) NotifyBidirectional(steam pb.OMS_NotifyBidirectionalServer) error {
		request, err := steam.Recv()
		if err != nil {
			return err
		}
		fmt.Println(request.Message)
		// 1 = Muertos /  2 = Infectados
		if request.Message == "1" || request.Message == "2" {
			fmt.Println("La ONU a preguntado por los datos de: " + request.Message)
			// Abre el archivo
			file, err := os.Open("Data.txt")
			if err != nil {
				fmt.Println("Error al abrir el archivo:", err)
				return err
			}
			// Lee el archivo
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				
				line := scanner.Text()
				// Lee la línea actual
				parts := strings.Split(line, " ")

				if len(parts) == 3 {
					id := parts[0]
					nodo := parts[1]
					estado := parts[2]

				if request.Message == "1" && estado == "Muerta" {
					Total +=1
					if nodo == "1" {
						fmt.Println("Se debe preguntar al nodo 1")
						id , _ := strconv.Atoi(id)
						go MandarDataDatanodes(id , 1, "Preguntar")

					}else {
						Total -=1
						fmt.Println("Se debe preguntar al nodo 2")
						//MandarDataDatanodes(0, 2, "Preguntar")

					}
				} else if request.Message == "2" && estado == "Infectada"{
					Total +=1
					fmt.Printf("ID: %s, NODO: %s, ESTADO: %s\n", id, nodo, estado)
					if nodo == "1" {
						fmt.Println("Se debe preguntar al nodo 1")
						id , _ := strconv.Atoi(id)
						MandarDataDatanodes(id, 1, "Preguntar")

					}else {
						fmt.Println("Se debe preguntar al nodo 2")
						//MandarDataDatanodes(0, 2, "Preguntar")
					}
				}
			}
			}	
			return nil
		}
		
		estado,nombre := ObtenerEstado(request.Message)
		if estado == "Data" {
			Total -=1
			fmt.Println("Los Datanodes nos diero el nombre de : " + nombre)
			DataEntregarONU = append(DataEntregarONU, nombre)
			fmt.Println(Total)
			if Total == 0 {
				fmt.Println("Le respondimos a la ONU")
				for i := 0; i < len(DataEntregarONU); i++ {
					MandarDataONU()
					DataEntregarONU = DataEntregarONU[1:]
				}
				DataEntregarONU = nil
			}
			return nil 
		}
		inicialApellido := RevisarInicial(request.Message)
		// Se decide a que nodo se debe enviar el mensaje.
		//Se guarda en el archivo DATA.txt el ID  NODO ESTADO 
		if inicialApellido <= 77{
			archivo.WriteString( fmt.Sprint(id) + " 1 " +  estado + "\n")
			MandarDataDatanodes(id, 1, nombre)
		}else {
			archivo.WriteString( fmt.Sprint(id) + " 2 " + estado + "\n")
			//MandarDataDatanodes(id, 2, nombre)
		}
		id+=1
		
	return nil
}

func MandarDataONU(){
	conn, err := grpc.Dial("localhost:50070", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
    }
	client := pb.NewOMSClient(conn)
    stream, err := client.NotifyBidirectional(context.Background())
	if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
    }
	mensaje := &pb.Request{Message: DataEntregarONU[0]}
	if err := stream.Send(mensaje); err != nil {
		log.Fatalf("Error al enviar mensaje: %v", err)
	}
}	

func MandarDataDatanodes(id , nodo int, nombre string ){
	msg := ""
	if nombre == "Preguntar" {
		msg = fmt.Sprint(id) + " P" 
	}else {
		msg = fmt.Sprint(id) + " " + nombre
	}
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
	Total = 0 
	id = 1
	pb.RegisterOMSServer(serv, &server{})
    archivo ,_  = os.Create("DATA.txt")
	log.Printf("server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}