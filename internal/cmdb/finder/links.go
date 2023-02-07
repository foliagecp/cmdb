// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package finder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb/query"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbfinder"
)

const (
	mask       = "FOR t IN %s FILTER %s RETURN t"
	collection = "links"
)

func (s *Server) Links(ctx context.Context, request *pbfinder.Request) (response *pbfinder.Response, err error) {
	response = &pbfinder.Response{}

	vars := make(map[string]any)
	var args []string

	if request.From != "" {
		args = append(args, "t._from == @from")
		vars["from"] = request.From
	}
	if request.To != "" {
		args = append(args, "t._to == @to")
		vars["to"] = request.To
	}
	if request.Name != "" {
		args = append(args, "t._name == @name")
		vars["name"] = request.Name
	}
	filter := strings.Join(args, " && ")
	q := fmt.Sprintf(mask, collection, filter)
	metas, resp, err := query.Query(ctx, s.client, q, vars)
	if err != nil {
		return
	}

	for i, meta := range metas {
		r := &pbcmdb.Response{}

		r.Meta = &pbcmdb.Meta{
			Key: meta.Key,
			Id:  meta.ID.String(),
			Rev: meta.Rev,
		}

		r.Payload, err = json.Marshal(resp[i])
		if err != nil {
			return
		}
		response.Links = append(response.Links, r)

	}

	return
}
