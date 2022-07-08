package metadata

import (
	"sync"

	"github.com/goropikari/simpledbgo/domain"
)

// StatManager is stat manager.
type StatManager struct {
	mu         *sync.Mutex
	tblMgr     *TableManager
	tableStats map[domain.TableName]domain.StatInfo
	numCalls   int
}

// NewStatManager constructs a StatManager.
// StatManager は一つの instance しか作らない(p.197)ので refreshStatistics の
// 呼び出しを lock する必要はない。
// GetStatInfo だけ lock すれば十分.
func NewStatManager(tblMgr *TableManager, txn domain.Transaction) (*StatManager, error) {
	statMgr := &StatManager{
		mu:     &sync.Mutex{},
		tblMgr: tblMgr,
	}

	if err := statMgr.refreshStatistics(txn); err != nil {
		return nil, err
	}

	return statMgr, nil
}

// GetStatInfo returns statistics of given table.
func (statMgr *StatManager) GetStatInfo(tblName domain.TableName, layout *domain.Layout, txn domain.Transaction) (domain.StatInfo, error) {
	statMgr.mu.Lock()
	defer statMgr.mu.Unlock()

	statMgr.numCalls++
	if statMgr.numCalls > updateTimes {
		// refreshStatistics で numCalls は 0 リセットされる
		if err := statMgr.refreshStatistics(txn); err != nil {
			return domain.StatInfo{}, err
		}
	}

	si, found := statMgr.tableStats[tblName]
	if !found {
		var err error
		si, err = statMgr.calcTableStats(tblName, layout, txn)
		if err != nil {
			return domain.StatInfo{}, err
		}
		statMgr.tableStats[tblName] = si
	}

	return si, nil
}

func (statMgr *StatManager) refreshStatistics(txn domain.Transaction) error {
	statMgr.tableStats = make(map[domain.TableName]domain.StatInfo)
	statMgr.numCalls = 0

	catLayout, err := statMgr.tblMgr.GetTableLayout(tableCatalog, txn)
	if err != nil {
		return err
	}

	tcat, err := domain.NewTableScan(txn, tableCatalog, catLayout)
	if err != nil {
		return err
	}
	defer tcat.Close()

	for tcat.HasNext() {
		tblNameStr, err := tcat.GetString(fldTableName)
		if err != nil {
			return err
		}

		tblName, err := domain.NewTableName(tblNameStr)
		if err != nil {
			return err
		}

		layout, err := statMgr.tblMgr.GetTableLayout(tblName, txn)
		if err != nil {
			return err
		}

		si, err := statMgr.calcTableStats(tblName, layout, txn)
		if err != nil {
			return err
		}
		statMgr.tableStats[tblName] = si
	}
	if err := tcat.Err(); err != nil {
		return err
	}

	return nil
}

func (statMgr *StatManager) calcTableStats(tblName domain.TableName, layout *domain.Layout, txn domain.Transaction) (domain.StatInfo, error) {
	numRecs := 0
	numBlocks := 0

	tbl, err := domain.NewTableScan(txn, tblName, layout)
	if err != nil {
		return domain.StatInfo{}, err
	}

	for tbl.HasNext() {
		numRecs++
		numBlocks = int(tbl.RecordID().BlockNumber() + 1)
	}
	if err := tbl.Err(); err != nil {
		return domain.StatInfo{}, tbl.Err()
	}

	return domain.NewStatInfo(numBlocks, numRecs), nil
}