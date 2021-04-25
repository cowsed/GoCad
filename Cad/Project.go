package cad

import (
	"log"

	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

//CurrentlySelected holds all of the currently selected objects
var CurrentlySelected = map[string]Selectable{}

//Project holds an entire cad project
//Can hold Bodies, parts, spreadsheets etc
type Project struct {
	Name  string
	Items []TreeItem

	ProjectionMatrix mgl32.Mat4

	CamType   render.CameraType
	CamPos    mgl32.Vec3
	CamAng    mgl32.Vec3
	CamMatrix mgl32.Mat4
}

//InitGL initializes all the sub parts' gl stuff
func (p *Project) InitGL() {
	//Create render.SketchProgram
	render.SketchProgram = render.MakeProgram(baseVertSource, baseFragSource)

	//Set Program Uniforms
	render.SetUniform3f(render.SketchProgram, "normal_color", render.PointColor)
	render.SetUniform3f(render.SketchProgram, "selected_color", render.SelectionColor)
	render.SetUniform3f(render.SketchProgram, "hovered_color", render.HoverColor)

	//Create Initial Camera Position
	p.UpdateMatrices()
	//Init Children
	for i := range p.Items {
		p.Items[i].InitGL()
	}
}
func (p *Project) UpdateMatrices() {
	p.CamMatrix = mgl32.LookAtV(p.CamPos, p.CamPos.Sub(mgl32.Vec3{0, 0, 3}), mgl32.Vec3{0, 1, 0})
	if p.CamType == render.Perspective {
		p.ProjectionMatrix = mgl32.Perspective(mgl32.DegToRad(render.FOV), float32(render.WindowWidth)/float32(render.WindowHeight), 0.1, 10)
	} else {
		var size float32 = 2
		var aspect float32 = float32(render.WindowWidth) / float32(render.WindowHeight)
		p.ProjectionMatrix = mgl32.Ortho(-size, size, -size/aspect, size/aspect, .1, 10.0)

	}
}

//DrawProject draws the project
func (p *Project) DrawProject() {
	p.UpdateMatrices()
	for i := range p.Items {
		p.Items[i].DrawChildren(p.CamMatrix, p.ProjectionMatrix)
	}
}
func (p *Project) HandleInput(window *glfw.Window) {
	//Camera controls
	act := window.GetMouseButton(glfw.MouseButton3)
	if act == glfw.Press {
		//Get mouse delta
		//do some math
		d := imgui.CurrentIO().MouseDelta()
		fx := float32(render.WindowWidth)
		fy := float32(render.WindowHeight)
		p.CamPos = p.CamPos.Add(mgl32.Vec3{render.MouseSensitivity * -d.X / fx, render.MouseSensitivity * d.Y / fy, 0})
		log.Println(p.CamPos)
	}
	//Select All

	//Rest of it
	for i := range p.Items {
		p.Items[i].HandleInput()
	}
}

//TreeItem is an interface for anything that can be held in a project
type TreeItem interface {
	SetPath(string)
	BuildTreeItem()
	InitGL()                             //Initialize all required VBOs, VAOs, and Programs
	DrawChildren(mgl32.Mat4, mgl32.Mat4) //Actually Draw the data
	HandleInput()
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
func (b *Body) DrawChildren(CamMatrix mgl32.Mat4, ProjectionMatrix mgl32.Mat4) {
	for i := range b.Parts {
		b.Parts[i].DrawOperations(CamMatrix, ProjectionMatrix)
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
	InitGL()                     //Prepare OpenGL stuff
	Draw(mgl32.Mat4, mgl32.Mat4) //Draw Opengl stuff
	BuildTreeItem()              //Draw ui in tree
	HandleInput()                //Handles Input
}
