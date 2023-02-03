// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package edge

import (
	"context"
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	driver "github.com/arangodb/go-driver"
)

func Create(ctx context.Context, client driver.Client, name string, payload any) (meta driver.DocumentMeta, resp map[string]any, err error) {
	graph, err := arangodb.Graph(ctx, client)
	if err != nil {
		return
	}

	collection, _, err := graph.EdgeCollection(ctx, name)
	if err != nil {
		return
	}

	ctx = driver.WithReturnNew(ctx, &resp)

	if b, ok := payload.([]byte); ok {
		var req map[string]any
		if err = json.Unmarshal(b, &req); err != nil {
			return
		}
		meta, err = collection.CreateDocument(ctx, req)
		return
	}
	meta, err = collection.CreateDocument(ctx, payload)
	return
}
