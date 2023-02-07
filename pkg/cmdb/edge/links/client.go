// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package links

import (
	"context"
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/edge"
)

const (
	collection = "links"
)

func Create(ctx context.Context, payload any) (meta *documents.BaseDocument, err error) {
	response, err := edge.Create(ctx, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Read(ctx context.Context, key string) (payload json.RawMessage, err error) {
	response, err := edge.Read(ctx, key, collection)
	if err != nil {
		return
	}
	payload = response.GetPayload()
	return
}

func Update(ctx context.Context, key string, payload any) (meta *documents.BaseDocument, err error) {
	response, err := edge.Update(ctx, key, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Replace(ctx context.Context, key string, payload any) (meta *documents.BaseDocument, err error) {
	response, err := edge.Replace(ctx, key, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Remove(ctx context.Context, key string) (meta *documents.BaseDocument, err error) {
	response, err := edge.Remove(ctx, key, collection)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}
