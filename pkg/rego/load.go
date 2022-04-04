package rego

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
)

func isRegoFile(name string) bool {
	return strings.HasSuffix(name, bundle.RegoExt) && !strings.HasSuffix(name, "_test"+bundle.RegoExt)
}

func (s *Scanner) loadPoliciesFromDirs(target fs.FS, paths []string) (map[string]*ast.Module, error) {
	modules := make(map[string]*ast.Module)
	for _, path := range paths {
		if err := fs.WalkDir(target, path, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !isRegoFile(info.Name()) {
				return nil
			}
			data, err := fs.ReadFile(target, path)
			if err != nil {
				return err
			}
			module, err := ast.ParseModuleWithOpts(path, string(data), ast.ParserOptions{})
			if err != nil {
				return err
			}
			modules[path] = module
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return modules, nil
}

func (s *Scanner) LoadPolicies(loadEmbedded bool, srcFS fs.FS, paths ...string) error {

	if s.policies == nil {
		s.policies = make(map[string]*ast.Module)
	}

	if loadEmbedded {
		loaded, err := loadEmbeddedPolicies()
		if err != nil {
			return fmt.Errorf("failed to load embedded rego policies: %w", err)
		}
		for name, policy := range loaded {
			s.policies[name] = policy
		}
		s.debug("Loaded %d embedded policies.", len(loaded))
	}

	var err error
	if len(paths) > 0 {
		loaded, err := s.loadPoliciesFromDirs(srcFS, paths)
		if err != nil {
			return fmt.Errorf("failed to load rego policies from %s: %w", paths, err)
		}
		for name, policy := range loaded {
			s.policies[name] = policy
		}
		s.debug("Loaded %d policies from disk.", len(loaded))
	}

	// gather namespaces
	uniq := make(map[string]struct{})
	for _, module := range s.policies {
		namespace := getModuleNamespace(module)
		uniq[namespace] = struct{}{}
	}
	var namespaces []string
	for namespace := range uniq {
		namespaces = append(namespaces, namespace)
	}
	store, err := initStore(s.dataDirs, namespaces)
	if err != nil {
		return fmt.Errorf("unable to load data: %w", err)
	}
	s.store = store

	compiler := ast.NewCompiler()
	compiler.Compile(s.policies)
	if compiler.Failed() {
		return compiler.Errors
	}
	s.compiler = compiler

	s.retriever = NewMetadataRetriever(compiler)

	return nil
}