package main

import (
	"log"
	"runtime"

	cad "github.com/cowsed/GoCad/Cad"
	renderers "github.com/cowsed/GoCad/Platform"
	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/imgui-go/v4"
)

var currentProject cad.Project = cad.Project{
	Name:   "Test Project",
	CamPos: mgl32.Vec3{0, 0, 3},
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
								{0, 1},
								{1, -0},
								{0, -1},
								{-1, -0},
							},
						},
					},
				},
			},
		},
	},
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	//Setup IMGUI and other graphics
	log.Println("Initializing...")
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := renderers.NewGLFW(io, renderers.GLFWClientAPIOpenGL3)
	if err != nil {
		log.Fatalf("Error Creating Platform: %v\n", err)

	}
	defer platform.Dispose()

	window := platform.Window()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		log.Fatalf("Error Creating Renderer: %v\n", err)
	}
	defer renderer.Dispose()

	windowWidth, windowHeight := window.GetSize()
	render.WindowWidth = windowWidth
	render.WindowHeight = windowHeight
	render.WindowAspect = float32(render.WindowWidth) / float32(render.WindowHeight)

	InitRender()

	Run(platform, renderer, window)
}

//Run runs the program
func Run(p *renderers.GLFW, r *renderers.OpenGL3, window *glfw.Window) {
	imgui.CurrentIO().SetClipboard(renderers.Clipboard{Platform: p})

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		//Draw UI
		cad.ShowMainMenuBar(&currentProject)
		cad.ShowUI(&currentProject)
		cad.ShowDebugWindow(&currentProject, window)
		// Render UI
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(render.ClearColor)

		//Accept input to the Cad Project
		if !imgui.CurrentIO().WantCaptureMouse() {
			HandleInput(window)
		}

		//Render the CAD Project
		RenderModel()

		//Actually draw to the screen
		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

	}
}
