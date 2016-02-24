package model

import "github.com/rancher/go-rancher/client"

type Node struct {
	client.Resource
	Common

	AgentState      string                 `json:"agentState,omitempty" yaml:"agent_state,omitempty"`
	Info            interface{}            `json:"info,omitempty" yaml:"info,omitempty"`
	Hostname        string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Labels          map[string]interface{} `json:"labels,omitempty" yaml:"labels,omitempty"`
	PublicEndpoints []interface{}          `json:"publicEndpoints,omitempty" yaml:"public_endpoints,omitempty"`
	ContainerIds    []ID                   `json:"containerIds"`
}

type NodeList struct {
	client.Collection
	Data []Node `json:"data,omitempty"`
}
