package models

type HttpResponseJSON struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
