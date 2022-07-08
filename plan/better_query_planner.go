package plan

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/lexer"
	"github.com/goropikari/simpledbgo/parser"
)

// BetterQueryPlanner is a basic query planner.
type BetterQueryPlanner struct {
	metadataMgr domain.MetadataManager
}

// NewBetterQueryPlanner constructs a BetterQueryPlanner.
func NewBetterQueryPlanner(metadataMgr domain.MetadataManager) *BetterQueryPlanner {
	return &BetterQueryPlanner{
		metadataMgr: metadataMgr,
	}
}

// CreatePlan creates a planner.
func (planner *BetterQueryPlanner) CreatePlan(data *domain.QueryData, txn domain.Transaction) (domain.Planner, error) {
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
		choice1 := NewProductPlan(plan, nextPlan)
		choice2 := NewProductPlan(nextPlan, plan)
		if choice1.EstNumBlocks() < choice2.EstNumBlocks() {
			plan = choice1
		} else {
			plan = choice2
		}
	}

	plan = NewSelectPlan(plan, data.Predicate())

	plan = NewProjectPlan(plan, data.Fields())

	return plan, nil
}
