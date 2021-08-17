package main;

import(
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/mkurban/lssolstore/lss"
)

// ToDo: Change this to have "service:port" be set with env viarables instead
const (
	port = "localhost:54351"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	lss.RegisterLotsizeSolutionStoreServer(s, lss.NewLSSServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
