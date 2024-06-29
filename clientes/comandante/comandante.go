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

func consultarEnemigos(direccion string, sector string, base string) (int32, error) {
	conn, err := grpc.Dial(direccion, grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pbFulcrum.NewFulcrumClient(conn)
	resp, err := client.GetEnemigos(context.Background(), &pbFulcrum.GetEnemigosRequest{Sector: sector, Base: base})
	if err != nil {
		return 0, err
	}
	return resp.Enemigos, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Uso: %s <sector> <base>", os.Args[0])
	}
	sector := os.Args[1]
	base := os.Args[2]

	direccion, err := enviarComandoAlBroker("consultar")
	if err != nil {
		log.Fatalf("Error enviando comando al broker: %v", err)
	}

	enemigos, err := consultarEnemigos(direccion, sector, base)
	if err != nil {
		log.Fatalf("Error consultando enemigos en el Fulcrum: %v", err)
	}

	log.Printf("Cantidad de enemigos en %s %s: %d", sector, base, enemigos)
}
