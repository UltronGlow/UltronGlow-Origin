package extrastate

import (
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/core/rawdb"
	"github.com/UltronGlow/UltronGlow-Origin/core/types"
	"github.com/UltronGlow/UltronGlow-Origin/crypto"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"github.com/UltronGlow/UltronGlow-Origin/rlp"
	"math/big"
	"testing"
	"time"
)

func _TestExtarState(t *testing.T) {
	root := common.HexToHash("uxb74f0ef64b82df3124f367fae46e2d541c9bd504a87e645f48da883633df172b")
	address1 := common.HexToAddress("0x823140710bf13990e4500136726d8b55")
	address2 := common.HexToAddress("0x823140710bf13990e4500136726d8b56")

	ethdb, err := newEthDataBase("")
	if err != nil {
		fmt.Println("newEthDataBase failed:", err)
		return
	}

	db := NewDatabase(ethdb)
	tr, err := db.OpenTrie(root)
	if err != nil {
		fmt.Println("OpenTrie failed:", err)
		return
	}

	sdb := &ExtraStateDB{
		db:                  db,
		trie:                tr,
		originalRoot:        root,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		hasher:              crypto.NewKeccakState(),
	}

	sdb.AddLockReward(address1, big.NewInt(int64(18000000000)))
	shash := sdb.IntermediateRoot(true)
	sdb.Commit(true)
	sdb.Database().TrieDB().Commit(shash, true, nil)
	fmt.Println("extradata root hash:", shash)
	fmt.Println("================================================")

	sdb.AddLockReward(address2, big.NewInt(int64(36000000000)))
	shash = sdb.IntermediateRoot(true)
	_, err = sdb.Commit(true)
	if err != nil {
		fmt.Println("sdb.Commit failed:", err)
		return
	}
	sdb.Database().TrieDB().Commit(shash, true, nil)
	fmt.Println("extradata root hash:", shash)
	fmt.Println("================================================")

	obj := sdb.getStateObject(address1)
	if obj == nil {
		fmt.Println("address:", address1, "not find in state")
		return
	}
	fmt.Println(address1, "totalLockReward:", obj.data.LockBalance, "Balance:", obj.data.Balance)
	fmt.Println("================================================")
	ethdb.Close()
}

func _TestGetObject(t *testing.T) {
	root := common.HexToHash("ux9fd46bee9c164811b6d63418021dd4902a08b8fb3b0082cda8ce1efe964488f0")
	address1 := common.HexToAddress("0x823140710bf13990e4500136726d8b55")
	address2 := common.HexToAddress("0x823140710bf13990e4500136726d8b56")

	ethdb, err := newEthDataBase("")
	if err != nil {
		fmt.Println("newEthDataBase failed:", err)
		return
	}
	defer ethdb.Close()

	db := NewDatabase(ethdb)
	tr, err := db.OpenTrie(root)
	if err != nil {
		fmt.Println("OpenTrie failed:", err)
		return
	}

	sdb := &ExtraStateDB{
		db:                  db,
		trie:                tr,
		originalRoot:        root,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		hasher:              crypto.NewKeccakState(),
	}

	obj := sdb.getStateObject(address1)
	if obj == nil {
		fmt.Println("address:", address1, "not find in state")
		return
	}
	fmt.Println("address1:", address1, "totalLockReward:", obj.data.LockBalance)
	fmt.Println("address1:", address1, "releaseQueue:", obj.data.ReleaseQueue)

	obj = sdb.getStateObject(address2)
	if obj == nil {
		fmt.Println("address2:", address2, "not find in state")
		return
	}
	fmt.Println("address2:", address2, "totalLockReward:", obj.data.LockBalance)
	fmt.Println("address2:", address2, "releaseQueue:", obj.data.ReleaseQueue)
}

