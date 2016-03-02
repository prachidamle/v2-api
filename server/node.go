package server

import (
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
)

func (s *Server) NodeByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	
	action := r.URL.Query().Get("action")
	if action != "" && r.Method == "POST" {
		return s.NodeAction(rw, r, action)
	}	
	
	return s.getNode(rw, r, vars["id"])
}

func (s *Server) NodeList(rw http.ResponseWriter, r *http.Request) error {
	return s.getNode(rw, r, "")
}

func (s *Server) getNode(rw http.ResponseWriter, r *http.Request, id string) error {
	id = s.deobfuscate(r, "node", id)

	rows, err := s.namedQuery(s.getNodesSQL(r, id), map[string]interface{}{
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
			ResourceType: "node",
		},
	}

	for rows.Next() {
		obj := &model.NodeV1{}
		obj.Type = "node"

		if err := rows.StructScan(obj); err != nil {
			return err
		}
		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "node", obj.Id)

		if err = s.parseData(obj.Data, obj); err != nil {
			return err
		}

		response.Data = append(response.Data, s.nodeV1ToV2(obj))
	}

	return s.writeResponse(rows.Err(), r, response)
}

func (s *Server) NodeAction(rw http.ResponseWriter, r *http.Request, action string) error {
	vars := mux.Vars(r)
	nodeObj := &model.NodeV2{}
	nodeObj.Resource.Type = "node"
	
	if s.isValidAction(action, &nodeObj.Resource) {
		switch(action){
		 	case "activate": return s.NodeActionActivate(rw, r)
		 	
		 	case "deactivate": return s.NodeActionDeactivate(rw, r)
		 	
		 	case "dockersocket": return s.NodeActionDockersocket(rw, r)
		 	
		 	case "purge": return s.NodeActionPurge(rw, r)
		 	
		 	case "remove": return s.NodeActionRemove(rw, r)
		 	
		 	case "restore": return s.NodeActionRestore(rw, r)
		 	
		 	case "update": return s.NodeActionUpdate(rw, r)
		}
	} else {
		return s.getContainer(rw, r, vars["id"])
	}
	
	return nil
}


func (s *Server) NodeActionActivate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionActivate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionDeactivate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionDeactivate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionDockersocket(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionDockersocket(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionPurge(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionPurge(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionRemove(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionRemove(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionRestore(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionRestore(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}

func (s *Server) NodeActionUpdate(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}
	
	v1Obj, err := rancherClient.Host.ById(id)
	if err != nil {
		return err
	}
	
	
	updatedV1Obj, err := rancherClient.Host.ActionUpdate(v1Obj)
	
	
	if err != nil {
		return err
	}

	return s.getNode(rw, r, updatedV1Obj.Id)
}


func (s *Server) nodeV1ToV2(obj *model.NodeV1) *model.NodeV2 {
	objV2 := &model.NodeV2{
		Resource:   obj.Resource,
		NodeCommon: obj.NodeCommon,
	}
	return objV2
}

func (s *Server) getNodesSQL(r *http.Request, id string) string {

	commonSql := s.getSql(r, &model.Common{})
	nodeSql := s.getSql(r, &model.NodeCommon{})
	if nodeSql != "" {
		nodeSql = `, ` + nodeSql
	}

	q := `SELECT COALESCE(id, '') as id, ` + commonSql + nodeSql +
		` FROM host
		  WHERE
			account_id = :account_id
			AND removed IS NULL`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}

