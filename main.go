package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	renderers "github.com/cowsed/GoCad/Renderer"
	"github.com/inkyblackness/imgui-go/v4"
)

var clearColor = [3]float32{0.0, 0.0, 0.0}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	//Set up glfw and gl
	fmt.Println("Initializing...")
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := renderers.NewGLFW(io, renderers.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer renderer.Dispose()

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
		ShowMainMenuBar(&currentProject)
		ShowUI(&currentProject)
		ShowDebugWindow()
		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		// app.RenderScene()

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

		// sleep to avoid 100% CPU usage for this demo
		//<-time.After(sleepDuration)
	}
}
