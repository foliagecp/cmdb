// Copyright 2022 Listware

package qdsl

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"git.fg-tech.ru/listware/cmdb/internal/arangodb"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
	driver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

var (
	arangoAddr     string
	arangoUser     string
	arangoPassword string
)

func init() {
	if value, ok := os.LookupEnv("ARANGO_ADDR"); ok {
		arangoAddr = value
	}
	if value, ok := os.LookupEnv("ARANGO_USER"); ok {
		arangoUser = value
	}
	if value, ok := os.LookupEnv("ARANGO_PASSWORD"); ok {
		arangoPassword = value
	}
}

type Server struct {
	pbqdsl.UnimplementedQdslServiceServer

	client driver.Client
}

func New(ctx context.Context) (s *Server, err error) {
	s = &Server{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Open a client connection
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Transport: tr,
		Endpoints: []string{arangoAddr},
	})
	if err != nil {
		return
	}

	s.client, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangoUser, arangoPassword),
	})

	return
}

func (s *Server) Qdsl(ctx context.Context, query *pbqdsl.Query) (elements *pbqdsl.Elements, err error) {
	elements = &pbqdsl.Elements{Elements: make([]*pbqdsl.Element, 0)}

	qdslElements, err := parse(query)
	if err != nil {
		return
	}

	graph, err := arangodb.Graph(ctx, s.client)
	if err != nil {
		return
	}

	for _, element := range qdslElements {
		documents, err := s.query(ctx, element)
		if err != nil {
			return elements, err
		}

		for _, document := range documents {
			elements.Elements = append(elements.Elements, document.ToElement())

			if query.Options.Remove {

				if query.Options.Id {

					collection, err := graph.VertexCollection(ctx, document.Id.Collection())
					if err != nil {
						return elements, err
					}
					if _, err := collection.RemoveDocument(ctx, document.Id.Key()); err != nil {
						return elements, err
					}
				}

				if query.Options.LinkId {
					collection, _, err := graph.EdgeCollection(ctx, document.LinkId.Collection())
					if err != nil {
						return elements, err
					}
					if _, err := collection.RemoveDocument(ctx, document.LinkId.Key()); err != nil {
						return elements, err
					}
				}

			}
		}
	}

	return
}
