//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package database

import (
	"github.com/google/wire"
	"github.com/goropikari/simpledbgo/buffer"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/file"
	"github.com/goropikari/simpledbgo/index/dummy"
	"github.com/goropikari/simpledbgo/index/hash"
	"github.com/goropikari/simpledbgo/log"
	"github.com/goropikari/simpledbgo/metadata"
	"github.com/goropikari/simpledbgo/plan"
	"github.com/goropikari/simpledbgo/tx"
)

var SetHashIndex = wire.NewSet(
	hash.NewIndexFactory,
	wire.Bind(new(domain.IndexFactory), new(*hash.IndexFactory)),
	hash.NewSearchCostCalculator,
	wire.Bind(new(domain.SearchCostCalculator), new(*hash.SearchCostCalculator)),
)

var SetDummyIndex = wire.NewSet(
	dummy.NewIndexFactory,
	wire.Bind(new(domain.IndexFactory), new(*dummy.IndexFactory)),
	dummy.NewSearchCostCalculator,
	wire.Bind(new(domain.SearchCostCalculator), new(*dummy.SearchCostCalculator)),
)

var Set = wire.NewSet(
	file.NewManagerConfig,
	file.NewManager,
	wire.Bind(new(domain.FileManager), new(*file.Manager)),
	log.NewManagerConfig,
	log.NewManager,
	wire.Bind(new(domain.LogManager), new(*log.Manager)),
	buffer.NewConfig,
	buffer.NewManager,
	wire.Bind(new(domain.BufferPoolManager), new(*buffer.Manager)),
	tx.NewLockTableConfig,
	tx.NewLockTable,
	tx.NewNumberGenerator,
	wire.Bind(new(domain.TxNumberGenerator), new(*tx.NumberGenerator)),
	SetDummyIndex,
	// SetHashIndex,
	domain.NewIndexDriver,
	metadata.NewManager,
	wire.Bind(new(domain.MetadataManager), new(*metadata.Manager)),
	plan.NewBasicQueryPlanner,
	wire.Bind(new(domain.QueryPlanner), new(*plan.BasicQueryPlanner)),
	plan.NewBasicUpdatePlanner,
	wire.Bind(new(domain.UpdateExecutor), new(*plan.BasicUpdatePlanner)),
	plan.NewExecutor,
	NewDB,
)

func InitializeDB() (*DB, error) {
	wire.Build(Set)
	return &DB{}, nil
}
