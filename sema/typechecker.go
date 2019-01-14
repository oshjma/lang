package sema

import (
	"github.com/oshjma/lang/ast"
	"github.com/oshjma/lang/types"
	"github.com/oshjma/lang/util"
)

/*
 Typechecker - do type checking
*/

type typechecker struct {
	refs  map[ast.Node]ast.Node
	types map[ast.Expr]types.Type
}

/* Program */

func (t *typechecker) typecheckProgram(prog *ast.Program) {
	for _, stmt := range prog.Stmts {
		t.typecheckStmt(stmt)
	}
}

/* Stmt */

func (t *typechecker) typecheckStmt(stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.BlockStmt:
		t.typecheckBlockStmt(v)
	case *ast.VarStmt:
		t.typecheckVarStmt(v)
	case *ast.FuncStmt:
		t.typecheckFuncStmt(v)
	case *ast.IfStmt:
		t.typecheckIfStmt(v)
	case *ast.ForStmt:
		t.typecheckForStmt(v)
	case *ast.ForInStmt:
		t.typecheckForInStmt(v)
	case *ast.ReturnStmt:
		t.typecheckReturnStmt(v)
	case *ast.AssignStmt:
		t.typecheckAssignStmt(v)
	case *ast.ExprStmt:
		t.typecheckExprStmt(v)
	}
}

func (t *typechecker) typecheckBlockStmt(stmt *ast.BlockStmt) {
	for _, stmt_ := range stmt.Stmts {
		t.typecheckStmt(stmt_)
	}
}

func (t *typechecker) typecheckVarStmt(stmt *ast.VarStmt) {
	for _, var_ := range stmt.Vars {
		t.typecheckVarDecl(var_)
	}
}

func (t *typechecker) typecheckFuncStmt(stmt *ast.FuncStmt) {
	t.typecheckFuncDecl(stmt.Func)
}

func (t *typechecker) typecheckIfStmt(stmt *ast.IfStmt) {
	t.typecheckExpr(stmt.Cond)
	ty := t.types[stmt.Cond]

	if _, ok := ty.(*types.Bool); !ok {
		util.Error("Expected bool value for if condition, but got %s", ty)
	}

	t.typecheckBlockStmt(stmt.Body)

	if stmt.Else != nil {
		t.typecheckStmt(stmt.Else)
	}
}

func (t *typechecker) typecheckForStmt(stmt *ast.ForStmt) {
	t.typecheckExpr(stmt.Cond)
	ty := t.types[stmt.Cond]

	if _, ok := ty.(*types.Bool); !ok {
		util.Error("Expected bool value for while condition, but got %s", ty)
	}

	t.typecheckBlockStmt(stmt.Body)
}

func (t *typechecker) typecheckForInStmt(stmt *ast.ForInStmt) {
	t.typecheckVarDecl(stmt.Array)
	ty := stmt.Array.VarType

	arr, ok := ty.(*types.Array)
	if !ok {
		util.Error("Expected array iterator, but got %s", ty)
	}

	stmt.Elem.VarType = arr.ElemType
	stmt.Index.VarType = &types.Int{}

	t.typecheckBlockStmt(stmt.Body)
}

func (t *typechecker) typecheckReturnStmt(stmt *ast.ReturnStmt) {
	ref := t.refs[stmt]

	var returnType types.Type

	switch v := ref.(type) {
	case *ast.FuncDecl:
		returnType = v.ReturnType
	case *ast.FuncLit:
		returnType = v.ReturnType
	}

	if stmt.Value == nil {
		if returnType != nil {
			f := "Expected %s return in function, but got nothing"
			util.Error(f, returnType)
		}
	} else {
		t.typecheckExpr(stmt.Value)
		ty := t.types[stmt.Value]

		if returnType == nil {
			f := "Expected no return in function, but got %s"
			util.Error(f, ty)
		}
		if !types.Same(ty, returnType) {
			f := "Expected %s return in function, but got %s"
			util.Error(f, returnType, ty)
		}
	}
}

