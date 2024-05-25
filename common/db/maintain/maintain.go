package maintain

import (
	"fmt"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/softwarecheng/ord-bridge/common/log"
	"github.com/softwarecheng/ord-bridge/common/util"
)

func GcFromPath(dbDir string, discardRatio float64) error {
	if !filepath.IsAbs(dbDir) {
		dbDir = filepath.Clean(dbDir) + string(filepath.Separator)
	}

	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("GcFromPath-> db directory isn't exist: %v", dbDir)
	} else if err != nil {
		return err
	}

	opts := badger.DefaultOptions(dbDir).
		WithLoggingLevel(badger.WARNING).
		WithSyncWrites(true)
	db, err := badger.Open(opts)
	if err != nil {
		return fmt.Errorf("GcFromPath-> open db error: %v", err)
	}
	defer db.Close()

	return Gc(db, discardRatio)
}

func Gc(db *badger.DB, discardRatio float64) error {
	if db.IsClosed() {
		return fmt.Errorf("db is Closed")
	}

	lastDbSize := int64(0)
	lastDbSize, err := util.GetDirSize(db.Opts().Dir)
	if err != nil {
		return err
	}
	lsmSize, logSize := db.Size()
	log.Log.Infof("Gc-> start, db dir: %v, DB size: %d MB, discardRatio: %v, lsmSize: %d, logSize: %d",
		db.Opts().Dir, lastDbSize/(1024*1024), discardRatio, lsmSize, logSize)

	gcCount := 0
	for {
		err := db.RunValueLogGC(discardRatio)
		if err == badger.ErrNoRewrite {
			break
		} else if err != nil {
			return err
		}
		dirSize, err := util.GetDirSize(db.Opts().Dir)
		if err != nil {
			log.Log.Infof("Gc-> filepath.Walk: error: %v, please ignore, no need to break", err)
		}
		if lastDbSize == dirSize {
			break
		}
		lastDbSize = dirSize
		gcCount += 1
		lsmSize, logSize := db.Size()
		log.Log.Infof("Gc-> continue, count: %v, db dir: %v, DB size: %d MB, discardRatio: %v, lsmSize: %d, logSize: %d",
			gcCount, db.Opts().Dir, lastDbSize/(1024*1024), discardRatio, lsmSize, logSize)

	}
	return nil
}

func LevelsInfo(db *badger.DB) {
	log.Log.Infof("ShowLevelsInfo-> info: %v", db.LevelsToString())
}

func Backup(fname string, db *badger.DB) error {
	backupFile, err := os.Create(fname)
	if err != nil {
		log.Log.Errorf("BackupDB-> create file %s failed. %v", fname, err)
		return err
	}
	defer backupFile.Close()
	since := uint64(0)
	latestVersion, err := db.Backup(backupFile, since)
	if err != nil {
		log.Log.Errorf("BackupDB-> Backup failed, error: %v", err)
		return err
	}
	log.Log.Infof("BackupDB-> Backup completed, new version: %d", latestVersion)
	return nil
}

func Restore(backupFile string, targetDir string) error {
	backup, err := os.Open(backupFile)
	if err != nil {
		log.Log.Errorf("RestoreDB-> Open file %s, error: %v", backupFile, err)
		return err
	}
	defer backup.Close()

	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		log.Log.Errorf("RestoreDB-> MkdirAll %s failed, error: %v", targetDir, err)
		return err
	}

	opts := badger.DefaultOptions(targetDir).WithInMemory(false)
	db, err := badger.Open(opts)
	if err != nil {
		log.Log.Errorf("RestoreDB-> Open DB failed. error: %v", err)
		return err
	}
	defer db.Close()

	err = db.Load(backup, 0)
	if err != nil {
		log.Log.Errorf("RestoreDB-> Load DB failed, error: %v", err)
		return err
	}

	log.Log.Info("RestoreDB-> DB restored")
	return nil
}
