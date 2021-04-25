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

var SketchProgram uint32
