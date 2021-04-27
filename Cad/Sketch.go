package cad

import (
	_ "embed" //Embed for loading shaders from file system
	"fmt"

	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

//go:embed gl-shader/base_vertex.vert
var baseVertSource string

//go:embed gl-shader/base_vertex.frag
var baseFragSource string

//SketchVertex is a vertex in a sketch. Need only be 2d
type SketchVertex struct {
	X, Y float32
}

//Sketch is the basic cad sketch
type Sketch struct {
	//Sketch Stuff
	Name     string
	Path     string
	Vertices []SketchVertex

	vao            uint32
	disposeVao     bool
	positionVbo    uint32 //Handle to VBO for draw_points
	stateVbo       uint32 //Handle to VBO for vertex_states
	disposeVbos    bool
	glPoints       []float32 //Holds opengl representation of points
	glVertexStates []uint32  //Holds selection, hover information about points
}

//InitGL sets up all necessary GL stuff
func (s *Sketch) InitGL() {
	//Delete Old OpenGL Stuff
	if s.disposeVbos {
		gl.DeleteBuffers(1, &s.positionVbo)
		gl.DeleteBuffers(1, &s.stateVbo)
	}
	if s.disposeVao {
		gl.DeleteVertexArrays(1, &s.vao)
	}

	//Generate VBO, VAO
	s.MakeDrawData()
	s.makeVao()
	s.UpdateDrawData()

}

//makeVao creates the sketch vao generates vbos
func (s *Sketch) makeVao() {
	//Bind Points
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	//Create postopm vbo
	gl.GenBuffers(1, &s.positionVbo)
	//Create state vbo
	gl.GenBuffers(1, &s.stateVbo)

}

//MakeDrawData transforms the sketch points and lines to opengl points and lines
func (s *Sketch) MakeDrawData() {
	s.glPoints = make([]float32, 3*len(s.Vertices))
	for i := range s.Vertices {
		s.glPoints[i*3] = s.Vertices[i].X
		s.glPoints[i*3+1] = s.Vertices[i].Y
		s.glPoints[i*3+2] = 0.0
	}

	s.glVertexStates = []uint32{2, 0, 2, 0}
}

//UpdateDrawData takes the newly made draw data and uploads it to the gpu
//Needs a way to check if the length has changed cuz then new vao and vbos must be created
func (s *Sketch) UpdateDrawData() {
	gl.BindVertexArray(s.vao)

	//Bind positions
	gl.BindBuffer(gl.ARRAY_BUFFER, s.positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.glPoints), gl.Ptr(s.glPoints), gl.STATIC_DRAW)

	//Enable positions
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	//Bind States
	gl.BindBuffer(gl.ARRAY_BUFFER, s.stateVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.glVertexStates), gl.Ptr(s.glVertexStates), gl.STATIC_DRAW)

	//Enable States
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.stateVbo)
	gl.VertexAttribPointer(1, 1, gl.UNSIGNED_INT, false, 0, nil)

}

//Draw Actually draws the sketch
func (s *Sketch) Draw(CamMatrix, ProjectionMatrix mgl32.Mat4) {

	MVP := ProjectionMatrix.Mul4(CamMatrix)
	MVPUniform := gl.GetUniformLocation(render.SketchProgram, gl.Str("MVP"+"\x00"))
	gl.UniformMatrix4fv(MVPUniform, 1, false, &MVP[0])

	gl.UseProgram(render.SketchProgram)

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(s.glPoints)/3))

	gl.DrawArrays(gl.POINTS, 0, int32(len(s.glPoints)/3))

}
func (s *Sketch) DrawIDs(CamMatrix, ProjectionMatrix mgl32.Mat4) {

	MVP := ProjectionMatrix.Mul4(CamMatrix)
	MVPUniform := gl.GetUniformLocation(render.SketchProgram, gl.Str("MVP"+"\x00"))
	gl.UniformMatrix4fv(MVPUniform, 1, false, &MVP[0])

	gl.UseProgram(render.SketchProgram)

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(s.glPoints)/3))

	gl.DrawArrays(gl.POINTS, 0, int32(len(s.glPoints)/3))

}

func (s *Sketch) HandleInput() {

}

//SetPath sets the ppath of a sketch
func (s *Sketch) SetPath(parent string) {
	s.Path = parent + "/" + s.Name
}

//BuildTreeItem creates the ui for the treeview
func (s *Sketch) BuildTreeItem() {
	open := imgui.TreeNodeV(s.Name+"{Sketch}", imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	if open {
		imgui.Button("Solve(N/A)")

		imgui.BeginTableV(s.Name+"Vertices", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		imgui.Text("Vertices")
		imgui.TableNextColumn()
		imgui.Text("X")
		imgui.TableNextColumn()
		imgui.Text("Y")
		imgui.TableNextRow()
		for i := range s.Vertices {
			imgui.TableNextColumn()

			imgui.Selectable(fmt.Sprint(s.Vertices[i].X))

			imgui.TableNextColumn()
			imgui.Selectable(fmt.Sprint(s.Vertices[i].Y))
			imgui.TableNextRow()
		}
		imgui.EndTable()
		imgui.TreePop()
	}
}
