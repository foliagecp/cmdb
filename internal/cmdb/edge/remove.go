// Copyright 2022 Listware

package edge

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb/edge"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

func (s *Server) Remove(ctx context.Context, request *pbcmdb.Request) (response *pbcmdb.Response, err error) {
	response = &pbcmdb.Response{}

	meta, err := edge.Remove(ctx, s.client, request.GetCollection(), request.GetKey())
	if err != nil {
		return
	}

	response.Meta = &pbcmdb.Meta{
		Key: meta.Key,
		Id:  meta.ID.String(),
		Rev: meta.Rev,
	}
	return
}
