// Copyright 2022 Listware

package arangodb

import (
	"context"

	driver "github.com/arangodb/go-driver"
)

const (
	systemCollection  = "system"
	typesCollection   = "types"
	objectsCollection = "objects"
	linksCollection   = "links"
)

var (
	allowUserKeysPtr = true
)

func Bootstrap(ctx context.Context) (err error) {
	client, err := Connect()
	if err != nil {
		return
	}

	ok, err := client.DatabaseExists(ctx, cmdbName)
	if err != nil {
		return
	}

	if !ok {
		options := &driver.CreateDatabaseOptions{}
		if _, err = client.CreateDatabase(ctx, cmdbName, options); err != nil {
			return
		}
	}

	db, err := client.Database(ctx, cmdbName)
	if err != nil {
		return
	}

	// system collection
	if ok, err = db.CollectionExists(ctx, systemCollection); err != nil {
		return
	}

	if !ok {
		options := &driver.CreateCollectionOptions{
			IsSystem: true,
			KeyOptions: &driver.CollectionKeyOptions{
				AllowUserKeysPtr: &allowUserKeysPtr,
			},
		}
		if _, err = db.CreateCollection(ctx, systemCollection, options); err != nil {
			return
		}
	}

	// types collection
	if ok, err = db.CollectionExists(ctx, typesCollection); err != nil {
		return
	}

	if !ok {
		options := &driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{
				AllowUserKeysPtr: &allowUserKeysPtr,
			},
		}
		if _, err = db.CreateCollection(ctx, typesCollection, options); err != nil {
			return
		}
	}

	// objects collection
	if ok, err = db.CollectionExists(ctx, objectsCollection); err != nil {
		return
	}

	if !ok {
		options := &driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{
				Type: driver.KeyGeneratorType("uuid")},
		}
		if _, err = db.CreateCollection(ctx, objectsCollection, options); err != nil {
			return
		}
	}

	// links collection
	if ok, err = db.CollectionExists(ctx, linksCollection); err != nil {
		return
	}

	if !ok {
		options := &driver.CreateCollectionOptions{
			Type: driver.CollectionTypeEdge,
		}
		if _, err = db.CreateCollection(ctx, linksCollection, options); err != nil {
			return
		}
	}

	// system graph
	if ok, err = db.GraphExists(ctx, systemGraphName); err != nil {
		return
	}

	if !ok {
		options := &driver.CreateGraphOptions{
			EdgeDefinitions: []driver.EdgeDefinition{
				driver.EdgeDefinition{
					Collection: linksCollection,
					From:       []string{systemCollection, typesCollection, objectsCollection},
					To:         []string{typesCollection, objectsCollection},
				},
			},
		}
		if _, err = db.CreateGraphV2(ctx, systemGraphName, options); err != nil {
			return
		}
	}
	return
}
