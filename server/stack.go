package server

import (
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
	"encoding/json"
	"fmt"
)

func (s *Server) StackByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	action := r.URL.Query().Get("action")
	if action != "" && r.Method == "POST" {
		return s.StackAction(rw, r, action)
	}
	return s.getStack(rw, r, vars["id"])
}

func (s *Server) StackList(rw http.ResponseWriter, r *http.Request) error {
	return s.getStack(rw, r, "")
}

func (s *Server) getStack(rw http.ResponseWriter, r *http.Request, id string) error {
	id = s.deobfuscate(r, "stack", id)

	fmt.Printf("\ngetStack ID returned: %v", id)
	rows, err := s.namedQuery(s.getStacksSQL(r, id), map[string]interface{}{
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
			ResourceType: "stack",
		},
	}

	for rows.Next() {
		obj := &model.StackV1{}
		obj.Type = "stack"

		if err := rows.StructScan(obj); err != nil {
			return err
		}
		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "stack", obj.Id)

		if err = s.parseData(obj.Data, obj); err != nil {
			return err
		}
		
		objV2, err := stackToV2(r, obj)
		
		s.addResourceActions(r, &objV2.Resource, objV2.State)
		
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

func (s *Server) StackCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	v2 := &model.StackV2{}
	if err := s.parseInputParameters(r, v2); err != nil {
		return err
	}
	
	v1, err := stackFromV2(v2)
	if err != nil {
		return err
	}
	
	v1.Type = "environment"

	environ, err := rancherClient.Environment.Create(v1)

	if err != nil {
		return err
	}

	return s.getStack(rw, r, environ.Id)
}

func (s *Server) StackAction(rw http.ResponseWriter, r *http.Request, action string) error {
	vars := mux.Vars(r)
	stackObj := &model.StackV2{}
	stackObj.Resource.Type = "stack"
	
	if s.isValidAction(action, &stackObj.Resource) {
		switch(action) {
			
		 	case "activateservices": return s.StackActionActivateservices(rw, r)
		 	
		 	case "addoutputs": return s.StackActionAddoutputs(rw, r)
		 	
		 	case "cancelrollback": return s.StackActionCancelrollback(rw, r)
		 	
		 	case "cancelupgrade": return s.StackActionCancelupgrade(rw, r)
		 	
		 	case "deactivateservices": return s.StackActionDeactivateservices(rw, r)
		 	
		 	case "error": return s.StackActionError(rw, r)
		 	
		 	case "exportconfig": return s.StackActionExportconfig(rw, r)
		 	
		 	case "finishupgrade": return s.StackActionFinishupgrade(rw, r)
		 	
		 	case "remove": return s.StackActionRemove(rw, r)
		 	
		 	case "rollback": return s.StackActionRollback(rw, r)
		 	
		 	case "update": return s.StackActionUpdate(rw, r)
		 	
		 	case "upgrade": return s.StackActionUpgrade(rw, r)
		 	
		}
	}else {
		return s.getStack(rw, r, vars["id"])
	}
	
	return nil
}

/*func (s *Server) StackActionAddoutputs(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	//bytes, _ := ioutil.ReadAll(r.Body)
	//fmt.Printf("\n\nIN the Stack addOutputs %v", string(bytes))

	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	environ, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}

	input := &client.AddOutputsInput{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	
	updatedEnviron, err := rancherClient.Environment.ActionAddoutputs(environ, input)

	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedEnviron.Id)
}


func (s *Server) StackActionActivateservices(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	//bytes, _ := ioutil.ReadAll(r.Body)
	//fmt.Printf("\n\nIN the Stack addOutputs %v", string(bytes))

	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	environ, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}

	updatedEnviron, err := rancherClient.Environment.ActionActivateservices(environ)

	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedEnviron.Id)
}
*/

func (s *Server) StackActionActivateservices(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionActivateservices(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionAddoutputs(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.AddOutputsInput{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Environment.ActionAddoutputs(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionCancelrollback(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionCancelrollback(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionCancelupgrade(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionCancelupgrade(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionDeactivateservices(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionDeactivateservices(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionError(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionError(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionExportconfig(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.ComposeConfigInput{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Environment.ActionExportconfig(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionFinishupgrade(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionFinishupgrade(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionRemove(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionRemove(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionRollback(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionRollback(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionUpdate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Environment.ActionUpdate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func (s *Server) StackActionUpgrade(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Environment.ById(id)
	if err != nil {
		return err
	}
	
	
	input := &client.EnvironmentUpgrade{}
	if err := s.parseInputParameters(r, input); err != nil {
		return err
	}
	updatedV1Obj, err := rancherClient.Environment.ActionUpgrade(v1Obj, input)
	
	
	if err != nil {
		return err
	}

	return s.getStack(rw, r, updatedV1Obj.Id)
}

func stackToV2(r *http.Request, v1 *model.StackV1) (*model.StackV2, error) {
	common := v1.StackCommon
	common.Transitioning = model.GetTransitioning(common.State, common.Transitioning)
	
	objV2 := &model.StackV2{
		Resource:    v1.Resource,
		StackCommon: common,
		Variables:   v1.Environment,
	}
	objV2.Type = "stack"

	return objV2, nil
}

func stackFromV2(v2 *model.StackV2) (*client.Environment, error) {
	b, err := json.Marshal(v2)
	if err != nil {
		return nil, err
	}
	v1 := &client.Environment{}
	json.Unmarshal(b, v1)
	v1.Environment = v2.Variables
	return v1, nil
}


func (s *Server) getStacksSQL(r *http.Request, id string) string {
	commonSql := s.getSql(r, &model.Common{})
	stackSql := s.getSql(r, &model.StackCommon{})
	if stackSql != "" {
		stackSql = `, ` + stackSql
	}

	q := `SELECT COALESCE(id, '') as id, ` + commonSql + stackSql +
		` FROM environment
		  WHERE
			account_id = :account_id
			AND removed IS NULL`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}



