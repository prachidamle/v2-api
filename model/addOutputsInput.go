package model

import (
	"github.com/rancher/go-rancher/client"
)	

type AddOutputsInput struct {
	client.Resource
	Outputs map[string]interface{} `json:"outputs,omitempty" yaml:"outputs,omitempty" schema:",type=map[string],required=true,create=true,update=false"`
}


func getAddOutputsInputSchema(schemas *client.Schemas) {
	addOutputs := AddType(schemas,"addOutputsInput", AddOutputsInput{}, (*client.AddOutputsInput)(nil))

	V2ToV1NamesMap["addOutputsInput"] = "addOutputsInput"
	
	addOutputs.CollectionMethods = []string{"GET", "POST"}
	addOutputs.ResourceMethods = []string{"GET", "POST"}
}
