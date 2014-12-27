package main

type ASTNode interface {
	GetName() string
}

type Identifier struct {
	Name string
}

func (id Identifier) GetName() string {
	return "Identifier"
}

type IntLiteralExpr struct {
	Value int
}

func (i IntLiteralExpr) GetName() string {
	return "IntLiteralExpr"
}

type Assignment struct {
	Target Identifier
	Value  ASTNode
}

func (a Assignment) GetName() string {
	return "Assignment"
}

type Program struct {
	Init []ASTNode
}

func (p Program) GetName() string {
	return "Program"
}
