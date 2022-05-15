package plan

import (
	"errors"

	"github.com/goropikari/simpledbgo/domain"
)

// BasicUpdatePlanner is a BasicUpdatePlanner.
type BasicUpdatePlanner struct {
	metadataMgr domain.MetadataManager
}

// NewBasicUpdatePlanner constructs a BasicUpdatePlanner.
func NewBasicUpdatePlanner(mgr domain.MetadataManager) *BasicUpdatePlanner {
	return &BasicUpdatePlanner{metadataMgr: mgr}
}

// ExecuteInsert executes insersion command.
func (p *BasicUpdatePlanner) ExecuteInsert(data *domain.InsertData, txn domain.Transaction) (int, error) {
	var plan domain.Planner
	plan, err := NewTablePlan(txn, data.TableName(), p.metadataMgr)
	if err != nil {
		return 0, err
	}

	s, err := plan.Open()
	if err != nil {
		return 0, err
	}

	us, ok := s.(domain.UpdateScanner)
	if !ok {
		return 0, errors.New("not updatable")
	}

	if err = us.AdvanceNextInsertSlotID(); err != nil {
		return 0, err
	}

	vals := data.Values()
	fields := data.Fields()
	for i := 0; i < len(vals); i++ {
		if err = us.SetVal(fields[i], vals[i]); err != nil {
			return 0, err
		}
	}
	us.Close()

	return 1, nil
}

// ExecuteDelete executes delete command.
func (p *BasicUpdatePlanner) ExecuteDelete(data *domain.DeleteData, txn domain.Transaction) (int, error) {
	var plan domain.Planner
	plan, err := NewTablePlan(txn, data.TableName(), p.metadataMgr)
	if err != nil {
		return 0, err
	}

	plan = NewSelectPlan(plan, data.Predicate())
	s, err := plan.Open()
	if err != nil {
		return 0, err
	}

	us, ok := s.(domain.UpdateScanner)
	if !ok {
		return 0, errors.New("not updatable")
	}

	cnt := 0
	for us.HasNext() {
		if err = us.Delete(); err != nil {
			return 0, err
		}
		cnt++
	}
	if us.Err() != nil {
		return 0, us.Err()
	}
	us.Close()

	return cnt, nil
}

// ExecuteModify executes update command.
func (p *BasicUpdatePlanner) ExecuteModify(data *domain.ModifyData, txn domain.Transaction) (int, error) {
	var plan domain.Planner
	plan, err := NewTablePlan(txn, data.TableName(), p.metadataMgr)
	if err != nil {
		return 0, err
	}
	plan = NewSelectPlan(plan, data.Predicate())

	s, err := plan.Open()
	if err != nil {
		return 0, err
	}

	us, ok := s.(domain.UpdateScanner)
	if !ok {
		return 0, errors.New("not updatable")
	}

	cnt := 0
	for us.HasNext() {
		val, err := data.Expression().Evaluate(us)
		if err != nil {
			return 0, err
		}

		if err = us.SetVal(data.FieldName(), val); err != nil {
			return 0, err
		}
		cnt++
	}
	us.Close()

	return cnt, nil
}

// ExecuteCreateTable executes create table command.
func (p *BasicUpdatePlanner) ExecuteCreateTable(data *domain.CreateTableData, txn domain.Transaction) (int, error) {
	return 0, p.metadataMgr.CreateTable(data.TableName(), data.Schema(), txn)
}

// ExecuteCreateView executes create view command.
func (p *BasicUpdatePlanner) ExecuteCreateView(data *domain.CreateViewData, txn domain.Transaction) (int, error) {
	return 0, p.metadataMgr.CreateView(data.ViewName(), data.ViewDef(), txn)
}

// ExecuteCreateIndex executes create index command.
func (p *BasicUpdatePlanner) ExecuteCreateIndex(data *domain.CreateIndexData, txn domain.Transaction) (int, error) {
	return 0, p.metadataMgr.CreateIndex(data.IndexName(), data.TableName(), data.FieldName(), txn)
}
