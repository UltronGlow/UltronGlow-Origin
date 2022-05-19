package extrastate

import (
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/crypto"
	"github.com/UltronGlow/UltronGlow-Origin/ethdb"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"github.com/UltronGlow/UltronGlow-Origin/rlp"
	"sort"
)

type LockAccounts struct {
	root common.Hash
	dirtied bool
	accountSlic sort.StringSlice
	accounts map[common.Address]struct{}
}

func newLockAccounts () *LockAccounts {
	return &LockAccounts{
		root : common.Hash{},
		dirtied: false,
		accountSlic: make([]string, 0),
		accounts: make(map[common.Address]struct{}),
	}
}

func (la *LockAccounts) Append(addr common.Address) {
	if addr == common.BigToAddress(common.Big0) {

		return
	}
	if _, ok := la.accounts[addr]; !ok {
		la.accounts[addr] = struct{}{}
		la.accountSlic = append(la.accountSlic, addr.String())
		la.dirtied = true
		log.Info("extrastate lockaccounts new", "addr", addr, "len", len(la.accountSlic))
	}
}

func (la *LockAccounts) Delete(addr common.Address) {
	if _, ok := la.accounts[addr]; ok {
		delete(la.accounts, addr)
		addrStr := addr.String()
		for i, _ := range la.accountSlic {
			if la.accountSlic[i] == addrStr {
				la.accountSlic = append(la.accountSlic[:i], la.accountSlic[i+1:]...)
				la.dirtied = true
				break
			}
		}
	}
}

func (la *LockAccounts) Commit() (common.Hash, error) {
	return la.commit(ldb)
}
func (la *LockAccounts) commit(db ethdb.Database) (common.Hash, error) {
	if la.dirtied {
		la.accountSlic.Sort()
	}
	data, err := rlp.EncodeToBytes(la.accountSlic)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %v", err))
		return common.Hash{}, err
	}
	key := crypto.Keccak256Hash(data)
	err = db.Put(key.Bytes(), data)
	if err != nil {
		panic(fmt.Errorf("extrastate commit lockaccounts failed, err = %v", err))
		return common.Hash{}, err
	}
	la.root = key
	log.Info("extrastate commit lockaccounts", "accountsLen", len(la.accountSlic), "root", key)
	return la.root, nil
}

func (la *LockAccounts) Load(root common.Hash) error {
	err := la.load(ldb, root)
	if err != nil {
		return err
	}
	return nil
}
func (la *LockAccounts) load(db ethdb.Database, root common.Hash) error {
	if root == (common.Hash{}) {
		return nil
	}
	data, err := db.Get(root.Bytes())
	if err != nil {
		//panic(fmt.Errorf("extrastate loadlockaccounts, err = %v, root = %v",  err,  root))
		log.Info("extrastate loadlockaccounts", "error", err, "root", root)
		return err
	}
	err = rlp.DecodeBytes(data, &la.accountSlic)
	if err != nil {
		panic(fmt.Errorf("extrastate decode account, err = %v", err))
	}
	la.accounts = make(map[common.Address]struct{}, len(la.accountSlic))
	for _, addr := range la.accountSlic {
		la.accounts[common.HexToAddress(addr)] = struct{}{}
	}
	la.root = root
	log.Info("extrastate load lockaccounts", "accountsLen", len(la.accountSlic), "root", root)
	return nil
}

func (la *LockAccounts) Len() int {
		return len(la.accountSlic)
}

func (la *LockAccounts) Root() common.Hash {
	if la.dirtied {
		la.accountSlic.Sort()
	}
	data, err := rlp.EncodeToBytes(la.accountSlic)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %v", err))
	}
	key := crypto.Keccak256Hash(data)
	la.root = key
	return la.root
}

func (la *LockAccounts) Accounts() sort.StringSlice {
	return la.accountSlic
}