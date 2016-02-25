package server

import (
	"net/http"
	"reflect"
	"strings"
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
