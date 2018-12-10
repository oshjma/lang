package ast

/*
 Interfaces
*/

// for all AST nodes
type Node interface {
	astNode()
}

// for statement nodes
type Stmt interface {
	Node
	stmtNode()
}

// for expression nodes
type Expr interface {
	Node
	exprNode()
}

/*
 Root node
*/

type Program struct {
	Stmts []Stmt
}

func (prog *Program) astNode() {}

/*
 Statement nodes
*/

type FuncDecl struct {
	Ident      string
	Params     []*VarDecl
	ReturnType string
	Body       *BlockStmt
}

type VarDecl struct {
	Ident string
	Type  string
	Value Expr
}

type BlockStmt struct {
	Stmts []Stmt
}

type IfStmt struct {
	Cond   Expr
	Conseq *BlockStmt
	Altern Stmt // *BlockStmt or *IfStmt
}

type ForStmt struct {
	Cond Expr
	Body *BlockStmt
}

type ReturnStmt struct {
	Value Expr
}

type ContinueStmt struct {
	dummy byte
}

type BreakStmt struct {
	dummy byte
}

type AssignStmt struct {
	Ident string
	Value Expr
}

type ExprStmt struct {
	Expr Expr
}

func (stmt *VarDecl) astNode()       {}
func (stmt *VarDecl) stmtNode()      {}
func (stmt *FuncDecl) astNode()      {}
func (stmt *FuncDecl) stmtNode()     {}
func (stmt *BlockStmt) astNode()     {}
func (stmt *BlockStmt) stmtNode()    {}
func (stmt *IfStmt) astNode()        {}
func (stmt *IfStmt) stmtNode()       {}
func (stmt *ForStmt) astNode()       {}
func (stmt *ForStmt) stmtNode()      {}
func (stmt *ReturnStmt) astNode()    {}
func (stmt *ReturnStmt) stmtNode()   {}
func (stmt *ContinueStmt) astNode()  {}
func (stmt *ContinueStmt) stmtNode() {}
func (stmt *BreakStmt) astNode()     {}
func (stmt *BreakStmt) stmtNode()    {}
func (stmt *AssignStmt) astNode()    {}
func (stmt *AssignStmt) stmtNode()   {}
func (stmt *ExprStmt) astNode()      {}
func (stmt *ExprStmt) stmtNode()     {}

/*
 Expression nodes
*/

type PrefixExpr struct {
	Op    string
	Right Expr
}

type InfixExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

type FuncCall struct {
	Ident  string
	Params []Expr
}

type VarRef struct {
	Ident string
}

type IntLit struct {
	Value int64
}

type BoolLit struct {
	Value bool
}

type StringLit struct {
	Value string
}

func (expr *PrefixExpr) astNode()  {}
func (expr *PrefixExpr) exprNode() {}
func (expr *InfixExpr) astNode()   {}
func (expr *InfixExpr) exprNode()  {}
func (expr *FuncCall) astNode()    {}
func (expr *FuncCall) exprNode()   {}
func (expr *VarRef) astNode()      {}
func (expr *VarRef) exprNode()     {}
func (expr *IntLit) astNode()      {}
func (expr *IntLit) exprNode()     {}
func (expr *BoolLit) astNode()     {}
func (expr *BoolLit) exprNode()    {}
func (expr *StringLit) astNode()   {}
func (expr *StringLit) exprNode()  {}