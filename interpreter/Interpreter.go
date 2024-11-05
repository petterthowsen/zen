package interpreter

import (
	"fmt"
	"zen/lang/common"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
	"zen/runtime/environment"
	"zen/runtime/types"
)

// Interpreter coordinates the execution of Zen programs
type Interpreter struct {
	// The current execution environment
	env *environment.Environment
	// Track whether we're in a function
	inFunction bool
	// Track whether we're in a loop
	inLoop bool
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter() *Interpreter {
	return &Interpreter{
		env:        environment.NewEnvironment(),
		inFunction: false,
		inLoop:     false,
	}
}

// Execute runs a complete Zen program
func (i *Interpreter) Execute(program *ast.ProgramNode) error {
	for _, stmt := range program.Statements {
		if err := i.ExecuteStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteStatement executes a single statement
func (i *Interpreter) ExecuteStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *statement.VarDeclarationNode:
		return i.executeVarDeclaration(s)
	case *statement.ExpressionStatement:
		return i.executeExpressionStatement(s)
	case *statement.IfStatement:
		return i.executeIfStatement(s)
	case *statement.WhileStatement:
		return i.executeWhileStatement(s)
	default:
		return &RuntimeError{
			Message:  "Unknown statement type",
			Location: stmt.GetLocation(),
		}
	}
}

// GetValue retrieves a variable's value from the current environment
func (i *Interpreter) GetValue(name string) (interface{}, error) {
	val, err := i.env.Get(name)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// RuntimeError represents an error that occurs during program execution
type RuntimeError struct {
	Message  string
	Location *common.SourceLocation
}

func (e *RuntimeError) Error() string {
	if e.Location != nil {
		return fmt.Sprintf("Runtime error at %s: %s", e.Location.String(), e.Message)
	}
	return fmt.Sprintf("Runtime error: %s", e.Message)
}

// EvaluateExpression evaluates an expression and returns a Zen value
func (i *Interpreter) EvaluateExpression(expr ast.Expression) (types.Value, error) {
	switch e := expr.(type) {
	case *expression.LiteralExpression:
		return i.evaluateLiteral(e)
	case *expression.IdentifierExpression:
		return i.evaluateIdentifier(e)
	case *expression.UnaryExpression:
		return i.evaluateUnary(e)
	case *expression.BinaryExpression:
		return i.evaluateBinary(e)
	default:
		return nil, &RuntimeError{
			Message:  "Unknown expression type",
			Location: expr.GetLocation(),
		}
	}
}

// evaluateLiteral handles literal values (numbers, strings, booleans)
func (i *Interpreter) evaluateLiteral(expr *expression.LiteralExpression) (types.Value, error) {
	val, err := types.FromGoValue(expr.Value)
	if err != nil {
		return nil, &RuntimeError{
			Message:  err.Error(),
			Location: expr.GetLocation(),
		}
	}
	return val, nil
}

// evaluateIdentifier handles variable references
func (i *Interpreter) evaluateIdentifier(expr *expression.IdentifierExpression) (types.Value, error) {
	value, err := i.env.Get(expr.Name)
	if err != nil {
		return nil, &RuntimeError{
			Message:  fmt.Sprintf("Undefined variable '%s'", expr.Name),
			Location: expr.GetLocation(),
		}
	}

	// Convert Go value to Zen value
	val, err := types.FromGoValue(value)
	if err != nil {
		return nil, &RuntimeError{
			Message:  err.Error(),
			Location: expr.GetLocation(),
		}
	}
	return val, nil
}

// evaluateUnary handles unary operations (-x, not x)
func (i *Interpreter) evaluateUnary(expr *expression.UnaryExpression) (types.Value, error) {
	operand, err := i.EvaluateExpression(expr.Expression)
	if err != nil {
		return nil, err
	}

	result, err := types.UnaryOp(operand, expr.Operator)
	if err != nil {
		return nil, &RuntimeError{
			Message:  err.Error(),
			Location: expr.GetLocation(),
		}
	}
	return result, nil
}

// evaluateBinary handles binary operations (+, -, *, /, etc.)
func (i *Interpreter) evaluateBinary(expr *expression.BinaryExpression) (types.Value, error) {
	// Handle assignment separately since right side should only be evaluated if needed
	if expr.Operator == "=" {
		if id, ok := expr.Left.(*expression.IdentifierExpression); ok {
			right, err := i.EvaluateExpression(expr.Right)
			if err != nil {
				return nil, err
			}
			if err := i.env.Assign(id.Name, types.ToGoValue(right)); err != nil {
				return nil, &RuntimeError{
					Message:  err.Error(),
					Location: expr.GetLocation(),
				}
			}
			return right, nil
		}
		return nil, &RuntimeError{
			Message:  "Invalid assignment target",
			Location: expr.GetLocation(),
		}
	}

	// Handle short-circuit evaluation for logical operators
	if expr.Operator == "and" || expr.Operator == "or" {
		left, err := i.EvaluateExpression(expr.Left)
		if err != nil {
			return nil, err
		}

		// Check that left operand is boolean
		if left.Type() != types.TypeBool {
			return nil, &RuntimeError{
				Message:  fmt.Sprintf("Left operand of %s must be boolean, got %s", expr.Operator, left.Type()),
				Location: expr.Left.GetLocation(),
			}
		}

		// Short-circuit evaluation
		if expr.Operator == "and" {
			if !left.IsTruthy() {
				return types.NewBool(false), nil
			}
		} else { // or
			if left.IsTruthy() {
				return types.NewBool(true), nil
			}
		}

		right, err := i.EvaluateExpression(expr.Right)
		if err != nil {
			return nil, err
		}

		// Check that right operand is boolean
		if right.Type() != types.TypeBool {
			return nil, &RuntimeError{
				Message:  fmt.Sprintf("Right operand of %s must be boolean, got %s", expr.Operator, right.Type()),
				Location: expr.Right.GetLocation(),
			}
		}

		result, err := types.BinaryOp(left, right, expr.Operator)
		if err != nil {
			return nil, &RuntimeError{
				Message:  err.Error(),
				Location: expr.GetLocation(),
			}
		}
		return result, nil
	}

	// Evaluate both operands for other operators
	left, err := i.EvaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.EvaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	result, err := types.BinaryOp(left, right, expr.Operator)
	if err != nil {
		return nil, &RuntimeError{
			Message:  err.Error(),
			Location: expr.GetLocation(),
		}
	}
	return result, nil
}

// executeVarDeclaration handles variable and constant declarations
func (i *Interpreter) executeVarDeclaration(stmt *statement.VarDeclarationNode) error {
	var value interface{}
	var err error

	// Evaluate initializer if present
	if stmt.Initializer != nil {
		val, err := i.EvaluateExpression(stmt.Initializer)
		if err != nil {
			return err
		}
		value = types.ToGoValue(val)
	}

	// If no initializer and not nullable, that's an error
	if stmt.Initializer == nil && !stmt.IsNullable {
		return &RuntimeError{
			Message:  fmt.Sprintf("Variable '%s' must be initialized or marked as nullable", stmt.Name),
			Location: stmt.GetLocation(),
		}
	}

	// Define the variable based on its properties
	if stmt.IsConstant {
		if stmt.Initializer == nil {
			return &RuntimeError{
				Message:  fmt.Sprintf("Constant '%s' must be initialized", stmt.Name),
				Location: stmt.GetLocation(),
			}
		}
		err = i.env.DefineConst(stmt.Name, value)
	} else if stmt.IsNullable {
		err = i.env.DefineNullable(stmt.Name, value)
	} else {
		err = i.env.Define(stmt.Name, value)
	}

	if err != nil {
		return &RuntimeError{
			Message:  err.Error(),
			Location: stmt.GetLocation(),
		}
	}

	return nil
}

// executeExpressionStatement handles expressions used as statements
func (i *Interpreter) executeExpressionStatement(stmt *statement.ExpressionStatement) error {
	_, err := i.EvaluateExpression(stmt.Expression)
	return err
}

// executeIfStatement handles if statements with else-if and else blocks
func (i *Interpreter) executeIfStatement(stmt *statement.IfStatement) error {
	// Evaluate primary condition
	condition, err := i.EvaluateExpression(stmt.PrimaryCondition)
	if err != nil {
		return err
	}

	// Check that condition is a boolean
	if condition.Type() != types.TypeBool {
		return &RuntimeError{
			Message:  fmt.Sprintf("If condition must be a boolean, got %s", condition.Type()),
			Location: stmt.GetLocation(),
		}
	}

	if condition.IsTruthy() {
		// Execute primary block in new scope
		i.env.BeginScope()
		defer i.env.EndScope()

		for _, statement := range stmt.PrimaryBlock {
			if err := i.ExecuteStatement(statement); err != nil {
				return err
			}
		}
		return nil
	}

	// Check else-if blocks
	for _, elseIfBlock := range stmt.ElseIfBlocks {
		condition, err := i.EvaluateExpression(elseIfBlock.Condition)
		if err != nil {
			return err
		}

		// Check that condition is a boolean
		if condition.Type() != types.TypeBool {
			return &RuntimeError{
				Message:  fmt.Sprintf("Elif condition must be a boolean, got %s", condition.Type()),
				Location: elseIfBlock.GetLocation(),
			}
		}

		if condition.IsTruthy() {
			// Execute else-if block in new scope
			i.env.BeginScope()
			defer i.env.EndScope()

			for _, statement := range elseIfBlock.Body {
				if err := i.ExecuteStatement(statement); err != nil {
					return err
				}
			}
			return nil
		}
	}

	// If no conditions were true and there's an else block, execute it
	if len(stmt.ElseBlock) > 0 {
		i.env.BeginScope()
		defer i.env.EndScope()

		for _, statement := range stmt.ElseBlock {
			if err := i.ExecuteStatement(statement); err != nil {
				return err
			}
		}
	}

	return nil
}

// executeWhileStatement handles while loops
func (i *Interpreter) executeWhileStatement(stmt *statement.WhileStatement) error {
	// Set loop state
	wasInLoop := i.inLoop
	i.inLoop = true
	defer func() { i.inLoop = wasInLoop }()

	for {
		// Evaluate condition
		condition, err := i.EvaluateExpression(stmt.Condition)
		if err != nil {
			return err
		}

		// Check that condition is a boolean
		if condition.Type() != types.TypeBool {
			return &RuntimeError{
				Message:  fmt.Sprintf("While condition must be a boolean, got %s", condition.Type()),
				Location: stmt.GetLocation(),
			}
		}

		if !condition.IsTruthy() {
			break
		}

		// Execute body in new scope
		i.env.BeginScope()
		for _, statement := range stmt.Body {
			if err := i.ExecuteStatement(statement); err != nil {
				i.env.EndScope()
				return err
			}
		}
		i.env.EndScope()
	}

	return nil
}
