package main

import (
	"context"
	"log"
	"os"

	pbBroker "github.com/felipefferrada/lab5/proto"
	pbFulcrum "github.com/felipefferrada/lab5/proto"
	"google.golang.org/grpc"
)

func enviarComandoAlBroker(comando string) (string, error) {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pbBroker.NewBrokerClient(conn)
	resp, err := client.RedirigirComando(context.Background(), &pbBroker.ComandoRequest{Comando: comando})
	if err != nil {
		return "", err
	}
	return resp.DireccionFulcrum, nil
}

func conectarYEnviarAlFulcrum(direccion string, comando string) error {
	conn, err := grpc.Dial(direccion, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbFulcrum.NewFulcrumClient(conn)
	// Aquí se puede ajustar para enviar comandos específicos al Fulcrum
	_, err = client.AgregarBase(context.Background(), &pbFulcrum.AgregarBaseRequest{
		Sector:   "Sector1",
		Base:     "Base1",
		Enemigos: 100,
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Uso: %s <comando>", os.Args[0])
	}
	comando := os.Args[1]

	direccion, err := enviarComandoAlBroker(comando)
	if err != nil {
		log.Fatalf("Error enviando comando al broker: %v", err)
	}

	err = conectarYEnviarAlFulcrum(direccion, comando)
	if err != nil {
		log.Fatalf("Error conectando y enviando comando al Fulcrum: %v", err)
	}

	log.Printf("Comando %s enviado exitosamente al servidor Fulcrum en %s", comando, direccion)
}
