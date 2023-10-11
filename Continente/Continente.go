package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"math/rand"
	"time"
	"log"
	"google.golang.org/grpc"
	"context"
	pb "main/proto"
	
)


func InfectadoMuerto(Persona string ) string  {
	numeroAleatorio := rand.Float64()
    umbralInfectada := 0.55
	var resultado string
    if numeroAleatorio < umbralInfectada {
        resultado = "Infectada"
    }else {
        resultado = "Muerta"
	}
	return Persona + resultado
}

func MandarDataOMS(Persona string, conn *grpc.ClientConn, Flag bool) string {
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
	if Flag == true {
		time.Sleep(50*time.Millisecond)
	}else{
		time.Sleep(3*time.Second)
	}
	
	return Persona
}


func main() {
	rand.Seed(time.Now().UnixNano())
	archivo, _ := os.Open("names.txt")
	// Crea un lector de bufio para el archivo
    lector := bufio.NewReader(archivo)
	Personas := []string{}
	
    // Lee las líneas del archivo
    for {
        linea, err := lector.ReadString('\n')
        if err != nil {
            break // Fin del archivo
        }
        linea = strings.TrimSuffix(linea, " ")

        Personas = append(Personas, linea)
    }
	archivo.Close()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
    }
    defer conn.Close()
	total := len(Personas)
	Flag := true
	for i := 0; i < total; i++ {
		if i == 4 {
			Flag = false
		}
		// Elegir un índice al azar en la lista
		indiceAleatorio := rand.Intn(len(Personas))
		// Guardar y eliminar el elemento seleccionado en la lista
		Resultado  := InfectadoMuerto(Personas[indiceAleatorio])
		MandarDataOMS(Resultado, conn, Flag)
		Personas = append(Personas[:indiceAleatorio], Personas[indiceAleatorio+1:]...)
	}
	

}
