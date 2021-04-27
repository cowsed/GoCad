package main

import (
	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

//InitRender Initializes all the objects that need to be initialized
func InitRender() {
	currentProject.InitGL()

	//Set initial parameters
	gl.PointSize(render.PointSize)
}

//RenderModel Renders the model defined by Current Project
func RenderModel() {

	currentProject.DrawIDs()
	currentProject.DrawProject()
}

//HandleInput Handles the input from the gl window
func HandleInput(window *glfw.Window) {
	currentProject.HandleInput(window)
}
