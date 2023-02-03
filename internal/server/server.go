// Copyright 2023 NJWS Inc.
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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// set max 100 MB
const maxMsgSize = 100 * 1024 * 1024

var (
	cmdbAddr = "127.0.0.1"
	cmdbPort = "31415"

	log = logrus.New()
)

func init() {
	if value, ok := os.LookupEnv("CMDB_ADDR"); ok {
		cmdbAddr = value
	}
	if value, ok := os.LookupEnv("CMDB_PORT"); ok {
		cmdbPort = value
	}
}

func Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)

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

	return serve(ctx)
}

func serve(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Info("bootstrap arangodb")

	if err = arangodb.Bootstrap(ctx); err != nil {
		return
	}
	port := fmt.Sprintf(":%s", cmdbPort)
	pc, err := net.Listen("tcp4", port)
	if err != nil {
		return
	}
	defer pc.Close()

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	defer server.Stop()

	log.Info("executing qdsl")
	qdsl, err := qdsl.New(ctx)
	if err != nil {
		return
	}
	pbqdsl.RegisterQdslServiceServer(server, pbqdsl.QdslServiceServer(qdsl))

	log.Info("executing finder")
	finder, err := finder.New(ctx)
	if err != nil {
		return
	}
	pbfinder.RegisterFinderServiceServer(server, pbfinder.FinderServiceServer(finder))

	log.Info("executing edge")
	edge, err := edge.New(ctx)
	if err != nil {
		return
	}

	pbcmdb.RegisterEdgeServiceServer(server, pbcmdb.EdgeServiceServer(edge))

	log.Info("executing vertex")
	vertex, err := vertex.New(ctx)
	if err != nil {
		return
	}

	pbcmdb.RegisterVertexServiceServer(server, pbcmdb.VertexServiceServer(vertex))

	log.Info("serving ", port)

	go func() {
		if err = server.Serve(pc); err != nil {
			cancel()
		}

	}()

	<-ctx.Done()

	return
}

func Client() (conn *grpc.ClientConn, err error) {
	return grpc.Dial(fmt.Sprintf("%s:%s", cmdbAddr, cmdbPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)))
}
