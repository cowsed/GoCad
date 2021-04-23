package main

type Project struct {
	Name   string
	Bodies []Body
}

type Body struct {
	Name        string
	Description string
}

type Part struct {
	Name        string
	Description string
	Chain       []Operations
}

//Can Hold Sketches or extrudes revolves and such
type Operations interface {
}