func BenchmarkAddBalance(b *testing.B) {
	root := common.HexToHash("ux9fd46bee9c164811b6d63418021dd4902a08b8fb3b0082cda8ce1efe964488f0")
	address1 := common.HexToAddress("0x823140710bf13990e4500136726d8b55")
	ethdb, err := newEthDataBase("")
	if err != nil {
		fmt.Println("newEthDataBase failed:", err)
		return
	}
	defer ethdb.Close()

	db := NewDatabase(ethdb)
	tr, err := db.OpenTrie(root)
	if err != nil {
		fmt.Println("OpenTrie failed:", err)
		return
	}
	defer ethdb.Close()

	sdb := &ExtraStateDB{
		db:                  db,
		trie:                tr,
		originalRoot:        root,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		hasher:              crypto.NewKeccakState(),
	}

	obj := sdb.getStateObject(address1)
	if obj == nil {
		fmt.Println("address:", address1, "not find in state")
		return
	}
	balanc := big.NewInt(int64(18000))
	for i := 0; i < 1000000; i++ {
		obj.AddBalance(balanc)
	}
}

func TestLockAccountStore(t *testing.T) {
	ldb, err := rawdb.NewLevelDBDatabase("e:\\home\\utg\\extrastate", 1, 0, "extrastate", false)
	if err != nil {
		fmt.Sprintf("open leveldb failed, err=%v", err)
		return
	}
	defer ldb.Close()
	la := newLockAccounts()
	lockcount := 10000
	begin := time.Now()
	var lockAddrs []string//[]common.Address

	/*
	for i:=0; i<lockcount; i++ {
		addr := common.BigToAddress(new(big.Int).SetInt64((int64(i))))
		lockAddrs = append(lockAddrs, addr)
	}
	*/

	lockAddrs = append(lockAddrs, "ux3EFbaA018F81f9d0ef8b3a3d0D0b767d3482B66b")
	lockAddrs = append(lockAddrs, "ux1EB90474374F2EA4deB3961cdD0A391821C0321b")
	data, err := rlp.EncodeToBytes(lockAddrs)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %v", err))
		return
	}
	key := crypto.Keccak256Hash(data)
	err = ldb.Put(key.Bytes(), data)
	if err != nil {
		log.Error("extrastate dbput", "error", err)
		return
	}
	fmt.Println(fmt.Sprintf("acountsLen=%d, root=%v", len(lockAddrs), key))
	end := time.Now()
	fmt.Println("account slic put spend", end.Sub(begin).Milliseconds(), "ms")
	//=======================================================================================================

	begin = time.Now()
	for i:=0; i<lockcount; i++ {
		addr := common.BigToAddress(new(big.Int).SetInt64((int64(i))))
		la.Append(addr)
	}
	la.commit(ldb)
	end2 := time.Now()
	fmt.Println("lockaccount object commit spend:", end2.Sub(begin).Milliseconds(), "ms, len:", la.Len(), "root:", la.Root())
/*
	begin = time.Now()
	la2  := newLockAccounts()
	for i:=lockcount-1; i>=0;  {
		addr := common.BigToAddress(new(big.Int).SetInt64((int64(i))))
		la2.append(addr)
		i -= 1
	}
	la2.commit(ldb)
	end2 = time.Now()
	fmt.Println("lockaccount object commit spend:", end2.Sub(begin).Milliseconds(), "ms, len:", la2.Len(), "root:", la.Root())
 */
}

func TestLockAccountLoad(t *testing.T) {
	ldb, err := rawdb.NewLevelDBDatabase("e:\\home\\utg\\extrastate", 1, 0, "extrastate", false)
	if err != nil {
		fmt.Sprintf("open leveldb failed, err=%v", err)
		return
	}
	defer ldb.Close()

	root := common.HexToHash("ux8b355c5c061716565f97e6dd342d271686369db158c72eb8477969232556310b")
	begin := time.Now()
	la := newLockAccounts()
	la.load(ldb, root)
	fmt.Println(fmt.Sprintf("loadlockaccounts acountsLen=%d, root=%v", len(la.Accounts()), root), "spend", time.Now().Sub(begin).Milliseconds(), "ms")
}