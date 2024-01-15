package config

import (
	"fmt"
	"strings"

	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/istio-mcp/pkg/config/schema/resource"
	"istio.io/istio-mcp/pkg/model"
)

func McpToPilot(rev string, m *mcp.Resource, gvkArr []string) (*model.Config, error) {
	if m == nil || m.Metadata == nil {
		return &model.Config{}, nil
	}

	c := &model.Config{
		Meta: model.Meta{
			ResourceVersion: m.Metadata.Version,
			Labels:          m.Metadata.Labels,
			Annotations:     m.Metadata.Annotations,
		},
	}
	gvk := resource.GroupVersionKind{Group: gvkArr[0], Version: gvkArr[1], Kind: gvkArr[2]}
	c.GroupVersionKind = gvk

	if rev != "" && !model.ObjectInRevision(c, rev) { // In case upstream does not support rev in node meta.
		// empty rev will be considered as legacy client and NEEDS ALL CONFIGS
		return nil, nil
	}

	nsn := strings.Split(m.Metadata.Name, "/")
	if len(nsn) != 2 {
		return nil, fmt.Errorf("invalid name %s", m.Metadata.Name)
	}
	c.Namespace = nsn[0]
	c.Name = nsn[1]
	c.CreationTimestamp = m.Metadata.CreateTime.AsTime()

	if m.Body != nil {
		pb, err := m.Body.UnmarshalNew()
		if err != nil {
			return nil, err
		}
		c.Spec = pb

	}
	return c, nil
}
