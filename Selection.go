package main

type SelectableType int

const (
	VertexType SelectableType = iota
	LineType
	FaceType
	SolidType
	BodyType
)

//Holds Primitives for selection
type Selectable interface {
	Path() string
	Type() SelectableType
}

type Vertex struct {
	path string
}

func (v *Vertex) Path() string {
	return v.path
}
func (v *Vertex) Type() SelectableType {
	return VertexType
}

type Line struct {
	path string
}

func (l *Line) Path() string {
	return l.path
}
func (l *Line) Type() SelectableType {
	return LineType
}

type Face struct {
	path string
}

func (f *Face) Path() string {
	return f.path
}
func (f *Face) Type() SelectableType {
	return FaceType
}
