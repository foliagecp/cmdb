// Copyright 2022 Listware

package types

import (
	"git.fg-tech.ru/listware/cmdb/internal/schema"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

// Type is a object type struct
type Type struct {
	Schema   *schema.Schema                        `json:"schema"`
	Triggers map[string]map[string]*pbcmdb.Trigger `json:"triggers"`
}

// NewType return new object type
func NewType(schema *schema.Schema) Type {
	return Type{
		Schema:   schema,
		Triggers: make(map[string]map[string]*pbcmdb.Trigger),
	}
}

// ReflectType get reflected profiletype
func ReflectType(v interface{}) *Type {
	r := &schema.Reflector{
		AllowAdditionalProperties: true,
		ExpandedStruct:            true,
		TitleNotation:             "dash",
	}
	return &Type{
		Schema:   r.Reflect(v),
		Triggers: make(map[string]map[string]*pbcmdb.Trigger),
	}
}
