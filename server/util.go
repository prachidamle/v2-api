package server

import (
	"net/http"
	"reflect"
	"strings"
	"encoding/json"
)

func (s *Server) getAccountID(r *http.Request) int64 {
	return 5
}

func (s *Server) getSql(r *http.Request, obj interface{}) string {
	var sql string

	val := reflect.ValueOf(obj).Elem()

	for i := 0; i < val.NumField(); i++ {

		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag.Get("db") != "" {
			sql = sql + "COALESCE(" + tag.Get("db") + ", '') as " + tag.Get("db") + ", "
		}
	}

	return strings.TrimSuffix(sql, ", ")
}

func (s *Server) parseData(dataStr string, obj interface{}) error {

	type Data struct {
		Fields interface{} `json:"fields"`
	}

	data := Data{}

	bytes := []byte(dataStr)

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	fieldsM, err := json.Marshal(data.Fields)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(fieldsM, &obj); err != nil {
		return err
	}

	return nil
}