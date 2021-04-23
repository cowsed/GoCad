package main

import (
	"log"

	"github.com/inkyblackness/imgui-go/v4"
)

func ShowUI(p *Project) {
	imgui.Begin("Project")
	if imgui.TreeNodeV(p.Name, imgui.TreeNodeFlagsFramed) {
		for _, b := range p.Bodies {
			if imgui.TreeNodeV(b.Name, imgui.TreeNodeFlagsFramed) {
				imgui.Text(b.Description)
				imgui.TreePop()
			}
		}
		imgui.TreePop()
	}
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
