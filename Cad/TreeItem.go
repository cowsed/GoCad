package cad

import (
	"github.com/go-gl/mathgl/mgl32"
)

//TreeItem is an interface for anything that can be held in a project
type TreeItem interface {
	SetPath(string)
	BuildTreeItem()
	InitGL()                             //Initialize all required VBOs, VAOs, and Programs
	DrawChildren(mgl32.Mat4, mgl32.Mat4) //Actually Draw the data
	DrawChildrenIDs(mgl32.Mat4, mgl32.Mat4)
	HandleInput()
}
