package cad

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/inkyblackness/imgui-go/v4"

	render "github.com/cowsed/GoCad/Render"
)

//ShowDebugWindow shows debug information
func ShowDebugWindow() {
	io := imgui.CurrentIO()
	imgui.Begin("Debug")
	imgui.Text(fmt.Sprintf("FPS: %.2f", io.Framerate()))

	imgui.Text(fmt.Sprint("Selected", CurrentlySelected))

	if imgui.TreeNode("Imgui Metrics") {
		imgui.Text(fmt.Sprintf("Imgui Capture Mouse: %v", io.WantCaptureMouse()))
		imgui.Text(fmt.Sprintf("Imgui Capture Keyboard: %v", io.WantCaptureKeyboard()))

		imgui.Text(fmt.Sprintf("Active Allocations: %v", io.MetricsActiveAllocations()))
		imgui.Text(fmt.Sprintf("Active Windows: %v", io.MetricsActiveWindows()))

		imgui.Text(fmt.Sprintf("Render Indices: %v", io.MetricsRenderIndices()))
		imgui.Text(fmt.Sprintf("Render Vertices: %v", io.MetricsRenderVertices()))
		imgui.TreePop()
	}
	imgui.End()

}

//ShowUI shows the treeview of the model
func ShowUI(p *Project) {
	imgui.Begin("Project")
	//imgui.PushStyleColor(imgui.StyleColor, imgui.Vec4{1, 0, 0, 1})

	if imgui.TreeNodeV(p.Name, imgui.TreeNodeFlagsOpenOnDoubleClick) {
		for i := range p.Items {
			p.Items[i].BuildTreeItem()
		}
		imgui.TreePop()
	}
	//imgui.PopStyleColor()

	imgui.End()

}

//ShowMainMenuBar creates the main menu bar
func ShowMainMenuBar(p *Project) {
	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("Save") {
				log.Println("Saving Project")
			}
			if imgui.MenuItem("Close") {
				log.Println("Closing Project")
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Edit") {
			if imgui.MenuItem("Undo") {
				log.Println("Undoing")
			}
			if imgui.MenuItem("Redo") {
				log.Println("Redoing")
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu("View") {
			if imgui.BeginMenu("FOV") {
				imgui.DragFloat("##FOV", &render.FOV)
				imgui.EndMenu()
			}
			if imgui.BeginMenu("Camera Type") {
				imgui.RadioButtonInt("Perspective", (*int)(&p.CamType), 0)
				imgui.RadioButtonInt("Orthographic", (*int)(&p.CamType), 1)

				imgui.EndMenu()
			}
			if imgui.BeginMenu("Clear Color") {
				imgui.ColorEdit3("##ClearColor", &render.ClearColor)
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Vertex Settings") {
				if imgui.BeginMenu("Normal Vertex Color") {
					if imgui.ColorEdit3("##VertexColor", &render.PointColor) {
						loc := gl.GetUniformLocation(render.SketchProgram, gl.Str("normal_color"+"\x00"))
						gl.ProgramUniform3f(render.SketchProgram, loc, render.PointColor[0], render.PointColor[1], render.PointColor[2])
					}
					imgui.EndMenu()
				}
				if imgui.BeginMenu("Selection Vertex Color") {
					if imgui.ColorEdit3("##SelectionColor", &render.SelectionColor) {
						loc := gl.GetUniformLocation(render.SketchProgram, gl.Str("selected_color"+"\x00"))
						gl.ProgramUniform3f(render.SketchProgram, loc, render.SelectionColor[0], render.SelectionColor[1], render.SelectionColor[2])
					}
					imgui.EndMenu()
				}
				if imgui.BeginMenu("Hovered Vertex Color") {
					if imgui.ColorEdit3("##HoverColor", &render.HoverColor) {
						loc := gl.GetUniformLocation(render.SketchProgram, gl.Str("hovered_color"+"\x00"))
						gl.ProgramUniform3f(render.SketchProgram, loc, render.HoverColor[0], render.HoverColor[1], render.HoverColor[2])
					}
					imgui.EndMenu()
				}

				imgui.EndMenu()
			}

			if imgui.BeginMenu("PointSize") {
				if imgui.SliderFloatV("##PointSizeControl", &render.PointSize, 0, 30, "%.3f", imgui.SliderFlagsNone) {
					gl.PointSize(render.PointSize)
				}
				imgui.EndMenu()
			}
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}

}
