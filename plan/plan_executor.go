package plan

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lexer"
	"github.com/goropikari/simpledbgo/parser"
)

// Executor is interface of plan executor.
type Executor struct {
	queryPlanner   domain.QueryPlanner
	updateExecutor domain.UpdateExecutor
}

// NewExecutor constructs a Executor.
func NewExecutor(qp domain.QueryPlanner, ue domain.UpdateExecutor) *Executor {
	return &Executor{
		queryPlanner:   qp,
		updateExecutor: ue,
	}
}

// CreateQueryPlan creates planner for select query.
func (pe Executor) CreateQueryPlan(query string, txn domain.Transaction) (domain.Planner, error) {
	lex := lexer.NewLexer(query)
	tokens, err := lex.ScanTokens()
	if err != nil {
		return nil, errors.Err(err, "ScanTokens")
	}

	parser := parser.NewParser(tokens)
	data, err := parser.Query()
	if err != nil {
		return nil, errors.Err(err, "Query")
	}

	err = pe.verifyQuery(data)
	if err != nil {
		return nil, errors.Err(err, "verifyQuery")
	}

	return pe.queryPlanner.CreatePlan(data, txn)
}

// ExecuteUpdate executes command.
func (pe Executor) ExecuteUpdate(cmd string, txn domain.Transaction) (int, error) {
	lex := lexer.NewLexer(cmd)
	tokens, err := lex.ScanTokens()
	if err != nil {
		return 0, errors.Err(err, "ScanTokens")
	}

	parser := parser.NewParser(tokens)
	data, err := parser.ExecCmd()
	if err != nil {
		return 0, errors.Err(err, "ExecCmd")
	}

	err = pe.verifyExec(data)
	if err != nil {
		return 0, errors.Err(err, "verifyExec")
	}

	switch v := data.(type) {
	case *domain.InsertData:
		return pe.updateExecutor.ExecuteInsert(v, txn)
	case *domain.DeleteData:
		return pe.updateExecutor.ExecuteDelete(v, txn)
	case *domain.ModifyData:
		return pe.updateExecutor.ExecuteModify(v, txn)
	case *domain.CreateTableData:
		return pe.updateExecutor.ExecuteCreateTable(v, txn)
	case *domain.CreateViewData:
		return pe.updateExecutor.ExecuteCreateView(v, txn)
	case *domain.CreateIndexData:
		return pe.updateExecutor.ExecuteCreateIndex(v, txn)
	default:
		return 0, errors.New("must not reach here")
	}
}

// SimpleDB does not verify queries, although it should.
func (pe Executor) verifyQuery(data *domain.QueryData) error {
	return nil
}

// SimpleDB does not verify execs, although it should.
func (pe Executor) verifyExec(data domain.ExecData) error {
	return nil
}
