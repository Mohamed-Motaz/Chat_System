package Server

import db "Server/Database"

type AddApplicationReq struct {
	Name string `json:"name"`
}
type AddApplicationRes struct {
	db.Application
}

type UpdateApplicationReq struct {
	Name string `json:"name"`
}

type UpdateApplicationRes struct {
}
