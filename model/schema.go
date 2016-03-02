package model

import (
	"github.com/rancher/go-rancher/client"
)




func NewSchema() *client.Schemas {
	schemas := &client.Schemas{}

	ResourceActionsToStates = make(map[string][]string)	
	V2ToV1NamesMap = make(map[string]string)


	apiVersion := schemas.AddType("apiVersion", client.Resource{})
	apiVersion.CollectionMethods = []string{}
	schemas.AddType("apiVersion", client.Resource{})
	schemas.AddType("schema", client.Schema{})
	schemas.AddType("service", Service{})

	getContainerSchema(schemas)
	getStackSchema(schemas)
	getNodeSchema(schemas)

	restartPolicy := schemas.AddType("restartPolicy", client.RestartPolicy{})
	restartPolicy.CollectionMethods = []string{}
	logConfig := schemas.AddType("logConfig", client.LogConfig{})
	logConfig.CollectionMethods = []string{}
	healthCheck := schemas.AddType("instanceHealthCheck", client.InstanceHealthCheck{})
	healthCheck.CollectionMethods = []string{}	

	getAddOutputsInputSchema(schemas)

	return schemas
}

