package plan

import (
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

// TablePlan is planner for table.
type TablePlan struct {
	tblName domain.TableName
	txn     domain.Transaction
	layout  *domain.Layout
	si      domain.StatInfo
}

// NewTablePlan constructs a TablePlan.
func NewTablePlan(txn domain.Transaction, tblName domain.TableName, md domain.MetadataManager) (*TablePlan, error) {
	layout, err := md.GetTableLayout(tblName, txn)
	if err != nil {
		return nil, errors.Err(err, "GetTableLayout")
	}

	si, err := md.GetStatInfo(tblName, layout, txn)
	if err != nil {
		return nil, errors.Err(err, "GetStatInfo")
	}

	return &TablePlan{
		tblName: tblName,
		txn:     txn,
		layout:  layout,
		si:      si,
	}, nil
}

// Open opens scanner.
func (plan *TablePlan) Open() (domain.Scanner, error) {
	return domain.NewTableScan(plan.txn, plan.tblName, plan.layout)
}

// EstNumBlocks estimates the number of block access.
func (plan *TablePlan) EstNumBlocks() int {
	return plan.si.EstNumBlocks()
}

// EstNumRecord estimates the number of record access.
func (plan *TablePlan) EstNumRecord() int {
	return plan.si.EstNumRecord()
}

// EstDistinctVals estimates the number of distinct value at given fldName.
func (plan *TablePlan) EstDistinctVals(fldName domain.FieldName) int {
	return plan.si.EstDistinctVals(fldName)
}

// Schema returns schema of table schema.
func (plan *TablePlan) Schema() *domain.Schema {
	return plan.layout.Schema()
}

// ProductPlan is planner for product of table.
type ProductPlan struct {
	lhsPlan domain.Planner
	rhsPlan domain.Planner
	schema  *domain.Schema
}

// NewProductPlan constructs a ProductPlan.
func NewProductPlan(lhs, rhs domain.Planner) *ProductPlan {
	sch := domain.NewSchema()
	sch.AddAllFields(lhs.Schema())
	sch.AddAllFields(rhs.Schema())

	return &ProductPlan{
		lhsPlan: lhs,
		rhsPlan: rhs,
		schema:  sch,
	}
}

// Open opens scanner.
func (plan *ProductPlan) Open() (domain.Scanner, error) {
	lhs, err := plan.lhsPlan.Open()
	if err != nil {
		return nil, errors.Err(err, "Open")
	}
	rhs, err := plan.rhsPlan.Open()
	if err != nil {
		return nil, errors.Err(err, "Open")
	}

	return domain.NewProductScan(lhs, rhs)
}

// EstNumBlocks estimates the number of block access.
func (plan *ProductPlan) EstNumBlocks() int {
	return plan.lhsPlan.EstNumBlocks() + plan.lhsPlan.EstNumRecord()*plan.rhsPlan.EstNumBlocks()
}

// EstNumRecord estimates the number of record access.
func (plan *ProductPlan) EstNumRecord() int {
	return plan.lhsPlan.EstNumRecord() * plan.rhsPlan.EstNumRecord()
}

// EstDistinctVals estimates the number of distinct value at given fldName.
func (plan *ProductPlan) EstDistinctVals(fldName domain.FieldName) int {
	if plan.lhsPlan.Schema().HasField(fldName) {
		return plan.lhsPlan.EstDistinctVals(fldName)
	}

	return plan.rhsPlan.EstDistinctVals(fldName)
}

// Schema returns schema of table schema.
func (plan *ProductPlan) Schema() *domain.Schema {
	return plan.schema
}

// SelectPlan is planner for select.
type SelectPlan struct {
	plan domain.Planner
	pred *domain.Predicate
}

// NewSelectPlan constructs a SelectPlan.
func NewSelectPlan(plan domain.Planner, pred *domain.Predicate) *SelectPlan {
	return &SelectPlan{
		plan: plan,
		pred: pred,
	}
}

// Open opens scanner.
func (p *SelectPlan) Open() (domain.Scanner, error) {
	s, err := p.plan.Open()
	if err != nil {
		return nil, errors.Err(err, "Open")
	}

	return domain.NewSelectScan(s, p.pred), nil
}

// EstNumBlocks estimates the number of block access.
func (p *SelectPlan) EstNumBlocks() int {
	return p.plan.EstNumBlocks()
}

// EstNumRecord estimates the number of record access.
func (p *SelectPlan) EstNumRecord() int {
	return p.plan.EstNumRecord() / p.pred.ReductionFactor(p.plan)
}

// EstDistinctVals estimates the number of distinct value at given fldName.
func (p *SelectPlan) EstDistinctVals(fldName domain.FieldName) int {
	if (p.pred.EquatesWithConstant(fldName) != domain.Constant{}) {
		return 1
	}

	fldName2 := p.pred.EquatesWithField(fldName)

	if fldName2 == "" {
		return p.plan.EstDistinctVals(fldName)
	}

	a, b := p.plan.EstDistinctVals(fldName), p.plan.EstDistinctVals(fldName2)

	if a < b {
		return a
	}

	return b
}

// Schema returns schema of table schema.
func (p *SelectPlan) Schema() *domain.Schema {
	return p.plan.Schema()
}

// ProjectPlan is planner of projection.
type ProjectPlan struct {
	plan   domain.Planner
	schema *domain.Schema
}

// NewProjectPlan constructs a ProjectPlan.
func NewProjectPlan(plan domain.Planner, fields []domain.FieldName) *ProjectPlan {
	sch := domain.NewSchema()
	for _, f := range fields {
		if f == "*" {
			for _, f2 := range plan.Schema().Fields() {
				sch.Add(f2, plan.Schema())
			}
		} else {
			sch.Add(f, plan.Schema())
		}
	}

	return &ProjectPlan{
		plan:   plan,
		schema: sch,
	}
}

// Open opens scanner.
func (p *ProjectPlan) Open() (domain.Scanner, error) {
	s, err := p.plan.Open()
	if err != nil {
		return nil, errors.Err(err, "Open")
	}

	return domain.NewProjectScan(s, p.schema.Fields()), nil
}

// EstNumBlocks estimates the number of block access.
func (p *ProjectPlan) EstNumBlocks() int {
	return p.plan.EstNumBlocks()
}

// EstNumRecord estimates the number of record access.
func (p *ProjectPlan) EstNumRecord() int {
	return p.plan.EstNumRecord()
}

// EstDistinctVals estimates the number of distinct value at given fldName.
func (p *ProjectPlan) EstDistinctVals(fld domain.FieldName) int {
	return p.plan.EstDistinctVals(fld)
}

// Schema returns schema of table schema.
func (p *ProjectPlan) Schema() *domain.Schema {
	return p.schema
}
