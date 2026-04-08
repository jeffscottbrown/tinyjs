package ast

// Program is the root node for the tiny language.
type Program struct {
	Assignments []*Assignment `@@*`
}

// Assignment represents the only supported statement shape:
//
// 	name = 123;
//
// This is intentionally tiny so the project can grow live.
type Assignment struct {
	Name  string `@Ident`
	Eq    string `"="`
	Value int64  `@Int`
	Semi  string `";"`
}