func (t *typechecker) typecheckAssignStmt(stmt *ast.AssignStmt) {
	for i, target := range stmt.Targets {
		value := stmt.Values[i]

		t.typecheckExpr(target)
		t.typecheckExpr(value)
		tty := t.types[target]
		vty := t.types[value]

		if !types.Same(tty, vty) {
			f := "Expected %s value in assignment, but got %s"
			util.Error(f, tty, vty)
		}
	}
}

func (t *typechecker) typecheckExprStmt(stmt *ast.ExprStmt) {
	t.typecheckExpr(stmt.Expr)
}

/* Expr */

func (t *typechecker) typecheckExpr(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.PrefixExpr:
		t.typecheckPrefixExpr(v)
	case *ast.InfixExpr:
		t.typecheckInfixExpr(v)
	case *ast.IndexExpr:
		t.typecheckIndexExpr(v)
	case *ast.CallExpr:
		t.typecheckCallExpr(v)
	case *ast.LibCallExpr:
		t.typecheckLibCallExpr(v)
	case *ast.Ident:
		t.typecheckIdent(v)
	case *ast.IntLit:
		t.types[v] = &types.Int{}
	case *ast.BoolLit:
		t.types[v] = &types.Bool{}
	case *ast.StringLit:
		t.types[v] = &types.String{}
	case *ast.ArrayLit:
		t.typecheckArrayLit(v)
	case *ast.FuncLit:
		t.typecheckFuncLit(v)
	}
}

func (t *typechecker) typecheckPrefixExpr(expr *ast.PrefixExpr) {
	t.typecheckExpr(expr.Right)
	ty := t.types[expr.Right]

	switch expr.Op {
	case "!":
		if _, ok := ty.(*types.Bool); !ok {
			util.Error("Expected bool operand for !, but got %s", ty)
		}
		t.types[expr] = &types.Bool{}
	case "-":
		if _, ok := ty.(*types.Int); !ok {
			util.Error("Expected int operand for -, but got %s", ty)
		}
		t.types[expr] = &types.Int{}
	}
}

func (t *typechecker) typecheckInfixExpr(expr *ast.InfixExpr) {
	t.typecheckExpr(expr.Left)
	t.typecheckExpr(expr.Right)
	lty := t.types[expr.Left]
	rty := t.types[expr.Right]

	switch expr.Op {
	case "+", "-", "*", "/", "%":
		_, lok := lty.(*types.Int)
		_, rok := rty.(*types.Int)
		if !lok || !rok {
			f := "Expected int operands for %s, but got %s, %s"
			util.Error(f, expr.Op, lty, rty)
		}
		t.types[expr] = &types.Int{}
	case "==", "!=":
		if lty == nil || rty == nil {
			util.Error("Unexpected void operand for %s", expr.Op)
		}
		if !types.Same(lty, rty) {
			f := "Expected same type operands for %s, but got %s, %s"
			util.Error(f, expr.Op, lty, rty)
		}
		t.types[expr] = &types.Bool{}
	case "<", "<=", ">", ">=":
		_, lok := lty.(*types.Int)
		_, rok := rty.(*types.Int)
		if !lok || !rok {
			f := "Expected int operands for %s, but got %s, %s"
			util.Error(f, expr.Op, lty, rty)
		}
		t.types[expr] = &types.Bool{}
	case "&&", "||":
		_, lok := lty.(*types.Bool)
		_, rok := rty.(*types.Bool)
		if !lok || !rok {
			f := "Expected bool operands for %s, but got %s, %s"
			util.Error(f, expr.Op, lty, rty)
		}
		t.types[expr] = &types.Bool{}
	case "=":
		if lty == nil || rty == nil {
			util.Error("Unexpected void operand for %s", expr.Op)
		}
		if !types.Same(lty, rty) {
			f := "Expected same type operands for %s, but got %s, %s"
			util.Error(f, expr.Op, lty, rty)
		}
		t.types[expr] = lty
	}
}

func (t *typechecker) typecheckIndexExpr(expr *ast.IndexExpr) {
	t.typecheckExpr(expr.Left)
	lty := t.types[expr.Left]

	arr, ok := lty.(*types.Array)
	if !ok {
		util.Error("Expected array to index, but got %s", lty)
	}

	t.typecheckExpr(expr.Index)
	ity := t.types[expr.Index]

	if _, ok := ity.(*types.Int); !ok {
		util.Error("Expected int index, but got %s", ity)
	}

	t.types[expr] = arr.ElemType
}

