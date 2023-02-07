// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package vertex

import (
	"context"
	"encoding/json"
	"errors"

	"git.fg-tech.ru/listware/cmdb/internal/server"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

var (
	ErrEmptyPayload = errors.New("empty payload")
)

func Create(ctx context.Context, collection string, payload any) (resp *pbcmdb.Response, err error) {
	if payload == nil {
		return nil, ErrEmptyPayload
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return
	}

	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewVertexServiceClient(conn)

	request := &pbcmdb.Request{Collection: collection, Payload: payloadRaw}
	return client.Create(ctx, request)
}

func Read(ctx context.Context, key, collection string) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewVertexServiceClient(conn)

	return client.Read(ctx, &pbcmdb.Request{Key: key, Collection: collection})
}

func Update(ctx context.Context, key, collection string, payload any) (resp *pbcmdb.Response, err error) {
	if payload == nil {
		return nil, ErrEmptyPayload
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return
	}

	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewVertexServiceClient(conn)

	request := &pbcmdb.Request{Key: key, Collection: collection, Payload: payloadRaw}

	return client.Update(ctx, request)
}

func Replace(ctx context.Context, key, collection string, payload any) (resp *pbcmdb.Response, err error) {
	if payload == nil {
		return nil, ErrEmptyPayload
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return
	}

	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewVertexServiceClient(conn)

	request := &pbcmdb.Request{Key: key, Collection: collection, Payload: payloadRaw}

	return client.Replace(ctx, request)
}

func Remove(ctx context.Context, key, collection string) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewVertexServiceClient(conn)

	request := &pbcmdb.Request{Key: key, Collection: collection}

	return client.Remove(ctx, request)
}
