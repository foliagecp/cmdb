// Copyright 2023 NJWS, Inc.

package edge

import (
	"context"
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb/edge"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

func (s *Server) Replace(ctx context.Context, request *pbcmdb.Request) (response *pbcmdb.Response, err error) {
	response = &pbcmdb.Response{}
	meta, resp, err := edge.Replace(ctx, s.client, request.GetCollection(), request.GetKey(), request.GetPayload())
	if err != nil {
		return
	}
	response.Meta = &pbcmdb.Meta{
		Key: meta.Key,
		Id:  meta.ID.String(),
		Rev: meta.Rev,
	}
	response.Payload, err = json.Marshal(resp)
	return
}