func (t *typechecker) typecheckCallExpr(expr *ast.CallExpr) {
	t.typecheckExpr(expr.Left)
	ty := t.types[expr.Left]

	fn, ok := ty.(*types.Func)
	if !ok {
		util.Error("Expected function to call, but got %s", ty)
	}

	if len(expr.Params) != len(fn.ParamTypes) {
		f := "Wrong number of parameters (expected %d, given %d)"
		util.Error(f, len(fn.ParamTypes), len(expr.Params))
	}
	for i, param := range expr.Params {
		paramType := fn.ParamTypes[i]

		t.typecheckExpr(param)
		ty := t.types[param]

		if !types.Same(ty, paramType) {
			f := "Expected %s value for #%d parameter, but got %s"
			util.Error(f, paramType, i+1, ty)
		}
	}

	t.types[expr] = fn.ReturnType
}

func (t *typechecker) typecheckLibCallExpr(expr *ast.LibCallExpr) {
	for _, param := range expr.Params {
		t.typecheckExpr(param)
	}
	t.types[expr] = nil // FIXME
}

func (t *typechecker) typecheckIdent(expr *ast.Ident) {
	ref := t.refs[expr]

	switch v := ref.(type) {
	case *ast.VarDecl:
		t.types[expr] = v.VarType
	case *ast.FuncDecl:
		paramTypes := make([]types.Type, 0, 4)
		for _, param := range v.Params {
			paramTypes = append(paramTypes, param.VarType)
		}
		t.types[expr] = &types.Func{ParamTypes: paramTypes, ReturnType: v.ReturnType}
	}
}

func (t *typechecker) typecheckArrayLit(expr *ast.ArrayLit) {
	for _, elem := range expr.Elems {
		t.typecheckExpr(elem)
		ty := t.types[elem]

		if !types.Same(ty, expr.ElemType) {
			f := "Expected %s value for array element, but got %s"
			util.Error(f, expr.ElemType, ty)
		}
	}
	t.types[expr] = &types.Array{Len: expr.Len, ElemType: expr.ElemType}
}

func (t *typechecker) typecheckFuncLit(expr *ast.FuncLit) {
	t.typecheckBlockStmt(expr.Body)

	paramTypes := make([]types.Type, 0, 4)
	for _, param := range expr.Params {
		paramTypes = append(paramTypes, param.VarType)
	}
	t.types[expr] = &types.Func{ParamTypes: paramTypes, ReturnType: expr.ReturnType}
}

/* Decl */

func (t *typechecker) typecheckVarDecl(decl *ast.VarDecl) {
	switch v := decl.Value.(type) {
	case *ast.FuncLit:
		paramTypes := make([]types.Type, 0, 4)
		for _, param := range v.Params {
			paramTypes = append(paramTypes, param.VarType)
		}
		ty := &types.Func{ParamTypes: paramTypes, ReturnType: v.ReturnType}
		t.types[v] = ty

		if decl.VarType == nil {
			decl.VarType = ty // type inference (write on AST node)
		} else {
			if !types.Same(ty, decl.VarType) {
				f := "Expected %s value for %s, but got %s"
				util.Error(f, decl.VarType, decl.Name, ty)
			}
		}
		t.typecheckBlockStmt(v.Body)
	default:
		t.typecheckExpr(v)
		ty := t.types[v]

		if decl.VarType == nil {
			if ty == nil {
				util.Error("No initial value for %s", decl.Name)
			}
			decl.VarType = ty // type inference (write on AST node)
		} else {
			if ty == nil {
				f := "Expected %s value for %s, but got nothing"
				util.Error(f, decl.VarType, decl.Name)
			}
			if !types.Same(ty, decl.VarType) {
				f := "Expected %s value for %s, but got %s"
				util.Error(f, decl.VarType, decl.Name, ty)
			}
		}
	}
}

func (t *typechecker) typecheckFuncDecl(decl *ast.FuncDecl) {
	t.typecheckBlockStmt(decl.Body)
}
