// Copyright 2022 Listware

package edge

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	driver "github.com/arangodb/go-driver"
)

type Server struct {
	pbcmdb.UnimplementedEdgeServiceServer

	client driver.Client
}

func New(ctx context.Context) (s *Server, err error) {
	s = &Server{}
	s.client, err = arangodb.Connect()
	return
}
