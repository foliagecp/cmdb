// Copyright 2022 Listware

package qdsl

import (
	"testing"

	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

func TestQdslToAql(t *testing.T) {
	query := &pbqdsl.Query{
		Query:   "*[?$._from == '47e98408-3d47-4730-ba94-c2314ce1982e'?]",
		Options: &pbqdsl.Options{},
	}
	elements, err := parse(query)
	if err != nil {
		t.Fatal(err)
	}

	for _, element := range elements {
		t.Log(element.Query)
	}
}
