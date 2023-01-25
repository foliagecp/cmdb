// Copyright 2022 Listware

package qdsl

import (
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

// OptionsOption query options
type OptionsOption func(*pbqdsl.Options)

// NewOptions return new query options
func NewOptions(opts ...OptionsOption) *pbqdsl.Options {
	h := &pbqdsl.Options{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// WithId return vertex '_id', default true
func WithId() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Id = true
	}
}

// WithKey return vertex '_key', default false
func WithKey() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Key = true
	}
}

// WithName return edge '_name', default false
func WithName() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Name = true
	}
}

// WithObject return vertex 'object', default false
func WithObject() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Object = true
	}
}

// WithLink return edge 'object', default false
func WithLink() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Link = true
	}
}

// WithLinkId return edge '_id', default false
func WithLinkId() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.LinkId = true
	}
}

// WithType return edge '_type', default false
func WithType() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Type = true
	}
}

// FIXME WithType now disabled
func WithPath() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Path = true
	}
}

// WithRemove remove all founded results
func WithRemove() OptionsOption {
	return func(h *pbqdsl.Options) {
		h.Remove = true
	}
}
