package parser

import (
	"github.com/aquasecurity/trivy-config-parsers/types"
)

// SourceFormat ...
type SourceFormat string

// YamlSourceFormat ...
const (
	YamlSourceFormat SourceFormat = "yaml"
	JsonSourceFormat SourceFormat = "json"
)

// FileContexts ...
type FileContexts []*FileContext

// FileContext ...
type FileContext struct {
	filepath     string
	lines        []string
	SourceFormat SourceFormat
	Parameters   map[string]*Parameter  `json:"Parameters" yaml:"Parameters"`
	Resources    map[string]*Resource   `json:"Resources" yaml:"Resources"`
	Globals      map[string]*Resource   `json:"Globals" yaml:"Globals"`
	Mappings     map[string]interface{} `json:"Mappings,omitempty" yaml:"Mappings"`
}

func (t *FileContext) GetResourceByLogicalID(name string) *Resource {
	for n, r := range t.Resources {
		if name == n {
			return r
		}
	}
	return nil
}

// GetResourceByType ...
func (t *FileContext) GetResourceByType(names ...string) []*Resource {
	var resources []*Resource
	for _, r := range t.Resources {
		for _, name := range names {
			if name == r.Type() {
				//
				resources = append(resources, r)
			}
		}
	}
	return resources
}

// Metadata ...
func (t *FileContext) Metadata() types.Metadata {
	rng := types.NewRange(t.filepath, 1, len(t.lines))

	return types.NewMetadata(rng, NewCFReference("Template", rng))
}