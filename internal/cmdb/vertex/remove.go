// Copyright 2022 Listware

package vertex

import (
	"context"

	vertex "git.fg-tech.ru/listware/cmdb/internal/arangodb/vertex"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

func (s *Server) Remove(ctx context.Context, request *pbcmdb.Request) (response *pbcmdb.Response, err error) {
	response = &pbcmdb.Response{}
	meta, err := vertex.Remove(ctx, s.client, request.GetCollection(), request.GetKey())
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
