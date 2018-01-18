// gabe/clients/fred-client.go

package fredClient

import (
	"context"
	"log"

	pb "github.com/dicurrio/protorepo/fred"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func getIndex(name string) string {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(":: GABE :: FRED-CLIENT :: Failed to connect: %v", err)
	}
	defer conn.Close()

	fredClient := pb.NewFredClient(conn)
	res, err := fredClient.GetIndex(context.Background(), &pb.Request{Name: name})
	if err != nil {
		log.Fatalf(":: GABE :: FRED-CLIENT :: GetIndex error: %v", err)
	}
	return res.Message
}
