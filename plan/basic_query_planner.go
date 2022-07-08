package plan

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lexer"
	"github.com/goropikari/simpledbgo/parser"
)

// BasicQueryPlanner is a basic query planner.
type BasicQueryPlanner struct {
	metadataMgr domain.MetadataManager
}

// NewBasicQueryPlanner constructs a BasicQueryPlanner.
func NewBasicQueryPlanner(metadataMgr domain.MetadataManager) *BasicQueryPlanner {
	return &BasicQueryPlanner{
		metadataMgr: metadataMgr,
	}
}

// CreatePlan creates a planner.
func (planner *BasicQueryPlanner) CreatePlan(data *domain.QueryData, txn domain.Transaction) (domain.Planner, error) {
	plans := make([]domain.Planner, 0)

	for _, tblName := range data.Tables() {
		viewDef, err := planner.metadataMgr.GetViewDef(tblName.ToViewName(), txn)
		if err != nil {
			return nil, err
		}

		if viewDef == "" {
			plan, err := NewTablePlan(txn, tblName, planner.metadataMgr)
			if err != nil {
				return nil, err
			}
			plans = append(plans, plan)
		} else {
			l := lexer.NewLexer(viewDef.String())
			tokens, err := l.ScanTokens()
			if err != nil {
				return nil, err
			}

			p := parser.NewParser(tokens)
			viewData, err := p.Query()
			if err != nil {
				return nil, err
			}

			plan, err := planner.CreatePlan(viewData, txn)
			if err != nil {
				return nil, err
			}
			plans = append(plans, plan)
		}
	}

	plan := plans[0]
	plans = plans[1:]
	for _, nextPlan := range plans {
		plan = NewProductPlan(plan, nextPlan)
	}

	plan = NewSelectPlan(plan, data.Predicate())

	plan = NewProjectPlan(plan, data.Fields())

	return plan, nil
}
