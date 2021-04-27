package cad

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

//Part is a part created by a chain of operations of sketches, extrudes, (analagous to the freecad PartDesign)
type Part struct {
	Name        string
	Path        string
	Description string
	Chain       []Operation
}

//BuildTreeItem builds the ui for the treeview
func (p *Part) BuildTreeItem() {
	open := imgui.TreeNodeV(p.Name+"{Part}", imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	if open {
		//Context Menu

		//Operations UI
		for i := range p.Chain {
			p.Chain[i].BuildTreeItem()
		}
		imgui.TreePop()
	}
}

//Draw draws all of its subparts in the opengl way
func (p *Part) DrawOperations(CamMatrix, ProjectionMatrix mgl32.Mat4) {
	for i := range p.Chain {
		p.Chain[i].Draw(CamMatrix, ProjectionMatrix)
	}
}
func (p *Part) DrawOperationIDs(CamMatrix, ProjectionMatrix mgl32.Mat4) {
	for i := range p.Chain {
		p.Chain[i].DrawIDs(CamMatrix, ProjectionMatrix)
	}
}

//InitGL initializes all of the graphics stuff for its sub parts
func (p *Part) InitGL() {
	for i := range p.Chain {
		p.Chain[i].InitGL()
	}
}
func (p *Part) HandleInput() {
	for i := range p.Chain {
		p.Chain[i].HandleInput()
	}
}

//Operation can Hold Sketches or extrudes revolves and such
type Operation interface {
	InitGL()                        //Prepare OpenGL stuff
	Draw(mgl32.Mat4, mgl32.Mat4)    //Draw Opengl stuff
	DrawIDs(mgl32.Mat4, mgl32.Mat4) //Draw Opengl stuff with IDs

	BuildTreeItem() //Draw ui in tree
	HandleInput()   //Handles Input
}
