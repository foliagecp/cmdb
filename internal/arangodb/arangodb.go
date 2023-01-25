// Copyright 2022 Listware

package arangodb

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	driver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

const (
	cmdbName        = "CMDB"
	systemGraphName = "system"
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

func Connect() (client driver.Client, err error) {
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

	return driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangoUser, arangoPassword),
	})

}

func Database(ctx context.Context, client driver.Client) (driver.Database, error) {
	return client.Database(ctx, cmdbName)
}

func Graph(ctx context.Context, client driver.Client) (graph driver.Graph, err error) {
	db, err := client.Database(ctx, cmdbName)
	if err != nil {
		return
	}
	return db.Graph(ctx, systemGraphName)
}
