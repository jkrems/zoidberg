package main

import (
	"fmt"
	"strconv"
)

type tokenFn func() token

type SyntaxError struct {
	Message  string
	Filename string
	Line     int
	Column   int
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s at %s:%d:%d", e.Message, e.Filename, e.Line, e.Column)
}

func syntaxError(tok token, format string, args ...interface{}) SyntaxError {
	return SyntaxError{
		fmt.Sprintf(format, args...),
		"unknown.berg",
		1,
		1,
	}
}

func unexpected(tok token) SyntaxError {
	return syntaxError(tok, "Unexpected %s", tok.typ)
}

func parseLExpr(read tokenFn) (*Identifier, error) {
	id := read()
	if id.typ != tokenIdentifier {
		return nil, unexpected(id)
	}
	return &Identifier{id.val}, nil
}

func parseDeclaration(read tokenFn) (ASTNode, error) {
	return nil, nil
}

func parseUnaryExpr(read tokenFn) (ASTNode, error) {
	tok := read()
	switch tok.typ {
	case tokenNumber:
		ival, err := strconv.ParseInt(tok.val, 0, 32)
		if err != nil {
			return nil, err
		}
		return &IntLiteralExpr{int(ival)}, nil
	}
	return nil, unexpected(tok)
}

func parseExpr(read tokenFn) (ASTNode, error) {
	left, err := parseUnaryExpr(read)
	// And now there's an optional next element...
	return left, err
}

func parseProgram(read tokenFn) (*Program, error) {
	init := []ASTNode{}
	for {
		lexpr, err := parseLExpr(read)
		if err != nil {
			return nil, err
		}
		op := read()
		if op.typ != tokenAssign {
			return nil, unexpected(op)
		}
		expr, err := parseExpr(read)
		if err != nil {
			return nil, err
		}
		init = append(init, Assignment{*lexpr, expr})
		break
	}
	return &Program{
		init,
	}, nil
}

func Parse(tokens chan token) (*Program, error) {
	read := func() token {
		return <-tokens
	}
	return parseProgram(read)
}
