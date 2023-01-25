// Copyright 2022 Listware

package qdsl

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/server"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

// RawQdsl query with options as object
func RawQdsl(ctx context.Context, query string, options *pbqdsl.Options) (nodes documents.Nodes, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()
	client := pbqdsl.NewQdslServiceClient(conn)

	elements, err := client.Qdsl(ctx, &pbqdsl.Query{Query: query, Options: options})
	if err != nil {
		return
	}

	nodes = documents.NewNodes()

	for _, element := range elements.GetElements() {
		nodes.AddElement(element)
	}
	return
}

// Qdsl query with options OptionsOption array
func Qdsl(ctx context.Context, query string, options ...OptionsOption) (documents.Nodes, error) {
	return RawQdsl(ctx, query, NewOptions(options...))
}
