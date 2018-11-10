package scan

import "github.com/oshjma/lang/token"

var punctuations = map[byte]string{
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'+': token.PLUS,
	'*': token.ASTERISK,
	'/': token.SLASH,
	',': token.COMMA,
	';': token.SEMICOLON,
}

var keywords = map[string]string{
	"func":     token.FUNC,
	"var":      token.VAR,
	"if":       token.IF,
	"else":     token.ELSE,
	"while":    token.WHILE,
	"return":   token.RETURN,
	"continue": token.CONTINUE,
	"break":    token.BREAK,
	"int":      token.INT,
	"bool":     token.BOOL,
	"true":     token.TRUE,
	"false":    token.FALSE,
}

var exprTerminators = []string{
	token.RPAREN,
	token.IDENT,
	token.NUMBER,
	token.TRUE,
	token.FALSE,
}
