// Copyright 2023 NJWS, Inc.

package vertex

import (
	"context"
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	driver "github.com/arangodb/go-driver"
)

func Replace(ctx context.Context, client driver.Client, name, key string, payload any) (meta driver.DocumentMeta, resp map[string]any, err error) {
	graph, err := arangodb.Graph(ctx, client)
	if err != nil {
		return
	}
	collection, err := graph.VertexCollection(ctx, name)
	if err != nil {
		return
	}
	ctx = driver.WithReturnNew(ctx, &resp)
	if b, ok := payload.([]byte); ok {
		var req map[string]any
		if err = json.Unmarshal(b, &req); err != nil {
			return
		}
		meta, err = collection.ReplaceDocument(ctx, key, req)
		return
	}
	meta, err = collection.ReplaceDocument(ctx, key, payload)
	return
}
