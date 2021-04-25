package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
)

//go:embed gl-shader/base_vertex.vert
var baseVertSource string

//go:embed gl-shader/base_vertex.frag
var baseFragSource string

type SketchVertex struct {
	x, y float32
}

type Sketch struct {
	//Sketch Stuff
	name     string
	path     string
	vertices []SketchVertex

	//Opengl Stuff
	program          uint32
	disposeProgram   bool
	vao              uint32
	disposeVao       bool
	positionVbo      uint32 //Handle to VBO for draw_points
	stateVbo         uint32 //Handle to VBO for vertex_states
	disposeVbos      bool
	gl_points        []float32 //Holds opengl representation of points
	gl_vertex_states []uint32  //Holds selection, hover information about points
}

func (s *Sketch) InitGL() {
	//Delete Old OpenGL Stuff
	if s.disposeVbos {
		gl.DeleteBuffers(1, &s.positionVbo)
		gl.DeleteBuffers(1, &s.stateVbo)
	}
	if s.disposeProgram {
		gl.DeleteProgram(s.program)
	}
	if s.disposeVao {
		gl.DeleteVertexArrays(1, &s.vao)
	}
	log.Println("Initiating sketch")
	//Create Program (TODO May want to make a single program for all sketches when there can be multiple)
	s.program = makeProgram(baseVertSource, baseFragSource)
	//Set Program Uniforms
	setUniform3f(s.program, "normal_color", pointColor)
	setUniform3f(s.program, "selected_color", selectionColor)
	setUniform3f(s.program, "hovered_color", hoverColor)

	//Generate VBO, VAO
	s.MakeDrawData()
	s.makeVao()
	s.UpdateDrawData()

	log.Println("vao", s.vao)
	log.Println("posVBO", s.positionVbo)
	log.Println("stateVBO", s.stateVbo)

}

func (s *Sketch) makeVao() {
	//Bind Points
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	//Create postopm vbo
	gl.GenBuffers(1, &s.positionVbo)
	//Create state vbo
	gl.GenBuffers(1, &s.stateVbo)

}

func (s *Sketch) MakeDrawData() {
	s.gl_points = make([]float32, 3*len(s.vertices))
	for i := range s.vertices {
		s.gl_points[i*3] = s.vertices[i].x
		s.gl_points[i*3+1] = s.vertices[i].y
		s.gl_points[i*3+2] = 0.0
	}
	log.Println(s.gl_points)

	s.gl_vertex_states = []uint32{2, 2, 2}
}

func (s *Sketch) UpdateDrawData() {
	gl.BindVertexArray(s.vao)

	//Bind positions
	gl.BindBuffer(gl.ARRAY_BUFFER, s.positionVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.gl_points), gl.Ptr(s.gl_points), gl.STATIC_DRAW)

	//Enable positions
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.positionVbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	//Bind States
	gl.BindBuffer(gl.ARRAY_BUFFER, s.stateVbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.gl_vertex_states), gl.Ptr(s.gl_vertex_states), gl.STATIC_DRAW)

	//Enable States
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.stateVbo)
	gl.VertexAttribPointer(1, 1, gl.UNSIGNED_INT, false, 0, nil)

}

func (s *Sketch) Draw() {

	gl.UseProgram(s.program)

	gl.BindVertexArray(s.vao)
	gl.DrawArrays(gl.POINTS, 0, int32(len(s.gl_points)/3))
	gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(s.gl_points)/3))

}

func (s *Sketch) SetPath(parent string) {
	s.path = parent + "/" + s.name
}
func (s *Sketch) BuildTreeItem() {
	open := imgui.TreeNodeV(s.name+"{Sketch}", imgui.TreeNodeFlagsAllowItemOverlap+imgui.TreeNodeFlagsOpenOnDoubleClick)
	if open {
		imgui.BeginTableV(s.name+"Vertices", 2, imgui.TableFlagsBorders, imgui.Vec2{}, 0)
		imgui.Text("Vertices")
		imgui.TableNextColumn()
		imgui.Text("X")
		imgui.TableNextColumn()
		imgui.Text("Y")
		imgui.TableNextRow()
		for i := range s.vertices {
			imgui.TableNextColumn()

			imgui.Selectable(fmt.Sprint(s.vertices[i].x))

			imgui.TableNextColumn()
			imgui.Selectable(fmt.Sprint(s.vertices[i].y))
			imgui.TableNextRow()
		}
		imgui.EndTable()
		imgui.TreePop()
	}
}
