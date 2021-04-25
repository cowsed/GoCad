package cad

import (
	render "github.com/cowsed/GoCad/Render"
	"github.com/inkyblackness/imgui-go/v4"
)

//CurrentlySelected holds all of the currently selected objects
var CurrentlySelected = map[string]Selectable{}

//Project holds an entire cad project
//Can hold Bodies, parts, spreadsheets etc
type Project struct {
	Name  string
	Items []TreeItem
}

//InitGL initializes all the sub parts' gl stuff
func (p *Project) InitGL() {
	//Create render.SketchProgram
	render.SketchProgram = render.MakeProgram(baseVertSource, baseFragSource)

	//Set Program Uniforms
	render.SetUniform3f(render.SketchProgram, "normal_color", render.PointColor)
	render.SetUniform3f(render.SketchProgram, "selected_color", render.SelectionColor)
	render.SetUniform3f(render.SketchProgram, "hovered_color", render.HoverColor)

	for i := range p.Items {
		p.Items[i].InitGL()
	}
}

//DrawProject draws the project
func (p *Project) DrawProject() {
	for i := range p.Items {
		p.Items[i].Draw()
	}
}

//TreeItem is an interface for anything that can be held in a project
type TreeItem interface {
	SetPath(string)
	BuildTreeItem()
	InitGL() //Initialize all required VBOs, VAOs, and Programs
	Draw()   //Actually Draw the data
}

//Body is a holder for many parts that can be moved together
type Body struct {
	Name        string
	ParentPath  string
	TreePath    string
	Description string
	Show        bool
	Selected    bool
	Parts       []Part
}

//SetPath sets the path of the body
func (b *Body) SetPath(parent string) {
	b.TreePath = parent + "/" + b.Name
}

//Path gets the path in the project of the body
func (b *Body) Path() string {
	return b.TreePath
}

//Type returns the Type of the body for selections
func (b *Body) Type() SelectableType {
	return BodyType
}

//BuildTreeItem builds the ui for the treeview
func (b *Body) BuildTreeItem() {

	open := imgui.TreeNodeV(b.Name+"{Body}", imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	if open {
		if imgui.BeginPopupContextItemV(b.Name+"ContextMenu", imgui.PopupFlagsMouseButtonRight) {
			if imgui.Checkbox("Selected", &b.Selected) {
				if b.Selected {
					b.SetPath(b.ParentPath)
					CurrentlySelected[b.Path()] = b
				} else {
					delete(CurrentlySelected, b.Path())
				}
			}
			imgui.Text(b.Description)
			imgui.EndPopup()
		}
		for i := range b.Parts {
			b.Parts[i].BuildTreeItem()
		}

		imgui.TreePop()
	}

}

//Draw draws the body to the screen the opengl way
func (b *Body) Draw() {
	for i := range b.Parts {
		b.Parts[i].Draw()
	}
}

//InitGL initializes all of the graphics stuff for its sub parts
func (b *Body) InitGL() {
	for i := range b.Parts {
		b.Parts[i].InitGL()
	}
}

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
func (p *Part) Draw() {
	for i := range p.Chain {
		p.Chain[i].Draw()
	}
}

//InitGL initializes all of the graphics stuff for its sub parts
func (p *Part) InitGL() {
	for i := range p.Chain {
		p.Chain[i].InitGL()
	}
}

//Operation can Hold Sketches or extrudes revolves and such
type Operation interface {
	InitGL()        //Prepare OpenGL stuff
	Draw()          //Draw Opengl stuff
	BuildTreeItem() //Draw ui in tree
}
