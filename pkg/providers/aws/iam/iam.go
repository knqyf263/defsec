package iam

import (
	"github.com/aquasecurity/defsec/internal/types"
	"github.com/liamg/iamgo"
)

type IAM struct {
	PasswordPolicy PasswordPolicy
	Policies       []Policy
	Groups         []Group
	Users          []User
	Roles          []Role
}

type Policy struct {
	types.Metadata
	Name     types.StringValue
	Document Document
}

type Document struct {
	types.Metadata
	Parsed   iamgo.Document
	IsOffset bool
	HasRefs  bool
}

type Group struct {
	types.Metadata
	Name     types.StringValue
	Users    []User
	Policies []Policy
}

type User struct {
	types.Metadata
	Name     types.StringValue
	Groups   []Group
	Policies []Policy
}

type Role struct {
	types.Metadata
	Name     types.StringValue
	Policies []Policy
}

func (d Document) MetadataFromIamGo(r ...iamgo.Range) types.Metadata {
	m := d.GetMetadata()
	if d.HasRefs {
		return m
	}
	newRange := m.Range()
	var start int
	if !d.IsOffset {
		start = newRange.GetStartLine()
	}
	for _, rng := range r {
		newRange := types.NewRange(
			newRange.GetLocalFilename(),
			start+rng.StartLine,
			start+rng.EndLine,
			newRange.GetSourcePrefix(),
			newRange.GetFS(),
		)
		m = types.NewMetadata(newRange, m.Reference()).WithParent(m)
	}
	return m
}
