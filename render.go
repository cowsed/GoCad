package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

var clearColor = [3]float32{0.0, 0.0, 0.0}

var pointColor = [3]float32{1.0, 0.0, 0.0}
var selectionColor = [3]float32{1.0, 1.0, 0.0}
var hoverColor = [3]float32{0.0, 0.0, 1.0}

var pointSize float32 = 8.0

// makeVao initializes and returns a vertex array from the points provided.
/*
func makeVao(points []float32) uint32 {

	//Bind Points

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	//Bind positions
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	var vbo2 uint32
	gl.GenBuffers(1, &vbo2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertex_states), gl.Ptr(vertex_states), gl.STATIC_DRAW)

	//Bind States
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	gl.VertexAttribPointer(1, 1, gl.UNSIGNED_INT, false, 0, nil)

	return vao
}
*/
func InitRender() {

	/*	testVao = makeVao(triangle)

		testProgram = makeProgram(baseVertSource, baseFragSource)

		//Set Uniforms
		loc := gl.GetUniformLocation(testProgram, gl.Str("normal_color"+"\x00"))
		gl.ProgramUniform3f(testProgram, loc, pointColor[0], pointColor[1], pointColor[2])
		loc = gl.GetUniformLocation(testProgram, gl.Str("selected_color"+"\x00"))
		gl.ProgramUniform3f(testProgram, loc, selectionColor[0], selectionColor[1], selectionColor[2])
		loc = gl.GetUniformLocation(testProgram, gl.Str("hovered_color"+"\x00"))
		gl.ProgramUniform3f(testProgram, loc, hoverColor[0], hoverColor[1], hoverColor[2])
	*/
	currentProject.InitGL()

	//Set initial parameters
	gl.PointSize(pointSize)
}

//Renders the model defined by Current Project
func RenderModel() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	currentProject.Draw()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func makeProgram(vShader, fShader string) uint32 {
	var program uint32
	vertexShaderSource := vShader + "\x00"
	fragmentShaderSource := fShader + "\x00"

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	//Make and error Check program
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		infolog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(infolog))
		log.Fatalf("Failed to link program: %v", infolog)
	}
	return program
}

func setUniform3f(program uint32, uniformName string, value [3]float32) {
	loc := gl.GetUniformLocation(program, gl.Str(uniformName+"\x00"))
	gl.ProgramUniform3f(program, loc, value[0], value[1], value[2])

}
