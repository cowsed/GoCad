package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

//go:embed gl-shader/base.vert
var baseVertSource string

//go:embed gl-shader/base.frag
var baseFragSource string

var clearColor = [3]float32{0.0, 0.0, 0.0}
var pointColor = [3]float32{1.0, 0.0, 0.0}
var selectionColor = [3]float32{0.0, 1.0, 0.0}

var pointSize float32 = 8.0

var triangle = []float32{
	0, 0.75, 0, //top
	-0.5, -0.5, 0, //left
	0.5, -0.5, 0, //right
}
var colors = []float32{
	1, 0, 0, 1, //top
	0, 0, 1, 1, //left
	1, 0, 1, 1, //right
}

var (
	testVao     uint32
	testProgram uint32
)

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	//Something is going wrong here and it may have to do with the order of
	//vbo, vao and program creation

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	log.Println("GL", gl.ARRAY_BUFFER, "VBO", vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), nil, gl.STATIC_DRAW)

	//Bind Points
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(points), gl.Ptr(points))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	var vbo2 uint32
	gl.GenBuffers(1, &vbo2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	log.Println("GL", gl.ARRAY_BUFFER, "VBO", vbo2)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)

	//Bind Colors
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(colors), gl.Ptr(colors))

	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)

	return vao
}

func InitRender() {
	testVao = makeVao(triangle)

	makeProgram(baseVertSource, baseFragSource)

	//Enable color attribute
	var color_attribute uint32 = uint32(gl.GetAttribLocation(testProgram, gl.Str("color"+"\x00")))
	log.Println("color attr @", color_attribute)
	gl.VertexAttribPointer(color_attribute, 3, gl.FLOAT, false, 0, gl.Ptr(colors))
	gl.EnableVertexAttribArray(color_attribute)

	loc := gl.GetUniformLocation(testProgram, gl.Str("normal_color"+"\x00"))
	gl.ProgramUniform3f(testProgram, loc, pointColor[0], pointColor[1], pointColor[2])
	loc = gl.GetUniformLocation(testProgram, gl.Str("selected_color"+"\x00"))
	gl.ProgramUniform3f(testProgram, loc, selectionColor[0], selectionColor[1], selectionColor[2])

	//Set initial parameters
	gl.PointSize(pointSize)
}

//Renders the model defined by Current Project
func RenderModel() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(testProgram)

	gl.BindVertexArray(testVao)
	gl.DrawArrays(gl.POINTS, 0, int32(len(triangle)/3))
	//gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(triangle)/3))

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

func makeProgram(vShader, fShader string) {
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

	testProgram = gl.CreateProgram()
	gl.AttachShader(testProgram, vertexShader)
	gl.AttachShader(testProgram, fragmentShader)

	gl.LinkProgram(testProgram)

}
