package server

import (
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
)

func (s *Server) NodeByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
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

func (s *Server) NodeCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}

	data := s.parseInputParameters(r)

	environ, err := rancherClient.Environment.Create(&client.Environment{
		Name: data.String("name"),
	})

	if err != nil {
		return err
	}

	return s.getNode(rw, r, environ.Id)
}
