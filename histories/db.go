package histories

import (
	"github.com/if1live/poloniex-history-viewer/balances"
	"github.com/if1live/poloniex-history-viewer/exchanges"
	"github.com/if1live/poloniex-history-viewer/lendings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(filepath string) (Database, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return Database{}, err
	}

	db.AutoMigrate(&exchanges.Exchange{})
	db.AutoMigrate(&lendings.Lending{})
	db.AutoMigrate(&balances.Transaction{})

	return Database{
		db: db,
	}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetORM() *gorm.DB {
	return d.db
}

func (d *Database) MakeExchangeSync(api *poloniex.Poloniex) *exchanges.Sync {
	return exchanges.NewSync(d.db, api)
}

func (d *Database) MakeLendingSync(api *poloniex.Poloniex) *lendings.Sync {
	return lendings.NewSync(d.db, api)
}
func (d *Database) MakeBalanceSync(api *poloniex.Poloniex) *balances.Sync {
	return balances.NewSync(d.db, api)
}
