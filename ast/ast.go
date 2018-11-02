package ast

/*
 Interfaces
*/

// for all AST nodes
type Node interface {
	astNode()
}

// for all statement nodes
type Stmt interface {
	Node
	stmtNode()
}

// for all expression nodes
type Expr interface {
	Node
	exprNode()
}

/*
 Root node
*/

type Program struct {
	TopLevel *BlockStmt
}

func (node *Program) astNode() {}

/*
 Statement nodes
*/

type BlockStmt struct {
	List []Stmt
}

type IfStmt struct {
	Cond   Expr
	Conseq *BlockStmt
	Altern Stmt // *BlockStmt or *IfStmt
}

type WhileStmt struct {
	Cond Expr
	Body *BlockStmt
}

type ContinueStmt struct{}

type BreakStmt struct{}

type LetStmt struct {
	Ident *Ident
	Type  string
	Expr  Expr
}

type AssignStmt struct {
	Ident *Ident
	Expr  Expr
}

type ExprStmt struct {
	Expr Expr
}

func (stmt *BlockStmt) astNode()     {}
func (stmt *BlockStmt) stmtNode()    {}
func (stmt *IfStmt) astNode()        {}
func (stmt *IfStmt) stmtNode()       {}
func (stmt *WhileStmt) astNode()     {}
func (stmt *WhileStmt) stmtNode()    {}
func (stmt *ContinueStmt) astNode()  {}
func (stmt *ContinueStmt) stmtNode() {}
func (stmt *BreakStmt) astNode()     {}
func (stmt *BreakStmt) stmtNode()    {}
func (stmt *LetStmt) astNode()       {}
func (stmt *LetStmt) stmtNode()      {}
func (stmt *AssignStmt) astNode()    {}
func (stmt *AssignStmt) stmtNode()   {}
func (stmt *ExprStmt) astNode()      {}
func (stmt *ExprStmt) stmtNode()     {}

/*
 Expression nodes
*/

type PrefixExpr struct {
	Operator string
	Right    Expr
}

type InfixExpr struct {
	Operator string
	Left     Expr
	Right    Expr
}

type Ident struct {
	Name string
}

type IntLit struct {
	Value int64
}

type BoolLit struct {
	Value bool
}

func (expr *PrefixExpr) astNode()  {}
func (expr *PrefixExpr) exprNode() {}
func (expr *InfixExpr) astNode()   {}
func (expr *InfixExpr) exprNode()  {}
func (expr *Ident) astNode()       {}
func (expr *Ident) exprNode()      {}
func (expr *IntLit) astNode()      {}
func (expr *IntLit) exprNode()     {}
func (expr *BoolLit) astNode()     {}
func (expr *BoolLit) exprNode()    {}
