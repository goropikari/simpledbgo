package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/goropikari/simpledbgo/database"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

const (
	payloadBytesLength = 4
	tagLength          = 1
	bufferSize         = 1024
)

var dbmsPORT = getEnvWithDefault("DBMS_PORT", "5432")
var dbmsHOST = getEnvWithDefault("DBMS_HOST", "127.0.0.1")

type ResultType int

const (
	queryResult ResultType = iota + 1
	commandResult
	beginResult
	commitResult
	rollbackResult
)

type Result struct {
	typ     ResultType
	records [][]any
	fields  []string
}

func (res Result) isQuery() bool {
	return res.typ == queryResult
}

type Server struct {
	db *database.DB
}

func NewServer() *Server {
	db := setupDB()
	return &Server{db: db}
}

type Connection struct {
	db    *database.DB
	conn  net.Conn
	txn   domain.Transaction
	inTxn bool
}

func NewConnection(db *database.DB, conn net.Conn) Connection {
	return Connection{
		db:   db,
		conn: conn,
	}
}

func (cn *Connection) Txn() (domain.Transaction, error) {
	if cn.inTxn {
		return cn.txn, nil
	}

	txn, err := cn.db.NewTx()
	if err != nil {
		return nil, errors.Err(err, "NewTx")
	}

	return txn, nil
}

// Run starts DBMS server
func (s *Server) Run() {
	ln, err := net.Listen("tcp", dbmsHOST+":"+dbmsPORT)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v\n", err)
		}
		cn := NewConnection(s.db, conn)
		go cn.handleConnection()
	}
}

func (cn *Connection) handleConnection() {
	cn.startup()
	defer cn.Close()
	for {
		tag, query, err := cn.readQuery()
		if err != nil {
			if err != io.EOF {
				log.Printf("%v\n", err)
				os.Exit(1)
			}
			break
		}
		if tag == 0x58 {
			// 0x58 -> X: terminate
			return
		}
		res, err := cn.handleQuery(query)
		if err != nil {
			log.Printf("%v\n", err)
			// Ideally, error msg should be sent if errors occur
			cn.conn.Write(makeCommandCompleteMsg(baseError(err).Error()))
			cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
			continue
		}
		// if res.isQuery() {
		// 	cn.sendResult(res)
		// } else {
		// 	// Query except for SELECT
		// 	cn.conn.Write(makeCommandCompleteMsg("OK"))
		// 	cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
		// }
		switch res.typ {
		case queryResult:
			cn.sendResult(res)
		case commandResult:
			cn.sendCommand(res)
		case beginResult:
			cn.sendBegin(res)
		case commitResult:
			cn.sendCommit(res)
		case rollbackResult:
			cn.sendRollback(res)
		}
		cn.sendReadyForQueryMsg()
	}
}

func (cn *Connection) sendCommand(res Result) {
	cn.conn.Write(makeCommandCompleteMsg("OK"))
}

func (cn *Connection) sendBegin(res Result) {
	cn.conn.Write(makeCommandCompleteMsg("BEGIN"))
}

func (cn *Connection) sendCommit(res Result) {
	cn.conn.Write(makeCommandCompleteMsg("COMMIT"))
}

func (cn *Connection) sendRollback(res Result) {
	cn.conn.Write(makeCommandCompleteMsg("ROLLBACK"))
}

func (cn *Connection) sendReadyForQueryMsg() {
	if cn.inTxn {
		cn.conn.Write(makeReadyForQueryMsg(Transaction))
	} else {
		cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
	}
}

func (cn *Connection) startup() error {
	// https://www.pgcon.org/2014/schedule/attachments/330_postgres-for-the-wire.pdf
	// https://www.postgresql.org/docs/12/protocol-message-formats.html
	sizeByte, err := read(cn.conn, payloadBytesLength)
	if err != nil {
		return errors.Err(err, "read")
	}
	cn.conn.Write([]byte{0x4e})

	size := int(binary.BigEndian.Uint32(sizeByte))
	if _, err := read(cn.conn, size-payloadBytesLength); err != nil {
		return errors.Err(err, "read")
	}
	// AuthenticationOk
	// 0x52 -> Z: ReadyForQuery
	cn.conn.Write([]byte{0x52, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00})
	// 0x53 -> S: ParameterStatus
	// fake client encoding for python PostgreSQL connector
	cn.conn.Write(makeParameterStatusMsg("client_encoding", "SQL_ASCII"))
	// fake server version
	cn.conn.Write(makeParameterStatusMsg("server_version", "0.0.0"))

	// ReadyForQuery
	cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))

	return nil
}

