package models

type HTTPResponseJSON struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
