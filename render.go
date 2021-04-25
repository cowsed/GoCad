package main

import (
	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func InitRender() {
	currentProject.InitGL()

	//Set initial parameters
	gl.PointSize(render.PointSize)
}

//Renders the model defined by Current Project
func RenderModel() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	currentProject.Draw()
}
