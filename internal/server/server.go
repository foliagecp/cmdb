// Copyright 2022 Listware

package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	"git.fg-tech.ru/listware/cmdb/internal/cmdb/edge"
	"git.fg-tech.ru/listware/cmdb/internal/cmdb/finder"
	"git.fg-tech.ru/listware/cmdb/internal/cmdb/qdsl"
	"git.fg-tech.ru/listware/cmdb/internal/cmdb/vertex"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbfinder"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
	"google.golang.org/grpc"
)

// set max 100 MB
const maxMsgSize = 100 * 1024 * 1024

var (
	cmdbAddr = "127.0.0.1"
	cmdbPort = "31415"
)

func init() {
	if value, ok := os.LookupEnv("CMDB_ADDR"); ok {
		cmdbAddr = value
	}
	if value, ok := os.LookupEnv("CMDB_PORT"); ok {
		cmdbPort = value
	}
}

func New() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	go func() {
		for {
			select {
			case sig := <-sigChan:
				switch sig {
				case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
					cancel()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := serve(ctx); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func serve(ctx context.Context) (err error) {

	if err = arangodb.Bootstrap(ctx); err != nil {
		return
	}

	pc, err := net.Listen("tcp", fmt.Sprintf(":%s", cmdbPort))
	if err != nil {
		return
	}
	defer pc.Close()

	server := grpc.NewServer(
		grpc.MaxMsgSize(maxMsgSize),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	defer server.Stop()

	qdsl, err := qdsl.New(ctx)
	if err != nil {
		return
	}
	pbqdsl.RegisterQdslServiceServer(server, pbqdsl.QdslServiceServer(qdsl))

	finder, err := finder.New(ctx)
	if err != nil {
		return
	}
	pbfinder.RegisterFinderServiceServer(server, pbfinder.FinderServiceServer(finder))

	edge, err := edge.New(ctx)
	if err != nil {
		return
	}

	pbcmdb.RegisterEdgeServiceServer(server, pbcmdb.EdgeServiceServer(edge))

	vertex, err := vertex.New(ctx)
	if err != nil {
		return
	}

	pbcmdb.RegisterVertexServiceServer(server, pbcmdb.VertexServiceServer(vertex))

	go server.Serve(pc)

	<-ctx.Done()

	return
}

func Client() (conn *grpc.ClientConn, err error) {

	return grpc.Dial(fmt.Sprintf("%s:%s", cmdbAddr, cmdbPort), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
}