func (cn *Connection) Close() error {
	return cn.conn.Close()
}

func (cn *Connection) sendResult(res Result) {
	// 	cols := []string{"hoge", "piyo"}
	// 	header := makeColDesc(cols)
	// 	c.Write(header)
	// 	recs := [][]any{
	// 		{"1", "taro"},
	// 		{"2", "hanako"},
	// 	}
	// 	if len(recs) != 0 {
	// 		rowByte := makeDataRows(recs)
	// 		c.Write(rowByte)
	// 	}

	// 	c.Write(selectFooter(len(recs)))
	// 	c.Write(makeReadyForQueryMsg(TransactionIdle))

	cols := res.fields
	header := makeColDesc(cols)
	cn.conn.Write(header)
	recs := res.records
	if len(recs) != 0 {
		rowByte := makeDataRows(recs)
		cn.conn.Write(rowByte)
	}

	cn.conn.Write(selectFooter(len(recs)))
	// cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
}

func selectFooter(n int) []byte {
	body := []byte("SELECT ")
	s := fmt.Sprintf("%v", n)
	body = append(body, []byte(s)...)
	body = append(body, 0x00)

	payload := make([]byte, 0)
	payload = append(payload, 0x43)
	lenBytes := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(body)+payloadBytesLength))
	payload = append(payload, lenBytes...)
	payload = append(payload, body...)

	return payload
}

func makeDataRow(rec []any) []byte {
	dataRow := make([]byte, 0)
	nc := len(rec)
	ncb := make([]byte, 2)
	binary.BigEndian.PutUint16(ncb, uint16(nc))
	dataRow = append(dataRow, ncb...)
	for _, val := range rec {
		if val == nil {
			dataRow = append(dataRow, []byte{0xff, 0xff, 0xff, 0xff}...)
		} else {
			s := fmt.Sprintf("%v", val)
			sb := []byte(s)
			slen := len(sb)
			lenByte := make([]byte, payloadBytesLength)
			binary.BigEndian.PutUint32(lenByte, uint32(slen))
			dataRow = append(dataRow, lenByte[:]...)
			dataRow = append(dataRow, sb[:]...)
		}
	}

	payload := make([]byte, 0)
	payload = append(payload, 0x44) // 0x44 -> D: DataRow
	lenByte := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lenByte, uint32(len(dataRow)+payloadBytesLength))
	payload = append(payload, lenByte...)
	payload = append(payload, dataRow...)

	return payload
}

func makeColDesc(cols []string) []byte {
	payload := make([]byte, 0)
	n := len(cols)
	numCols := make([]byte, 2)
	binary.BigEndian.PutUint16(numCols, uint16(n))
	payload = append(payload, numCols[:]...)

	for k, col := range cols {
		payload = append(payload, []byte(col)...)
		payload = append(payload, 0x00)
		payload = append(payload, []byte{0x00, 0x00, 0x40, 0x06}...) // object id
		idx := make([]byte, 2)
		binary.BigEndian.PutUint16(idx, uint16(k+1))
		payload = append(payload, idx[:]...)                         // col id
		payload = append(payload, []byte{0x00, 0x00, 0x04, 0x13}...) // data type
		payload = append(payload, []byte{0xff, 0xff}...)             // data type size
		payload = append(payload, []byte{0xff, 0xff, 0xff, 0xff}...) // type modifier
		payload = append(payload, []byte{0x00, 0x00}...)             // format code
	}

	length := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(length, uint32(len(payload)+payloadBytesLength))
	packet := make([]byte, 0)
	packet = append(packet, 0x54) // 0x54 -> T: RowDescription
	packet = append(packet, length[:]...)
	packet = append(packet, payload[:]...)

	return packet
}

func makeDataRows(recs [][]any) []byte {
	dataRows := make([]byte, 0)
	for _, rec := range recs {
		dataRows = append(dataRows, makeDataRow(rec)...)
	}

	return dataRows
}

