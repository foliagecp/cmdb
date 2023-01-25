// Copyright 2022 Listware

package finder

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbfinder"
	driver "github.com/arangodb/go-driver"
)

type Server struct {
	pbfinder.UnimplementedFinderServiceServer

	client driver.Client
}

func New(ctx context.Context) (s *Server, err error) {
	s = &Server{}
	s.client, err = arangodb.Connect()
	return
}
