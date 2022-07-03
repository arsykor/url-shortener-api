package entities

import "reflect"

type Link struct {
	Id       string `json:"id"`
	URL      string `json:"url"`
	Requests int    `json:"requests"`
}

type LinkRequest struct {
	URL string `json:"url"`
}

func (lr *LinkRequest) ParamName() string {
	return reflect.Indirect(reflect.ValueOf(&LinkRequest{})).Type().Field(0).Name
}

type LinkResponse struct {
	Link string `json:"link"`
}

type LinkInMemory struct {
	Id       string `json:"id"`
	URL      string `json:"url"`
	Requests int    `json:"requests"`
	Deleted  bool   `json:"deleted"`
}
