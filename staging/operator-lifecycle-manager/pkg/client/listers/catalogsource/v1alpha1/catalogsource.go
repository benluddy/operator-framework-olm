/*
Copyright 2018 CoreOS, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was automatically generated by lister-gen

package v1alpha1

import (
	v1alpha1 "github.com/coreos-inc/alm/pkg/apis/catalogsource/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CatalogSourceLister helps list CatalogSources.
type CatalogSourceLister interface {
	// List lists all CatalogSources in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.CatalogSource, err error)
	// CatalogSources returns an object that can list and get CatalogSources.
	CatalogSources(namespace string) CatalogSourceNamespaceLister
	CatalogSourceListerExpansion
}

// catalogSourceLister implements the CatalogSourceLister interface.
type catalogSourceLister struct {
	indexer cache.Indexer
}

// NewCatalogSourceLister returns a new CatalogSourceLister.
func NewCatalogSourceLister(indexer cache.Indexer) CatalogSourceLister {
	return &catalogSourceLister{indexer: indexer}
}

// List lists all CatalogSources in the indexer.
func (s *catalogSourceLister) List(selector labels.Selector) (ret []*v1alpha1.CatalogSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CatalogSource))
	})
	return ret, err
}

// CatalogSources returns an object that can list and get CatalogSources.
func (s *catalogSourceLister) CatalogSources(namespace string) CatalogSourceNamespaceLister {
	return catalogSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CatalogSourceNamespaceLister helps list and get CatalogSources.
type CatalogSourceNamespaceLister interface {
	// List lists all CatalogSources in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.CatalogSource, err error)
	// Get retrieves the CatalogSource from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.CatalogSource, error)
	CatalogSourceNamespaceListerExpansion
}

// catalogSourceNamespaceLister implements the CatalogSourceNamespaceLister
// interface.
type catalogSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CatalogSources in the indexer for a given namespace.
func (s catalogSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.CatalogSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CatalogSource))
	})
	return ret, err
}

// Get retrieves the CatalogSource from the indexer for a given namespace and name.
func (s catalogSourceNamespaceLister) Get(name string) (*v1alpha1.CatalogSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("catalogsource-v1"), name)
	}
	return obj.(*v1alpha1.CatalogSource), nil
}
