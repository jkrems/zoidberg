package main

import "testing"

func verifyToken(t *testing.T, source string, verifier func(tok token) bool) {
	tokens := Tokenize("testfile.berg", source)

	found := false
	for tok := range tokens {
		found = found || verifier(tok)
	}
	if !found {
		t.Errorf("The expected token was not found")
	}
}

func verifyTokenValue(t *testing.T, source string, typ tokenType, val string) {
	verifyToken(t, source, func(tok token) bool {
		if tok.typ == typ {
			if tok.val != val {
				t.Errorf("Expected '%s', got %s", val, tok.val)
			}
			return true
		}
		return false
	})
}

func TestTokenizeInt(t *testing.T) {
	verifyTokenValue(t, "13", tokenNumber, "13")
}

func TestTokenizeFloat(t *testing.T) {
	verifyTokenValue(t, "42.531", tokenNumber, "42.531")
}

func TestTokenizeTwoInts(t *testing.T) {
	verifyToken(t, "42 53", func(tok token) bool {
		return tok.typ == tokenNumber && tok.val == "42"
	})
	verifyToken(t, "42 53", func(tok token) bool {
		return tok.typ == tokenNumber && tok.val == "53"
	})
}

func TestTokenizeTwoIntsWhitespace(t *testing.T) {
	verifyToken(t, "42\t \n  53", func(tok token) bool {
		return tok.typ == tokenNumber && tok.val == "42"
	})
	verifyToken(t, "42\t \n  53", func(tok token) bool {
		return tok.typ == tokenNumber && tok.val == "53"
	})
}

func TestTokenizeStringDoubleQuotes(t *testing.T) {
	verifyTokenValue(t, "\"abc\"", tokenString, "abc")
}

func TestTokenizeStringEscapedDoubleQuote(t *testing.T) {
	verifyTokenValue(t, "\"With \\\"stuff\\\"!\"", tokenString, "With \\\"stuff\\\"!")
}

func TestTokenizeIdentifier(t *testing.T) {
	verifyTokenValue(t, "x", tokenIdentifier, "x")
	verifyTokenValue(t, "10 y", tokenIdentifier, "y")
}

func TestTokenizeOperators(t *testing.T) {
	verifyTokenValue(t, "4 +5", tokenBinaryOrUnaryOp, "+")
	verifyTokenValue(t, "x = 14", tokenAssign, "=")
	verifyTokenValue(t, "x == 14", tokenBinaryOp, "==")
	verifyTokenValue(t, "x / &14", tokenBinaryOrUnaryOp, "&")
	verifyTokenValue(t, "x && 14", tokenBinaryOp, "&&")
	verifyTokenValue(t, "x && !y", tokenUnaryOp, "!")
	verifyTokenValue(t, "x ** !y", tokenBinaryOp, "**")
	verifyTokenValue(t, "*y", tokenBinaryOrUnaryOp, "*")
	verifyTokenValue(t, "a ++ b", tokenBinaryOp, "++")
}

func TestTokenizeDeclare(t *testing.T) {
	verifyTokenValue(t, "val x = 10", tokenDeclare, "val")
	verifyTokenValue(t, "var z = \"foo\"", tokenDeclare, "var")
}

func TestTokenizeBrackets(t *testing.T) {
	verifyTokenValue(t, "(", tokenLParen, "(")
	verifyTokenValue(t, ")", tokenRParen, ")")
	verifyTokenValue(t, "[", tokenLSquare, "[")
	verifyTokenValue(t, "]", tokenRSquare, "]")
	verifyTokenValue(t, "{", tokenLCurly, "{")
	verifyTokenValue(t, "}", tokenRCurly, "}")
}
