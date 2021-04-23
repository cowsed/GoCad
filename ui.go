package main

import (
	"fmt"
	"log"

	"github.com/inkyblackness/imgui-go/v4"
)

func ShowDebugWindow() {
	imgui.Begin("Debug")
	imgui.Text(fmt.Sprintf("FPS: %.2f", imgui.CurrentIO().Framerate()))
	imgui.Text(fmt.Sprint("Selected", CurrentlySelected))
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
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}

}
