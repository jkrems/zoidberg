package main

import (
	"strings"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenNumber
	tokenString
	tokenIdentifier
	tokenBinaryOp
	tokenBinaryOrUnaryOp
	tokenUnaryOp
	tokenAssign
	tokenDeclare
	tokenLParen
	tokenRParen
	tokenLSquare
	tokenRSquare
	tokenLCurly
	tokenRCurly
)
const eof = -1

type token struct {
	typ tokenType
	val string
}

type lexer struct {
	filename string
	source   string
	start    int
	width    int
	pos      int
	tokens   chan token
}

type stateFn func(l *lexer) stateFn

var simpleTokens = map[rune]tokenType{
	'(': tokenLParen,
	')': tokenRParen,
	'[': tokenLSquare,
	']': tokenRSquare,
	'{': tokenLCurly,
	'}': tokenRCurly,
	'-': tokenBinaryOrUnaryOp,
	'/': tokenBinaryOp,
}

func lexNumber(l *lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	l.emit(tokenNumber)
	return lexRoot
}

func lexDoubleQuotedString(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case '\\':
			l.next()
		case '"':
			l.backup()
			l.emit(tokenString)
			l.next()
			return lexRoot
		case eof:
			l.emit(tokenError)
			return nil
		}
	}
}

func isIdentifierStart(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isIdentifierPart(r rune) bool {
	return isIdentifierStart(r)
}

func lexIdentifier(l *lexer) stateFn {
	for isIdentifierPart(l.next()) {
	}
	l.backup()
	switch id := l.current(); {
	case id == "val" || id == "var":
		l.emit(tokenDeclare)
	default:
		l.emit(tokenIdentifier)
	}
	return lexRoot
}

func repeatOperator(l *lexer, op rune, singleType, repeatType tokenType) {
	peek := l.next()
	if peek == op {
		l.emit(repeatType)
	} else {
		l.backup()
		l.emit(singleType)
	}
}

func lexRoot(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(tokenEOF)
			return nil
		case ' ' == r || '\t' == r || '\n' == r:
			l.ignore()
		case '"' == r:
			l.ignore()
			return lexDoubleQuotedString
		case isIdentifierStart(r):
			l.backup()
			return lexIdentifier
		case r == '-' || r == '/':
			l.emit(tokenBinaryOp)
		case r == '+':
			repeatOperator(l, r, tokenBinaryOrUnaryOp, tokenBinaryOp)
		case r == '=':
			repeatOperator(l, r, tokenAssign, tokenBinaryOp)
		case r == '&':
			repeatOperator(l, r, tokenBinaryOrUnaryOp, tokenBinaryOp)
		case r == '|':
			repeatOperator(l, r, tokenBinaryOp, tokenBinaryOp)
		case r == '!':
			l.emit(tokenUnaryOp)
		case r == '*':
			repeatOperator(l, r, tokenBinaryOrUnaryOp, tokenBinaryOp)
		case '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		default:
			typ, ok := simpleTokens[r]
			if !ok {
				panic("Unexpected character")
				return nil
			}
			l.emit(typ)
		}
	}
	panic("Should never be reached")
}

func (l *lexer) emit(typ tokenType) {
	l.tokens <- token{typ, l.source[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) current() string {
	return l.source[l.start:l.pos]
}

func (l *lexer) next() (rune rune) {
	if l.pos >= len(l.source) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.source[l.pos:])
	l.pos += l.width
	return rune
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) run() {
	for state := lexRoot; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func Tokenize(filename string, source string) (*lexer, chan token) {
	l := &lexer{
		filename: filename,
		source:   source,
		tokens:   make(chan token),
	}
	go l.run()
	return l, l.tokens
}
