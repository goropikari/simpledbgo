package embedded

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"

	"github.com/goropikari/simpledbgo/database"
	"github.com/goropikari/simpledbgo/domain"
)

// Driver is sql driver.
type Driver struct{}

// Open opens database.
func (d *Driver) Open(name string) (driver.Conn, error) {
	// fmt.Println("*Driver.Open: " + name)
	db, err := database.InitializeDB()
	if err != nil {
		return nil, err
	}

	return &Conn{db: db, inTxn: false}, nil
}

// Conn is connection of simple database.
type Conn struct {
	db    *database.DB
	inTxn bool
	txn   domain.Transaction
}

// Stmt is statement of query/command.
type Stmt struct {
	cn   *Conn
	plan domain.Planner
	cmd  string
}

// Prepare satisfies driver.Conn interface.
func (cn *Conn) Prepare(query string) (driver.Stmt, error) {
	// fmt.Println("*Conn.Prepare: " + query)
	var txn domain.Transaction
	var err error
	if cn.inTxn {
		txn = cn.txn
	} else {
		txn, err = cn.db.NewTx()
		if err != nil {
			return nil, err
		}
	}

	cn.txn = txn

	var p domain.Planner
	var cmd string
	if strings.HasPrefix(query, "select") {
		p, err = cn.db.Query(txn, query)
		if err != nil {
			return nil, err
		}
	} else {
		cmd = query
	}

	return &Stmt{cn: cn, plan: p, cmd: cmd}, nil
}

// Close satisfies driver.Conn interface.
func (cn *Conn) Close() error {
	// fmt.Println("Conn.Close")

	return nil
}

// Begin satisfies driver.Conn interface.
func (cn *Conn) Begin() (driver.Tx, error) {
	// fmt.Println("Conn.Begin")
	cn.inTxn = true
	txn, err := cn.db.NewTx()
	if err != nil {
		return nil, err
	}

	cn.txn = txn

	return cn, nil
}

// Close satisfies driver.Stmt interface.
func (stmt *Stmt) Close() error {
	// fmt.Println("Stmt.Close")
	if !stmt.cn.inTxn {
		return stmt.cn.txn.Commit()
	}

	return nil
}

// NumInput satisfies driver.Stmt interface.
func (stmt *Stmt) NumInput() int {
	// fmt.Println("Stmt.NumInput")

	return -1
}

// Exec satisfies driver.Stmt interface.
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	// fmt.Println("Stmt.Exec")
	_, err := stmt.cn.db.Exec(stmt.cn.txn, stmt.cmd)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Query satisfies driver.Stmt interface.
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("not implemented")
}

// QueryContext satisfies driver.StmtQueryContext interface.
func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	// fmt.Println("Stmt.QueryContext")
	scan, err := stmt.plan.Open()
	if err != nil {
		return nil, err
	}

	fields := stmt.plan.Schema().Fields()

	if !stmt.cn.inTxn {
		return nil, stmt.cn.txn.Commit()
	}

	return &Rows{
		scan:   scan,
		fields: fields,
	}, nil
}

// Commit satisfies driver.Tx interface.
func (cn *Conn) Commit() error {
	// fmt.Println("Tx.Commit")
	cn.inTxn = false

	return cn.txn.Commit()
}

// Rollback satisfies driver.Tx interface.
func (cn *Conn) Rollback() error {
	// fmt.Println("Tx.Rollback")

	return cn.txn.Rollback()
}

// Rows implements driver.Rows.
type Rows struct {
	scan   domain.Scanner
	fields []domain.FieldName
}

// Columns satisfies drivers.Rows.
func (r *Rows) Columns() []string {
	// fmt.Println("Rows.Columns")
	cols := make([]string, len(r.fields))
	for i, f := range r.fields {
		cols[i] = string(f)
	}

	return cols
}

// Close satisfies drivers.Rows.
func (r *Rows) Close() error {
	// fmt.Println("Rows.Close")
	r.scan.Close()

	return nil
}

// Next satisfies drivers.Rows.
// dest は Rows.Columns が返すスライスと同じ長さ.
func (r *Rows) Next(dest []driver.Value) error {
	// fmt.Println("Rows.Next")
	if r.scan.HasNext() {
		for i, f := range r.fields {
			v, err := r.scan.GetVal(f)
			if err != nil {
				return err
			}
			dest[i] = v.AsVal()
		}

		return nil
	}

	return io.EOF
}

func init() {
	sql.Register("simpledb", &Driver{})
}
