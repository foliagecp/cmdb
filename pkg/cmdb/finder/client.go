// Copyright 2023 NJWS, INC

package finder

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/internal/server"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbfinder"
)

// Links try search links with parameters
// pbcmdb.Response.Payload will be decoded to 'cmdb/documents.EdgeDocument'
func Links(ctx context.Context, from, to, name, linkType string) (resp []*pbcmdb.Response, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pbfinder.NewFinderServiceClient(conn)

	r, err := client.Links(ctx, &pbfinder.Request{From: from, To: to, Name: name, Type: linkType})
	if err != nil {
		return
	}

	resp = r.GetLinks()
	return
}
