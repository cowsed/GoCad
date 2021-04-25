package main

import (
	"github.com/inkyblackness/imgui-go/v4"
)

var currentProject Project = Project{
	Name: "Test Project",
	Items: []TreeItem{
		&Body{
			Name:        "Square",
			TreePath:    "Test Project/Square",
			Description: "Just a square",
			Show:        true,
			Selected:    false,
			Parts: []Part{
				{
					Name:        "Square",
					Path:        "Test Project/Square/Square",
					Description: "Interior Part of it yknow",
					Chain: []Operation{
						&Sketch{
							name: "SquareSketch",
							path: "Test Project/Square/Square/SquareSketch",
							vertices: []SketchVertex{
								{0, 0.75},
								{-0.5, -0.5},
								{0.5, -0.5},
							},
						},
					},
				},
			}},
	},
}
var CurrentlySelected = map[string]Selectable{}

type Project struct {
	Name  string
	Items []TreeItem
}

func (p *Project) InitGL() {
	for i := range p.Items {
		p.Items[i].InitGL()
	}
}
func (p *Project) Draw() {
	for i := range p.Items {
		p.Items[i].Draw()
	}
}

type TreeItem interface {
	SetPath(string)
	BuildTreeItem()
	InitGL() //Initialize all required VBOs, VAOs, and Programs
	Draw()   //Actually Draw the data
}

type Body struct {
	Name        string
	TreePath    string
	Description string
	Show        bool
	Selected    bool
	Parts       []Part
}

func (b *Body) SetPath(parent string) {
	b.TreePath = parent + "/" + b.Name
}
func (b *Body) Path() string {
	return b.TreePath
}
func (b *Body) Type() SelectableType {
	return BodyType
}
func (b *Body) BuildTreeItem() {

	open := imgui.TreeNodeV(b.Name+"{Body}", imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	if open {
		if imgui.BeginPopupContextItemV(b.Name+"ContextMenu", imgui.PopupFlagsMouseButtonRight) {
			if imgui.Checkbox("Selected", &b.Selected) {
				if b.Selected {
					b.SetPath(currentProject.Name)
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
func (b *Body) Draw() {
	for i := range b.Parts {
		b.Parts[i].Draw()
	}
}
func (b *Body) InitGL() {
	for i := range b.Parts {
		b.Parts[i].InitGL()
	}
}

type Part struct {
	Name        string
	Path        string
	Description string
	Chain       []Operation
}

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
func (p *Part) Draw() {
	for i := range p.Chain {
		p.Chain[i].Draw()
	}
}
func (p *Part) InitGL() {
	for i := range p.Chain {
		p.Chain[i].InitGL()
	}
}

//Can Hold Sketches or extrudes revolves and such
type Operation interface {
	InitGL()        //Prepare OpenGL stuff
	Draw()          //Draw Opengl stuff
	BuildTreeItem() //Draw part in tree
}
