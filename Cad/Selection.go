package cad

//SelectableType holds different types of things that can be selected
type SelectableType int

//All the types of SelectableType
const (
	VertexType SelectableType = iota
	LineType
	FaceType
	SolidType
	BodyType
)

//Selectable holds cad primitives for selection
type Selectable interface {
	Path() string
	Type() SelectableType
}
