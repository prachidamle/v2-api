package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
)

type Container struct{}

func (s *Server) getContainersSQL(r *http.Request, id string) string {
	q := `
	  SELECT
	      COALESCE(name, '') as name, id, uuid, state, data
	  FROM instance
	  WHERE
	      account_id = :account_id
	      AND removed IS NULL
	      AND kind = 'container'`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}

func (s *Server) ContainerCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	v2 := &model.ContainerV2{}
	if err := s.parseInputParameters(r, v2); err != nil {
		return err
	}

	v1, err := containerFromV2(v2)
	if err != nil {
		return err
	}
	container, err := rancherClient.Container.Create(v1)

	if err != nil {
		return err
	}

	return s.getContainer(rw, r, container.Id)
}

func (s *Server) ContainerByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	
	action := r.URL.Query().Get("action")
	if action != "" && r.Method == "POST" {
		return s.ContainerAction(rw, r, action)
	}
	
	return s.getContainer(rw, r, vars["id"])
}

func (s *Server) ContainerList(rw http.ResponseWriter, r *http.Request) error {
	return s.getContainer(rw, r, "")
}

func (s *Server) getContainer(rw http.ResponseWriter, r *http.Request, id string) error {
	resourceType := "container"

	id = s.deobfuscate(r, resourceType, id)

	rows, err := s.namedQuery(s.getContainersSQL(r, id), map[string]interface{}{
		"account_id": s.getAccountID(r),
		"id":         id,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	response := &client.GenericCollection{
		Collection: client.Collection{
			Type:         "collection",
			ResourceType: resourceType,
		},
	}

	for rows.Next() {

		obj := &model.ContainerV1{}
		obj.Type = resourceType

		var data string

		if err := rows.Scan(&obj.Name, &obj.Id, &obj.UUID, &obj.State, &data); err != nil {
			return err
		}

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, resourceType, obj.Id)
		if err = s.parseData(data, obj); err != nil {
			return err
		}

		objV2, err := containerToV2(obj)
		if err != nil {
			return err
		}
		if id != "" {
			return s.writeResponse(rows.Err(), r, objV2)
		}
		response.Data = append(response.Data, objV2)
	}

	return s.writeResponse(rows.Err(), r, response)
}

func (s *Server) ContainerAction(rw http.ResponseWriter, r *http.Request, action string) error {
	vars := mux.Vars(r)
	containerObj := &model.ContainerV2{}
	containerObj.Resource.Type = "container"
	containerObj.ParentType = "instance"
	
	if s.isValidAction(action, &containerObj.Resource) {
		switch(action){
			
		 	case "allocate": return s.ContainerActionAllocate(rw, r)
		 	
		 	case "console": return s.ContainerActionConsole(rw, r)
		 	
		 	case "deallocate": return s.ContainerActionDeallocate(rw, r)
		 	
		 	case "execute": return s.ContainerActionExecute(rw, r)
		 	
		 	case "logs": return s.ContainerActionLogs(rw, r)
		 	
		 	case "migrate": return s.ContainerActionMigrate(rw, r)
		 	
		 	case "purge": return s.ContainerActionPurge(rw, r)
		 	
		 	case "remove": return s.ContainerActionRemove(rw, r)
		 	
		 	case "restart": return s.ContainerActionRestart(rw, r)
		 	
		 	case "restore": return s.ContainerActionRestore(rw, r)
		 	
		 	case "setlabels": return s.ContainerActionSetlabels(rw, r)
		 	
		 	case "start": return s.ContainerActionStart(rw, r)
		 	
		 	case "stop": return s.ContainerActionStop(rw, r)
		 	
		 	case "update": return s.ContainerActionUpdate(rw, r)
		 	
		 	case "updatehealthy": return s.ContainerActionUpdatehealthy(rw, r)
		 	
		 	case "updatereinitializing": return s.ContainerActionUpdatereinitializing(rw, r)
		 	
		 	case "updateunhealthy": return s.ContainerActionUpdateunhealthy(rw, r)
		}
	} else {
		return s.getContainer(rw, r, vars["id"])
	}
	
	return nil
}

func (s *Server) ContainerActionAllocate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionAllocate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionConsole(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.InstanceConsoleInput{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Container.ActionConsole(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionDeallocate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionDeallocate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionExecute(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.ContainerExec{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Container.ActionExecute(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionLogs(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.ContainerLogs{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Container.ActionLogs(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionMigrate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionMigrate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionPurge(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionPurge(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionRemove(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionRemove(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionRestart(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionRestart(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionRestore(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionRestore(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionSetlabels(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.SetLabelsInput{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Container.ActionSetlabels(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionStart(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionStart(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionStop(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.InstanceStop{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Container.ActionStop(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionUpdate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionUpdate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionUpdatehealthy(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionUpdatehealthy(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionUpdatereinitializing(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionUpdatereinitializing(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}

func (s *Server) ContainerActionUpdateunhealthy(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Container.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Container.ActionUpdateunhealthy(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getContainer(rw, r, updatedV1Obj.Id)
}



func containerToV2(v1 *model.ContainerV1) (*model.ContainerV2, error) {
	common := v1.ContainerCommon
	common.Transitioning = model.GetTransitioning(common.State, common.Transitioning)
	return &model.ContainerV2{
		Resource:        v1.Resource,
		ContainerCommon: common,
		Image:           v1.ImageUUID,
		WorkDir:         v1.WorkingDir,
		Logging:         v1.LogConfig,
		MemSwap:         v1.MemorySwap,
	}, nil
}

func containerFromV2(v2 *model.ContainerV2) (*client.Container, error) {
	b, err := json.Marshal(v2)
	if err != nil {
		return nil, err
	}
	v1 := &client.Container{}
	json.Unmarshal(b, v1)
	v1.ImageUuid = v2.Image
	v1.WorkingDir = v2.WorkDir
	v1.MemorySwap = v2.MemSwap

	return v1, nil
}
