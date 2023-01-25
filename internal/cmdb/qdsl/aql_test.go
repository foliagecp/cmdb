// Copyright 2022 Listware

package qdsl

import (
	"testing"

	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

func TestPathToAql(t *testing.T) {
	query := &pbqdsl.Query{
		Query: "*[?$._name == 'init'?].exmt.functions.objects",
		Options: &pbqdsl.Options{
			LinkId: true,
			Object: true,
		},
	}
	elements, err := parse(query)
	if err != nil {
		t.Fatal(err)
	}

	for _, element := range elements {
		t.Log(element.Query)
	}
}
