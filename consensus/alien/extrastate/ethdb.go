package extrastate

import (
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/core/rawdb"
	"github.com/UltronGlow/UltronGlow-Origin/ethdb"
	"os"
	"os/user"
	//"github.com/UltronGlow/UltronGlow-Origin/log"
)

var extradb Database
var ldb ethdb.Database

func InitExtraDB(dbpath string) error {
	var err error
	db, err := newEthDataBase(dbpath)
	extradb = NewDatabase(db)
	return err
}

func newEthDataBase(dbpath string) (ethdb.Database, error) {
	var (
		//db  ethdb.Database
		err error
	)

	if dbpath == "" {
		ldb, err = rawdb.NewLevelDBDatabase(homeDir()+"/utg/extrastate", 1, 0, "extrastate", false)
	} else {
		ldb, err = rawdb.NewLevelDBDatabase(dbpath, 1, 0, "extrastate", false)
	}

	if nil != err {
		panic(fmt.Sprintf("open extrastate database failed, err=%s", err.Error()))
	}

	return ldb, nil
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func DBPut(key, value []byte) error {
	return ldb.Put(key, value)
}

func DBGet(key []byte) ([]byte, error) {
	return ldb.Get(key)
}
