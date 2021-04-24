package main

import (
	"log"

	"github.com/inkyblackness/imgui-go/v4"
)

var currentProject Project = Project{"Test Project", []TreeItem{&Body{"Square", "Test Project/Square", "Just a square", true, false}}}
var CurrentlySelected = map[string]Selectable{}

type Project struct {
	Name  string
	Items []TreeItem
}
type TreeItem interface {
	SetPath(string)
	BuildTreeItem()
}

type Body struct {
	Name        string
	TreePath    string
	Description string
	Show        bool
	Selected    bool
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

	open := imgui.TreeNodeV(b.Name, imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	imgui.SameLine()
	if imgui.Button("Edit") {
		log.Println("Edit", b)
	}

	if open {
		if imgui.Checkbox("Selected", &b.Selected) {
			if b.Selected {
				b.SetPath(currentProject.Name)
				CurrentlySelected[b.Path()] = b
			} else {
				delete(CurrentlySelected, b.Path())
			}
		}
		imgui.Text(b.Description)
		imgui.TreePop()
	}

}

type Part struct {
	Name        string
	Description string
	Chain       []Operations
}

//Can Hold Sketches or extrudes revolves and such
type Operations interface {
}
