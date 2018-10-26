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
	EQUAL
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[string]int{
	token.EQ:       EQUAL,
	token.NE:       EQUAL,
	token.LT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.OR:       SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.AND:      PRODUCT,
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
	for p.tk.Type != token.EOF {
		statements = append(statements, p.parseStmt())
	}
	return &ast.Program{Statements: statements}
}

func (p *parser) parseStmt() ast.Stmt {
	return p.parseExprStmt()
}

func (p *parser) parseExprStmt() *ast.ExprStmt {
	expr := p.parseExpr(LOWEST)
	if p.tk.Type != token.SEMICOLON {
		util.Error("Expected ; but got %s", p.tk.Literal)
	}
	p.next()
	return &ast.ExprStmt{Expr: expr}
}

func (p *parser) parseExpr(precedence int) ast.Expr {
	var expr ast.Expr
	switch p.tk.Type {
	case token.LPAREN:
		expr = p.parseGroupedExpr()
	case token.BANG:
		expr = p.parsePrefixExpr()
	case token.MINUS:
		expr = p.parsePrefixExpr()
	case token.NUMBER:
		expr = p.parseIntLit()
	case token.TRUE:
		expr = p.parseBoolLit()
	case token.FALSE:
		expr = p.parseBoolLit()
	default:
		util.Error("Unexpected %s", p.tk.Literal)
	}
	for p.lookPrecedence() > precedence {
		expr = p.parseInfixExpr(expr)
	}
	return expr
}

func (p *parser) parseGroupedExpr() ast.Expr {
	p.next()
	expr := p.parseExpr(LOWEST)
	if p.tk.Type != token.RPAREN {
		util.Error("Expected ) but got %s", p.tk.Literal)
	}
	p.next()
	return expr
}

func (p *parser) parsePrefixExpr() *ast.PrefixExpr {
	operator := p.tk.Literal
	p.next()
	right := p.parseExpr(PREFIX)
	return &ast.PrefixExpr{Operator: operator, Right: right}
}

func (p *parser) parseInfixExpr(left ast.Expr) *ast.InfixExpr {
	operator := p.tk.Literal
	precedence := p.lookPrecedence()
	p.next()
	right := p.parseExpr(precedence)
	return &ast.InfixExpr{Operator: operator, Left: left, Right: right}
}

func (p *parser) parseIntLit() *ast.IntLit {
	value, err := strconv.ParseInt(p.tk.Literal, 0, 64)
	if err != nil {
		util.Error("Could not parse %s as integer", p.tk.Literal)
	}
	p.next()
	return &ast.IntLit{Value: value}
}

func (p *parser) parseBoolLit() *ast.BoolLit {
	ty := p.tk.Type
	p.next()
	return &ast.BoolLit{Value: ty == token.TRUE}
}
