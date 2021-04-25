package render

//PointSize Controls the size of vertices in rendering
var PointSize float32 = 8.0

//ClearColor is the opengl Clear Color
var ClearColor = [3]float32{0.0, 0.0, 0.0}

//PointColor is the color of a normal point
var PointColor = [3]float32{1.0, 0.0, 0.0}

//SelectionColor is the color of a selected point
var SelectionColor = [3]float32{1.0, 1.0, 0.0}

//HoverColor is the color of a hovered point
var HoverColor = [3]float32{0.0, 0.0, 1.0}

//FOV of the cameras in degrees
var FOV float32 = 45.0

var WindowWidth int
var WindowHeight int

var MouseSensitivity float32 = 3.0

//SketchProgram is the program used by all sketches
//For now created in the project, realistically should be done in a render.InitGL
var SketchProgram uint32
