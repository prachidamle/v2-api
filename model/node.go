package model

import "github.com/rancher/go-rancher/client"

type NodeCommon struct {
	Common
	AgentState      string                 `json:"agentState,omitempty" yaml:"agent_state,omitempty" db:"agent_state"`
	Hostname        string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Info            interface{}            `json:"info,omitempty" yaml:"info,omitempty"`
	Labels          map[string]interface{} `json:"labels,omitempty" yaml:"labels,omitempty"`
	PhysicalHostId  string                 `json:"physicalHostId,omitempty" yaml:"physical_host_id,omitempty" db:"physical_host_id"`
	PublicEndpoints []interface{}          `json:"publicEndpoints,omitempty" yaml:"public_endpoints,omitempty"`
}

type NodeV1 struct {
	client.Resource
	NodeCommon
	AgentId      string `json:"agentId,omitempty" yaml:"agent_id,omitempty"`
	ApiProxy     string `json:"apiProxy,omitempty" yaml:"api_proxy,omitempty"`
	ComputeTotal int64  `json:"computeTotal,omitempty" yaml:"compute_total,omitempty"`
	Kind         string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

type NodeV2 struct {
	client.Resource
	NodeCommon
	InstanceIds []ID `json:"instanceIds"`
}

type NodeList struct {
	client.Collection
	Data []NodeV2 `json:"data,omitempty"`
}

func getNodeSchema(schemas *client.Schemas) {
	node := AddType(schemas, "node", NodeV2{}, (*client.HostOperations)(nil))
	node.CollectionMethods = []string{"GET", "POST"}
	V2ToV1NamesMap["node"] = "host"	
}
