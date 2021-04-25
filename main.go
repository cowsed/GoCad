package main

import (
	"fmt"
	"log"
	"runtime"

	cad "github.com/cowsed/GoCad/Cad"
	renderers "github.com/cowsed/GoCad/Platform"
	render "github.com/cowsed/GoCad/Render"

	"github.com/inkyblackness/imgui-go/v4"
)

var currentProject cad.Project = cad.Project{
	Name: "Test Project",
	Items: []cad.TreeItem{
		&cad.Body{
			Name:        "Square",
			TreePath:    "Test Project/Square",
			Description: "Just a square",
			Show:        true,
			Selected:    false,
			Parts: []cad.Part{
				{
					Name:        "Square",
					Path:        "Test Project/Square/Square",
					Description: "Interior Part of it yknow",
					Chain: []cad.Operation{
						&cad.Sketch{
							Name: "SquareSketch",
							Path: "Test Project/Square/Square/SquareSketch",
							Vertices: []cad.SketchVertex{
								{0, 0.75},
								{-0.5, -0.5},
								{0.5, -0.5},
							},
						},
					},
				},
			}},
	},
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	//Setup IMGUI and other graphics
	fmt.Println("Initializing...")
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := renderers.NewGLFW(io, renderers.GLFWClientAPIOpenGL3)
	if err != nil {
		log.Fatalf("Error Creating Platform: %v\n", err)

	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		log.Fatalf("Error Creating Renderer: %v\n", err)
	}
	defer renderer.Dispose()

	InitRender()

	Run(platform, renderer)
}

func Run(p *renderers.GLFW, r *renderers.OpenGL3) {
	imgui.CurrentIO().SetClipboard(renderers.Clipboard{Platform: p})

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		//Draw UI
		cad.ShowMainMenuBar(&currentProject)
		cad.ShowUI(&currentProject)
		cad.ShowDebugWindow()
		// Render UI
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(render.ClearColor)

		//Accept input to the Cad Project
		//if imgui.CurrentIO().WantCaptureMouse() {
		//}

		//Render the CAD Project
		RenderModel()
		//Actually draw to the screen
		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

		// sleep to avoid 100% CPU usage for this demo
		//<-time.After(sleepDuration)
	}
}
