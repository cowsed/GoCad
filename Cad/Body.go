package cad

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

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
func (b *Body) DrawChildren(CamMatrix mgl32.Mat4, ProjectionMatrix mgl32.Mat4) {
	for i := range b.Parts {
		b.Parts[i].DrawOperations(CamMatrix, ProjectionMatrix)
	}
}
func (b *Body) DrawChildrenIDs(CamMatrix mgl32.Mat4, ProjectionMatrix mgl32.Mat4) {
	for i := range b.Parts {
		b.Parts[i].DrawOperationIDs(CamMatrix, ProjectionMatrix)
	}
}

//InitGL initializes all of the graphics stuff for its sub parts
func (b *Body) InitGL() {
	for i := range b.Parts {
		b.Parts[i].InitGL()
	}
}

//HandleInput handles the input on children
func (b *Body) HandleInput() {
	for i := range b.Parts {
		b.Parts[i].HandleInput()
	}
}
