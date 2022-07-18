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
	tagLength           = 1
	payloadBytesLength  = 4
	payloadLengthOffset = tagLength
	bufferSize          = 1024
)

type ResultType int

const (
	queryResult ResultType = iota + 1
	commandResult
	beginResult
	commitResult
	rollbackResult
)

type Config struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{
		Host: getEnvWithDefault("SIMPLEDB_HOST", "127.0.0.1"),
		Port: getEnvWithDefault("SIMPLEDB_PORT", "5432"),
	}
}

type Server struct {
	cfg Config
	db  *database.DB
}

func NewServer(cfg Config) *Server {
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	return &Server{cfg: cfg, db: db}
}

// Run starts DBMS server
func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.cfg.Host+":"+s.cfg.Port)
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
	// fake postgres server version
	cn.conn.Write(makeParameterStatusMsg("server_version", "14.0.0"))

	// ReadyForQuery
	cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))

	return nil
}

func (cn *Connection) Close() error {
	return cn.conn.Close()
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

func (cn *Connection) handleConnection() {
	cn.startup()
	defer cn.Close()
	for {
		tag, query, err := cn.readQuery()
		if err != nil {
			cn.sendError(err)
			continue
		}
		if tag == 0x58 {
			// 0x58 -> X: terminate
			return
		}
		res, err := cn.handleQuery(query)
		if err != nil {
			cn.sendError(err)
			continue
		}
		switch res.typ {
		case queryResult:
			cn.sendResult(res)
		case commandResult:
			cn.sendCommand()
		case beginResult:
			cn.sendBegin()
		case commitResult:
			cn.sendCommit()
		case rollbackResult:
			cn.sendRollback()
		}
		cn.sendReadyForQueryMsg()
	}
}

func (cn *Connection) sendError(err error) {
	log.Printf("%v\n", err)
	// Ideally, error msg should be sent if errors occur
	cn.conn.Write(makeErrorMsg(baseError(err)))
	cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
}

func (cn *Connection) sendCommand() {
	cn.conn.Write(makeCommandCompleteMsg("OK"))
}

func (cn *Connection) sendBegin() {
	cn.conn.Write(makeCommandCompleteMsg("BEGIN"))
}

func (cn *Connection) sendCommit() {
	cn.conn.Write(makeCommandCompleteMsg("COMMIT"))
}

func (cn *Connection) sendRollback() {
	cn.conn.Write(makeCommandCompleteMsg("ROLLBACK"))
}

func (cn *Connection) sendReadyForQueryMsg() {
	if cn.inTxn {
		cn.conn.Write(makeReadyForQueryMsg(Transaction))
	} else {
		cn.conn.Write(makeReadyForQueryMsg(TransactionIdle))
	}
}

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
func (cn *Connection) sendResult(res Result) {
	cols := res.fields
	header := makeColDesc(cols)
	cn.conn.Write(header)
	recs := res.records
	if len(recs) != 0 {
		rowByte := makeDataRows(recs)
		cn.conn.Write(rowByte)
	}

	cn.conn.Write(makeCommandCompleteMsg(fmt.Sprintf("SELECT %v", len(recs))))
}

func (cn *Connection) handleQuery(query string) (Result, error) {
	log.Println(query)
	query = strings.TrimSpace(query)
	query = strings.TrimRight(query, ";")
	prefix := strings.ToLower(strings.Fields(query)[0])
	switch prefix {
	case "select":
		return cn.handleSelect(query)
	case "begin":
		return cn.handleBegin(query)
	case "commit":
		return cn.handleCommit(query)
	case "rollback":
		return cn.handleRollback(query)
	default:
		return cn.handleCommand(query)
	}
}

func (cn *Connection) handleSelect(query string) (Result, error) {
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
	txn, err := cn.Txn()
	if err != nil {
		return Result{}, errors.Err(err, "Txn")
	}
	if err := txn.Commit(); err != nil {
		return Result{}, errors.Err(err, "Commit")
	}
	cn.inTxn = false
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

func (cn *Connection) handleCommand(query string) (Result, error) {
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

func (cn *Connection) readQuery() (byte, string, error) {
	data := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		n, err := cn.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, "", errors.Err(err, "Read")
		}
		if n < bufferSize {
			data = append(data, buf[:n]...)
			break
		}
		data = append(data, buf[:]...)
	}
	tag := data[0]
	size := parseSize(data[payloadLengthOffset : payloadLengthOffset+payloadBytesLength])
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
