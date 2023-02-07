// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package types

import (
	"context"
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex"
)

const (
	collection = "types"
)

func Create(ctx context.Context, payload any) (meta *documents.BaseDocument, err error) {
	response, err := vertex.Create(ctx, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Read(ctx context.Context, key string) (meta *documents.BaseDocument, payload json.RawMessage, err error) {
	response, err := vertex.Read(ctx, key, collection)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	payload = response.GetPayload()
	return
}

func Update(ctx context.Context, key string, payload any) (meta *documents.BaseDocument, err error) {
	response, err := vertex.Update(ctx, key, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Replace(ctx context.Context, key string, payload any) (meta *documents.BaseDocument, err error) {
	response, err := vertex.Replace(ctx, key, collection, payload)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}

func Remove(ctx context.Context, key string) (meta *documents.BaseDocument, err error) {
	response, err := vertex.Remove(ctx, key, collection)
	if err != nil {
		return
	}
	meta = documents.NewBaseDocument(response.Meta)
	return
}
