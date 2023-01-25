// Copyright 2022 Listware

package qdsl

import (
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

func normalizeOptions(options *pbqdsl.Options) {
	if !options.GetId() &&
		!options.GetKey() &&
		!options.GetName() &&
		!options.GetType() &&
		!options.GetLink() &&
		!options.GetLinkId() &&
		!options.GetPath() &&
		!options.GetObject() {
		options.Id = true
	}
}
func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
