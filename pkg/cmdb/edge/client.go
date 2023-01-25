// Copyright 2022 Listware

package edge

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/server"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

func Create(ctx context.Context, collection string, payload any) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewEdgeServiceClient(conn)

	request := &pbcmdb.Request{Collection: collection}
	return client.Create(ctx, request)
}

func Read(ctx context.Context, key, collection string) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewEdgeServiceClient(conn)

	return client.Read(ctx, &pbcmdb.Request{Key: key, Collection: collection})
}

func Update(ctx context.Context, key, collection string, payload any) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewEdgeServiceClient(conn)

	request := &pbcmdb.Request{Key: key, Collection: collection}

	return client.Update(ctx, request)
}

func Remove(ctx context.Context, key, collection string) (resp *pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbcmdb.NewEdgeServiceClient(conn)

	request := &pbcmdb.Request{Key: key, Collection: collection}

	return client.Remove(ctx, request)
}
