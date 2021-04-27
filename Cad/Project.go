package cad

import (
	"image"
	"image/png"
	"log"
	"os"
	"unsafe"

	render "github.com/cowsed/GoCad/Render"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

//CurrentlySelected holds all of the currently selected objects
var CurrentlySelected = map[string]Selectable{}

//Project holds an entire cad project
//Can hold Bodies, parts, spreadsheets etc
type Project struct {
	Name  string
	Items []TreeItem

	ProjectionMatrix mgl32.Mat4

	IDFramebuffer uint32
	IDTex         uint32
	CamType       render.CameraType
	CamPos        mgl32.Vec3
	CamAng        mgl32.Vec3
	CamMatrix     mgl32.Mat4
}

//InitGL initializes all the sub parts' gl stuff
func (p *Project) InitGL() {
	//Create ID Framebuffer
	gl.CreateFramebuffers(1, &p.IDFramebuffer)

	//Create and setup texture to view ID Framebuffer
	gl.GenTextures(1, &p.IDTex)
	gl.BindTexture(gl.TEXTURE_2D, p.IDTex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(render.WindowWidth), int32(render.WindowHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//Link the two
	gl.BindFramebuffer(gl.FRAMEBUFFER, p.IDFramebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, p.IDTex, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) == gl.FRAMEBUFFER_COMPLETE {
		log.Println("Framebuffer OK")
	}

	//Create render.SketchProgram
	render.SketchProgram = render.MakeProgram(baseVertSource, baseFragSource)

	//Set Program Uniforms
	render.SetUniform3f(render.SketchProgram, "normal_color", render.PointColor)
	render.SetUniform3f(render.SketchProgram, "selected_color", render.SelectionColor)
	render.SetUniform3f(render.SketchProgram, "hovered_color", render.HoverColor)

	//Create Initial Camera Position
	p.UpdateMatrices()
	//Init Children
	for i := range p.Items {
		p.Items[i].InitGL()
	}
}

//SaveBuf somewhat of a debug function, saves the ID Framebuffer texture to a file
func (p *Project) SaveBuf() {
	outImage := image.NewRGBA(image.Rect(0, 0, render.WindowWidth, render.WindowHeight))
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, p.IDTex)

	gl.GetTexImage(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&outImage.Pix[0]))
	f, _ := os.Create("out.png")
	png.Encode(f, outImage)
}
func (p *Project) UpdateMatrices() {
	p.CamMatrix = mgl32.LookAtV(p.CamPos, p.CamPos.Sub(mgl32.Vec3{0, 0, 3}), mgl32.Vec3{0, 1, 0})
	if p.CamType == render.Perspective {
		p.ProjectionMatrix = mgl32.Perspective(mgl32.DegToRad(render.FOV), float32(render.WindowWidth)/float32(render.WindowHeight), 0.1, 10)
	} else {
		var size float32 = 2

		p.ProjectionMatrix = mgl32.Ortho(-size, size, -size/render.WindowAspect, size/render.WindowAspect, .1, 10.0)

	}
}

//DrawProject draws the project
func (p *Project) DrawProject() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	p.UpdateMatrices()
	for i := range p.Items {
		p.Items[i].DrawChildren(p.CamMatrix, p.ProjectionMatrix)
	}
}
func (p *Project) DrawIDs() {
	//Bind ID Framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, p.IDFramebuffer)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for i := range p.Items {
		p.Items[i].DrawChildrenIDs(p.CamMatrix, p.ProjectionMatrix)
	}

	//Go Back to default framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
func (p *Project) HandleInput(window *glfw.Window) {
	//Camera controls

	act := window.GetMouseButton(glfw.MouseButton3)
	if act == glfw.Press {
		//Get mouse delta
		//do some math

		d := imgui.CurrentIO().MouseDelta()
		fx := float32(render.WindowWidth)
		fy := float32(render.WindowHeight)
		p.CamPos = p.CamPos.Add(mgl32.Vec3{render.MouseSensitivity * -d.X / fx, render.MouseSensitivity * d.Y / fy, 0})
		log.Println(p.CamPos)
	}
	//Select All

	//Rest of it
	for i := range p.Items {
		p.Items[i].HandleInput()
	}
}
