// Created by interfacer; DO NOT EDIT

package api_if

import (
	"github.com/hashicorp/consul/api"
)

// Catalog is an interface generated for "github.com/hashicorp/consul/api".Catalog.
type Catalog interface {
	Datacenters() ([]string, error)
	Deregister(*api.CatalogDeregistration, *api.WriteOptions) (*api.WriteMeta, error)
	Node(string, *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error)
	Nodes(*api.QueryOptions) ([]*api.Node, *api.QueryMeta, error)
	Register(*api.CatalogRegistration, *api.WriteOptions) (*api.WriteMeta, error)
	Service(string, string, *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error)
	Services(*api.QueryOptions) (map[string][]string, *api.QueryMeta, error)
}
