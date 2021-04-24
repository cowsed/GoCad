package main

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
)

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
			if imgui.BeginMenu("Clear Color") {
				imgui.ColorEdit3("##ClearColor", &clearColor)
				imgui.EndMenu()
			}
			if imgui.BeginMenu("Vertex Settings") {
				if imgui.BeginMenu("Normal Vertex Color") {
					if imgui.ColorEdit3("##VertexColor", &pointColor) {
						loc := gl.GetUniformLocation(testProgram, gl.Str("normal_color"+"\x00"))
						gl.ProgramUniform3f(testProgram, loc, pointColor[0], pointColor[1], pointColor[2])
					}
					imgui.EndMenu()
				}
				if imgui.BeginMenu("Selection Vertex Color") {
					if imgui.ColorEdit3("##SelectionColor", &selectionColor) {
						loc := gl.GetUniformLocation(testProgram, gl.Str("selected_color"+"\x00"))
						gl.ProgramUniform3f(testProgram, loc, selectionColor[0], selectionColor[1], selectionColor[2])
					}
					imgui.EndMenu()
				}
				imgui.EndMenu()
			}

			if imgui.BeginMenu("PointSize") {
				if imgui.DragFloatV("##PointSizeControl", &pointSize, .01, 0, 30, "%.3f", imgui.SliderFlagsNone) {
					gl.PointSize(pointSize)
				}
				imgui.EndMenu()
			}
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}

}
