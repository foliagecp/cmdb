// Copyright 2022 Listware

package qdsl

import (
	"context"
	"sort"
	"strings"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

func parse(query *pbqdsl.Query) (elements []*Element, err error) {
	got, err := ParseReader("", strings.NewReader(query.GetQuery()))
	if err != nil {
		return
	}

	if query.GetOptions() == nil {
		query.Options = &pbqdsl.Options{}
	}

	normalizeOptions(query.GetOptions())

	for _, i := range got.([]any) {
		if element, ok := i.(*Element); ok {
			pathToAql(element, query.GetOptions())
			elements = append(elements, element)
		}
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Action > elements[j].Action
	})
	return
}

func (s *Server) query(ctx context.Context, element *Element) (nodes documents.Nodes, err error) {
	db, err := arangodb.Database(ctx, s.client)
	if err != nil {
		return
	}
	cursor, err := db.Query(ctx, element.Query, nil)
	if err != nil {
		return
	}
	defer cursor.Close()

	for cursor.HasMore() {
		var node documents.Node

		if _, err = cursor.ReadDocument(ctx, &node); err != nil {
			return
		}

		nodes.Add(&node)
	}

	return
}
