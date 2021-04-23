package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go/v4"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	//Set up glfw and gl
	fmt.Println("Starting")
	err := glfw.Init()
	if err != nil {
		log.Fatal("Could not initialize GLFW", err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "GoCad", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		log.Fatal(err)
	}
	//Setup IMGUI
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		log.Fatalf("%v\n", err)

	}
	defer renderer.Dispose()

	Run(window)
}
func Run(window *glfw.Window) {
	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)
		gl.ClearColor(1, 0, 0, 1)

		//Do Imgui Stuff
		imgui.Render()

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
