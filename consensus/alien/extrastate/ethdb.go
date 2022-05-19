package extrastate

import (
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/core/rawdb"
	"github.com/UltronGlow/UltronGlow-Origin/ethdb"
	"os"
	"os/user"
	"sync"

	//"github.com/UltronGlow/UltronGlow-Origin/log"
)

var extradb Database
var ldb ethdb.Database
var dbpath string
var dbMutex sync.RWMutex
/*
type ldbController struct {
	extradb Database
	ldb ethdb.Database
	mutex sync.RWMutex
}

var globalLdbCtrl *ldbController
*/

func InitExtraDB(dbpath string) error {
	/*
	db, err := newEthDataBase(dbpath)
	if err != nil {
		return err
	}
	globalLdbCtrl =&ldbController{
		extradb: NewDatabase(db),
		ldb: db,
	}
	return nil

	 */

	var err error
	db, err := newEthDataBase(dbpath)
	extradb = NewDatabase(db)
	return err
}
/*
func (l *ldbController) DataBase() Database{
	return l.extradb
}

func (l *ldbController) ethDataBase() ethdb.Database {
	return l.ldb
}
*/
func newEthDataBase(path string) (ethdb.Database, error) {
	var (
		//db  ethdb.Database
		err error
	)
	if path == "" {
		dbpath = homeDir()+"/utg/extrastate"
	}else {
		dbpath = path
	}
	ldb, err = rawdb.NewLevelDBDatabase(dbpath, 1, 0, "extrastate", false)
	if nil != err {
		panic(fmt.Sprintf("open extrastate database failed, err=%s", err.Error()))
	}

	return ldb, nil
}

/*
func Commi(es *ExtraStateDB,node common.Hash, report bool, callback func(common.Hash)) error {
	globalLdbCtrl.mutex.Lock()
	defer globalLdbCtrl.mutex.Unlock()
	return globalLdbCtrl.extradb.TrieDB().Commit(node, report, callback)
}
*/
func DBPut(key, value []byte) error {
	//globalLdbCtrl.mutex.Lock()
	//defer globalLdbCtrl.mutex.Unlock()
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return ldb.Put(key, value)
}

func DBGet(key []byte) ([]byte, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return ldb.Get(key)
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