package token

import (
	"fmt"
)

type TokenType uint32

const (
	STRING TokenType = iota
	NUMBER
	IDENTIFIER
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	AND
	CONTRACT
	ABSTRACT
	INTERFACE
	LIB
	FUNCTION
	TRUE
	FALSE
	ELSE
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	VAR
	LET
	CONST
	WHILE
	ADDRESS
	PROXY
	UINT256
	UINT512
	BRIDGE
	PRAGMA
	VERSION
	SYNTAX
)

type Token struct {
	typ        Type
	lexeme     string
	lext, rext int
	literal    interface{}
	line       int
	input      []rune
}

func NewToken(typ TokenType, lext, rext int, input []rune) *Token {

	return &Token{
		typ, "", lext, rext, nil, 0, input,
	}

}

// GetLineColumn returns the line and column of the left extent of t
func (t *Token) GetLineColumn() (line, col int) {
	line, col = 1, 1
	for j := 0; j < t.lext; j++ {
		switch t.input[j] {
		case '\n':
			line++
			col = 1
		case '\t':
			col += 4
		default:
			col++
		}
	}
	return
}

// GetInput returns the input from which t was parsed.
func (t *Token) GetInput() []rune {
	return t.input
}

// Lext returns the left extent of t
func (t *Token) Lext() int {
	return t.lext
}

// Literal returns the literal runes of t scanned by the lexer
func (t *Token) Literal() []rune {
	return t.input[t.lext:t.rext]
}

// LiteralString returns string(t.Literal())
func (t *Token) LiteralString() string {
	return string(t.Literal())
}

// LiteralStripEscape returns the literal runes of t scanned by the lexer
func (t *Token) LiteralStripEscape() []rune {
	lit := t.Literal()
	strip := make([]rune, 0, len(lit))
	for i := 0; i < len(lit); i++ {
		if lit[i] == '\\' {
			i++
			switch lit[i] {
			case 't':
				strip = append(strip, '\t')
			case 'r':
				strip = append(strip, '\r')
			case 'n':
				strip = append(strip, '\r')
			default:
				strip = append(strip, lit[i])
			}
		} else {
			strip = append(strip, lit[i])
		}
	}
	return strip
}

// LiteralStringStripEscape returns string(t.LiteralStripEscape())
func (t *Token) LiteralStringStripEscape() string {
	return string(t.LiteralStripEscape())
}

// Rext returns the right extent of t in the input
func (t *Token) Rext() int {
	return t.rext
}

func (t *Token) String() string {
	return fmt.Sprintf("%s (%d,%d) %s",
		t.TypeID(), t.lext, t.rext, t.LiteralString())
}

// Suppress returns true iff t is suppressed by the lexer
func (t *Token) Suppress() bool {
	return Suppress[t.typ]
}

// Type returns the token Type of t
func (t *Token) Type() Type {
	return t.typ
}

// TypeID returns the token Type ID of t.
// This may be different from the literal of token t.
func (t *Token) TypeID() string {
	return t.Type().ID()
}
