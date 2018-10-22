package parse

import (
	"strconv"
	"github.com/oshjma/lang/ast"
	"github.com/oshjma/lang/token"
	"github.com/oshjma/lang/util"
)

func Parse(tokens []*token.Token) *ast.Program {
	p := &parser{tokens: tokens, pos: -1}
	p.next()
	return p.parseProgram()
}

const (
	LOWEST int = iota
	SUM
	PRODUCT
)

var precedences = map[string]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

type parser struct {
	tokens []*token.Token // input tokens
	pos int               // current position
	tk *token.Token       // current token
}

func (p *parser) next() {
	p.pos += 1
	p.tk = p.tokens[p.pos]
}

func (p *parser) lookPrecedence() int {
	if pr, ok := precedences[p.tk.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *parser) parseProgram() *ast.Program {
	var statements []ast.Stmt
	var stmt ast.Stmt
	for p.tk.Type != token.EOF {
		stmt = p.parseStmt()
		statements = append(statements, stmt)
	}
	return &ast.Program{Statements: statements}
}

func (p *parser) parseStmt() ast.Stmt {
	return p.parseExprStmt()
}

func (p *parser) parseExprStmt() *ast.ExprStmt {
	expr := p.parseExpr(LOWEST)
	if p.tk.Type != token.SEMICOLON {
		if p.tk.Type == token.EOF {
			util.Error("Expected %q but got <EOF>", ";")
		} else {
			util.Error("Expected %q but got %q", ";", p.tk.Source)
		}
	}
	p.next()
	return &ast.ExprStmt{Expr: expr}
}

func (p *parser) parseExpr(precedence int) ast.Expr {
	var left ast.Expr
	switch p.tk.Type {
	case token.LPAREN:
		left = p.parseGroupedExpr()
	case token.INT:
		left = p.parseIntLit()
	case token.EOF:
		util.Error("Unexpected <EOF>")
	default:
		util.Error("Unexpected %q", p.tk.Source)
	}
	for p.lookPrecedence() > precedence {
		left = p.parseInfixExpr(left)
	}
	return left
}

func (p *parser) parseGroupedExpr() ast.Expr {
	p.next()
	expr := p.parseExpr(LOWEST)
	if p.tk.Type != token.RPAREN {
		util.Error("Expected %q but got %q", ")", p.tk.Source)
	}
	p.next()
	return expr
}

func (p *parser) parseInfixExpr(left ast.Expr) *ast.InfixExpr {
	operator := p.tk.Source
	precedence := p.lookPrecedence()
	p.next()
	right := p.parseExpr(precedence)
	return &ast.InfixExpr{Operator: operator, Left: left, Right: right}
}

func (p *parser) parseIntLit() *ast.IntLit {
	value, err := strconv.ParseInt(p.tk.Source, 0, 64)
	if err != nil {
		util.Error("Could not parse %q as integer", p.tk.Source)
	}
	p.next()
	return &ast.IntLit{Value: value}
}