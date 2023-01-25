// Copyright 2022 Listware

package query

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	driver "github.com/arangodb/go-driver"
)

func Query(ctx context.Context, client driver.Client, query string, vars map[string]any) (metas []driver.DocumentMeta, resp []map[string]any, err error) {
	db, err := arangodb.Database(ctx, client)
	if err != nil {
		return
	}

	cursor, err := db.Query(ctx, query, vars)
	if err != nil {
		return
	}
	defer cursor.Close()

	for cursor.HasMore() {
		var obj map[string]any

		meta, err := cursor.ReadDocument(ctx, &obj)
		if err != nil {
			return nil, nil, err
		}
		resp = append(resp, obj)
		metas = append(metas, meta)
	}
	return
}
