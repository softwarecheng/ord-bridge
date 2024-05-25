package g

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/common/db/maintain"
	"github.com/softwarecheng/ord-bridge/common/log"
)

func releaseDB(db *badger.DB) {
	if db == nil || db.IsClosed() {
		return
	}
	err := db.Sync()
	if err != nil {
		log.Log.Errorf("releaseDB-> syncing db error: %v", err)
	}
	err = db.Close()
	if err != nil {
		log.Log.Errorf("releaseDB-> closing db error: %v", err)
	}
}

func openDB(dbname, dbDir string, opts badger.Options) (db *badger.DB, err error) {
	dbDir += dbname
	opts = opts.WithDir(dbDir).WithValueDir(dbDir).WithLoggingLevel(badger.WARNING)
	db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}
	log.Log.WithField("db", dbDir).Info("openDB->")
	return db, nil
}

func InitDB() (err error) {
	log.Log.Info("InitDB-> start...")
	err = InitOrdDB(Cfg.DB.Path, 3000<<20)
	if err != nil {
		return err
	}
	return nil
}

func InitOrdDB(dbDir string, size int64) (err error) {
	opts := badger.DefaultOptions("ord").
		WithBlockCacheSize(size)
	DB.Ord, err = openDB("ord", dbDir, opts)
	if err != nil {
		return err
	}
	log.Log.Info("InitOrdDB-> gc db...")
	maintain.Gc(DB.Ord, 0.5)
	return nil
}

func ReleaseAllDB() {
	log.Log.Info("ReleaseOrdDB")
	ReleaseOrdDB()
}

func ReleaseOrdDB() {
	releaseDB(DB.Ord)
}
