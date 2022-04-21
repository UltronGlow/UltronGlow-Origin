// Copyright 2021 The utg Authors
// This file is part of the utg library.
//
// The utg library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The utg library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the utg library. If not, see <http://www.gnu.org/licenses/>.

package extrastate

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/UltronGlow/UltronGlow-Origin/common"
)

func TestAddLockReward(t *testing.T) {
	InitExtraDB("e:\\home\\utg\\extrastate")
	es, _ := ExtraStateAt(common.Hash{})
	acc1 := Account{
		LockBalance: common.Big0,
		ReleaseQueue: make([]*big.Int, maxLockDays+maxReleasePeriod),
		ReleaseIdx: 0,
	}
	address1 := common.HexToAddress("ux5ea73097c65ab6fcf541e34504219efaffb13a2f")
	obj1 := newObject(es, address1, acc1)

	acc2 := Account{
		LockBalance: common.Big0,
		ReleaseQueue: make([]*big.Int, maxLockDays+maxReleasePeriod),
		ReleaseIdx: 0,
	}
	address2 := common.HexToAddress("ux1453ba0a7caea29cd653fa0d7ff266221e7b03ab")
	obj2 := newObject(es, address2, acc2)

	acc3 := Account{
		LockBalance: common.Big0,
		ReleaseQueue: make([]*big.Int, maxLockDays+maxReleasePeriod),
		ReleaseIdx: 0,
	}
	address3 := common.HexToAddress("ux37c9c11cc147797832a472279d915ada3ab949cd")
	obj3 := newObject(es, address3, acc3)

	//测试数据, 第一天锁仓记录
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))

	//测试数据，第二天
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj2.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj3.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))
	obj1.AddReward(new(big.Int).SetUint64(500000000000000000))

	fmt.Println(new(big.Int).Div(new(big.Int).SetUint64(500000000000000000), big.NewInt(int64(rewardReleasPeriod))))
	var amount *big.Int
	for i:=59; i<241; i++ {
		amount = obj1.PayReward(uint64(i), nil)
		if amount.Cmp(common.Big0) > 0 {
			fmt.Println("abj1", obj1.address, "release:", amount ,"number:", i, "releasenumberperday", obj1.paymentNumberPerDay)
		}
		amount = obj2.PayReward(uint64(i), nil)
		if amount.Cmp(common.Big0) > 0 {
			fmt.Println("abj2",  obj2.address, "release:", amount ,"number:", i, "releasenumberperday", obj2.paymentNumberPerDay)
		}
		amount = obj3.PayReward(uint64(i), nil)
		if amount.Cmp(common.Big0) > 0 {
			fmt.Println("abj3",  obj3.address, " release:", amount ,"number:", i, "releasenumberperday", obj3.paymentNumberPerDay)
		}
	}
}

func BenchmarkCutOriginal(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		bytes.TrimLeft(value[:], "\x00")
	}
}

func BenchmarkCutsetterFn(b *testing.B) {
	value := common.HexToHash("0x01")
	cutSetFn := func(r rune) bool { return r == 0 }
	for i := 0; i < b.N; i++ {
		bytes.TrimLeftFunc(value[:], cutSetFn)
	}
}

func BenchmarkCutCustomTrim(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		common.TrimLeftZeroes(value[:])
	}
}
