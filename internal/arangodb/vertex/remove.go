// Copyright 2022 Listware

package vertex

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	driver "github.com/arangodb/go-driver"
)

func Remove(ctx context.Context, client driver.Client, name, key string) (meta driver.DocumentMeta, err error) {
	graph, err := arangodb.Graph(ctx, client)
	if err != nil {
		return
	}
	collection, err := graph.VertexCollection(ctx, name)
	if err != nil {
		return
	}
	return collection.RemoveDocument(ctx, key)
}