func (cn *Connection) handleQuery(query string) (Result, error) {
	fmt.Println(query)
	query = strings.TrimRight(query, ";")
	if strings.HasPrefix(query, "select") {
		return cn.handleSelect(query)
	}
	if strings.HasPrefix(query, "begin") {
		return cn.handleBegin(query)
	}
	if strings.HasPrefix(query, "commit") {
		return cn.handleCommit(query)
	}
	if strings.HasPrefix(query, "rollback") {
		return cn.handleRollback(query)
	}
	return cn.handleCommand(query)
}

func (cn *Connection) handleBegin(query string) (Result, error) {
	cn.inTxn = false
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "Txn")
	}
	cn.txn = txn
	cn.inTxn = true

	return Result{typ: beginResult}, nil
}

func (cn *Connection) handleCommit(query string) (Result, error) {
	cn.inTxn = false
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "Txn")
	}
	if err := txn.Commit(); err != nil {
		return Result{}, errors.Err(err, "Commit")
	}
	cn.txn = nil

	return Result{typ: commitResult}, nil
}

func (cn *Connection) handleRollback(query string) (Result, error) {
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "Txn")
	}
	if err := txn.Rollback(); err != nil {
		return Result{}, errors.Err(err, "Rollback")
	}
	cn.txn = nil
	cn.inTxn = false

	return Result{typ: rollbackResult}, nil
}

func (cn *Connection) handleSelect(query string) (Result, error) {
	// txn, err := cn.db.NewTx()
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "NewTx")
	}
	p, err := cn.db.Query(txn, query)
	if err != nil {
		return Result{}, cn.rollback(txn, err)
	}

	scan, err := p.Open()
	if err != nil {
		return Result{}, errors.Err(err, "Open")
	}

	rows := &Rows{scan: scan, fields: p.Schema().Fields()}
	result, err := cn.makeResult(rows)
	if err != nil {
		return Result{}, cn.rollback(txn, err)
	}

	if !cn.inTxn {
		if err := txn.Commit(); err != nil {
			return Result{}, cn.rollback(txn, err)
		}
	}

	return result, nil
}

func (cn *Connection) handleCommand(query string) (Result, error) {
	// txn, err := cn.db.NewTx()
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "NewTx")
	}
	if _, err = cn.db.Exec(txn, query); err != nil {
		return Result{}, cn.rollback(txn, err)
	}

	if !cn.inTxn {
		if err := txn.Commit(); err != nil {
			return Result{}, cn.rollback(txn, err)
		}
	}

	return Result{typ: commandResult}, nil
}

func (cn *Connection) rollback(txn domain.Transaction, err error) error {
	cn.inTxn = false
	if err2 := txn.Rollback(); err2 != nil {
		panic(err2)
	}

	return err
}

type Rows struct {
	scan   domain.Scanner
	fields []domain.FieldName
}

func (cn *Connection) makeResult(rows *Rows) (Result, error) {
	recs := make([][]any, 0)
	for rows.scan.HasNext() {
		rec := make([]any, 0)
		for _, fld := range rows.fields {
			v, err := rows.scan.GetVal(fld)
			if err != nil {
				return Result{}, errors.Err(err, "GetVal")
			}
			rec = append(rec, v)
		}
		recs = append(recs, rec)
	}

	fields := make([]string, 0)
	for _, fld := range rows.fields {
		fields = append(fields, string(fld))
	}

	return Result{
		typ:     queryResult,
		records: recs,
		fields:  fields,
	}, nil
}

func (cn *Connection) readQuery() (byte, string, error) {
	data := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := cn.conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				os.Exit(1)
			}
			break
		}
		if n < bufferSize {
			data = append(data, buf[:n]...)
			break
		}
		data = append(data, buf[:]...)
	}
	tag := data[0]
	size := parseSize(data[1:5])
	var query string
	if size >= 5 {
		query = string(data[5:size][:])
	}

	return tag, query, nil
}

func parseSize(bs []byte) int {
	return int(binary.BigEndian.Uint32(bs))
}

// read n bytes from c.
func read(c net.Conn, n int) ([]byte, error) {
	reader := bufio.NewReader(c)
	data := make([]byte, 0)
	for i := 0; i < n; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, errors.Err(err, "ReadByte")
		}
		data = append(data, b)
	}

	return data, nil
}

func setupDB() *database.DB {
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func baseError(err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	return err
}

func getEnvWithDefault(key string, d string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return d
}
