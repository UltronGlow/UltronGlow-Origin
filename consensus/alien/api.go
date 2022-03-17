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

// Package alien implements the delegated-proof-of-stake consensus engine.

package alien

import (
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/consensus"
	"github.com/UltronGlow/UltronGlow-Origin/core/types"
	"github.com/UltronGlow/UltronGlow-Origin/ethdb"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"github.com/UltronGlow/UltronGlow-Origin/rlp"
	"github.com/UltronGlow/UltronGlow-Origin/rpc"
	"math/big"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the delegated-proof-of-stake scheme.
type API struct {
	chain consensus.ChainHeaderReader
	alien *Alien
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)

}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
}

// GetSnapshotAtNumber retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtNumber(number uint64) (*Snapshot, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
}

// GetSnapshotByHeaderTime retrieves the state snapshot by timestamp of header.
// snapshot.header.time <= targetTime < snapshot.header.time + period
// todo: add confirm headertime in return snapshot, to minimize the request from side chain
func (api *API) GetSnapshotByHeaderTime(targetTime uint64, scHash common.Hash) (*Snapshot, error) {
	header := api.chain.CurrentHeader()
	if header == nil {
		return nil, errUnknownBlock
	}
	period := new(big.Int).SetUint64(api.chain.Config().Alien.Period)
	target := new(big.Int).SetUint64(targetTime)
	ceil := new(big.Int).Add(new(big.Int).SetUint64(header.Time), period)
	if target.Cmp(ceil) > 0 {
		target = new(big.Int).SetUint64(header.Time)
	}

	minN := new(big.Int).SetUint64(api.chain.Config().Alien.MaxSignerCount)
	maxN := new(big.Int).Set(header.Number)
	nextN := new(big.Int).SetInt64(0)
	isNext := false
	for {
		ceil = new(big.Int).Add(new(big.Int).SetUint64(header.Time), period)
		if target.Cmp(new(big.Int).SetUint64(header.Time)) >= 0 && target.Cmp(ceil) < 0 {
			snap, err := api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)

			// replace coinbase by signer settings
			var scSigners    []*common.Address
			//for _, signer := range snap.Signers {
			//	replaced := false
			//	if _, ok := snap.SCCoinbase[*signer]; ok {
			//		if addr, ok := snap.SCCoinbase[*signer][scHash]; ok {
			//			replaced = true
			//			scSigners = append(scSigners, &addr)
			//		}
			//	}
			//	if !replaced {
			//		scSigners = append(scSigners, signer)
			//	}
			//}
			for signer, _ := range snap.SCCoinbase[scHash] {
				scSigners = append(scSigners, &signer)
			}
			mcs := Snapshot{
				LoopStartTime: snap.LoopStartTime,
				Period: snap.Period,
				Signers: scSigners,
				Number: snap.Number,
				SCFULBalance: make(map[common.Address]*big.Int),
				SCMinerRevenue: make(map[common.Address]common.Address),
				SCFlowPledge: make(map[common.Address]bool),
			}
			for address, item := range snap.FULBalance {
				balance := new(big.Int).Set(item.Balance)
				for sc, cost := range item.CostTotal {
					if sc.String() == scHash.String() {
						continue
					}
					balance = new(big.Int).Sub(balance,cost)
					if 0 >= balance.Cmp(big.NewInt(0)) {
						break
					}
				}
				mcs.SCFULBalance[address] = balance
			}
			for address, revenue := range snap.RevenueFlow {
				mcs.SCMinerRevenue[address] = revenue.RevenueAddress
			}
			for address, pledge := range snap.FlowPledge {
				if 0 == pledge.StartHigh {
					mcs.SCFlowPledge[address] = true
				}
			}
			if _, ok := snap.SCNoticeMap[scHash]; ok {
				mcs.SCNoticeMap = make(map[common.Hash]*CCNotice)
				mcs.SCNoticeMap[scHash] = snap.SCNoticeMap[scHash]
			}
			return &mcs, err
		} else {
			if minNext := new(big.Int).Add(minN, big.NewInt(1)); maxN.Cmp(minN) == 0 || maxN.Cmp(minNext) == 0 {
				if !isNext && maxN.Cmp(minNext) == 0 {
					var maxHeaderTime, minHeaderTime *big.Int
					maxH := api.chain.GetHeaderByNumber(maxN.Uint64())
					if maxH != nil {
						maxHeaderTime = new(big.Int).SetUint64(maxH.Time)
					} else {
						break
					}
					minH := api.chain.GetHeaderByNumber(minN.Uint64())
					if minH != nil {
						minHeaderTime = new(big.Int).SetUint64(minH.Time)
					} else {
						break
					}
					period = period.Sub(maxHeaderTime, minHeaderTime)
					isNext = true
				} else {
					break
				}
			}
			// calculate next number
			nextN.Sub(target, new(big.Int).SetUint64(header.Time))
			nextN.Div(nextN, period)
			nextN.Add(nextN, header.Number)

			// if nextN beyond the [minN,maxN] then set nextN = (min+max)/2
			if nextN.Cmp(maxN) >= 0 || nextN.Cmp(minN) <= 0 {
				nextN.Add(maxN, minN)
				nextN.Div(nextN, big.NewInt(2))
			}
			// get new header
			header = api.chain.GetHeaderByNumber(nextN.Uint64())
			if header == nil {
				break
			}
			// update maxN & minN
			if new(big.Int).SetUint64(header.Time).Cmp(target) >= 0 {
				if header.Number.Cmp(maxN) < 0 {
					maxN.Set(header.Number)
				}
			} else if new(big.Int).SetUint64(header.Time).Cmp(target) <= 0 {
				if header.Number.Cmp(minN) > 0 {
					minN.Set(header.Number)
				}
			}

		}
	}
	return nil, errUnknownBlock
}

//y add method
func (api *API) GetSnapshotSignerAtNumber(number uint64) (*SnapshotSign, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
	if err != nil {
		log.Warn("Fail to GetSnapshotSignAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	snapshotSign := &SnapshotSign{
		LoopStartTime:snapshot.LoopStartTime,
		Signers: snapshot.Signers,
		Punished: snapshot.Punished,
	}
	return snapshotSign, err
}


type SnapshotSign struct {
	LoopStartTime   uint64                                              `json:"loopStartTime"`
	Signers         []*common.Address                                   `json:"signers"`
	Punished        map[common.Address]uint64                           `json:"punished"`
}


func (api *API) GetSnapshotReleaseAtNumber(number uint64,part string) (*SnapshotRelease, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
	if err != nil {
		log.Warn("Fail to GetSnapshotSignAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	snapshotRelease := &SnapshotRelease{
		CandidatePledge:make(map[common.Address]*PledgeItem),
		FlowPledge: make(map[common.Address]*PledgeItem),
		FlowRevenue: make(map[common.Address]*LockBalanceData),
	}
	if part!=""{
		if part =="candidatepledge"{
			snapshotRelease.CandidatePledge=snapshot.CandidatePledge
		}else if part =="flowminerpledge"{
			snapshotRelease.FlowPledge=snapshot.FlowPledge
		}else if part =="rewardlock"{
			snapshotRelease.appendFRlockData(snapshot.FlowRevenue.RewardLock,api.alien.db)
		}else if part =="flowlock"{
			snapshotRelease.appendFRlockData(snapshot.FlowRevenue.FlowLock,api.alien.db)
		}else if part =="bandwidthlock"{
			snapshotRelease.appendFRlockData(snapshot.FlowRevenue.BandwidthLock,api.alien.db)
		}
	}else{
		snapshotRelease.CandidatePledge=snapshot.CandidatePledge
		snapshotRelease.FlowPledge=snapshot.FlowPledge
		snapshotRelease.appendFRlockData(snapshot.FlowRevenue.RewardLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.FlowRevenue.FlowLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.FlowRevenue.BandwidthLock,api.alien.db)
	}
	return snapshotRelease, err
}

func (s *SnapshotRelease) appendFRItems(items []*PledgeItem) {
	for _, item := range items {
		if _, ok := s.FlowRevenue[item.TargetAddress]; !ok {
			s.FlowRevenue[item.TargetAddress] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		flowRevenusTarget := s.FlowRevenue[item.TargetAddress]
		if _, ok := flowRevenusTarget.LockBalance[item.StartHigh]; !ok {
			flowRevenusTarget.LockBalance[item.StartHigh] = make(map[uint32]*PledgeItem)
		}
		lockBalance := flowRevenusTarget.LockBalance[item.StartHigh]
		lockBalance[item.PledgeType] = item
	}
}

func (sr *SnapshotRelease) appendFR(FlowRevenue map[common.Address]*LockBalanceData) (error) {
	fr1:=FlowRevenue
	for t1, item1 := range fr1 {
		if _, ok := sr.FlowRevenue[t1]; !ok {
			sr.FlowRevenue[t1] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		rewardBalance:=item1.RewardBalance
		for t2, item2 := range rewardBalance {
			sr.FlowRevenue[t1].RewardBalance[t2]=item2
		}
		lockBalance:=item1.LockBalance
		for t3, item3 := range lockBalance {
			if _, ok := sr.FlowRevenue[t1].LockBalance[t3]; !ok {
				sr.FlowRevenue[t1].LockBalance[t3] = make(map[uint32]*PledgeItem)
			}
			sr.FlowRevenue[t1].LockBalance[t3]=item3
		}
	}
	return nil
}


func (sr *SnapshotRelease) appendFRlockData(lockData *LockData,db ethdb.Database) (error) {
	sr.appendFR(lockData.FlowRevenue)
	items, err := lockData.loadCacheL1(db)
	if err == nil {
		sr.appendFRItems(items)
	}
	items, err = lockData.loadCacheL2(db)
	if err == nil {
		sr.appendFRItems(items)
	}
	return nil
}


type SnapshotRelease struct {
	CandidatePledge map[common.Address]*PledgeItem                      `json:"candidatepledge"`
	FlowPledge      map[common.Address]*PledgeItem                      `json:"flowminerpledge"`
	FlowRevenue     map[common.Address]*LockBalanceData                 `json:"flowrevenve"`
}

func (api *API) GetSnapshotFlowAtNumber(number uint64) (*SnapshotFlow, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	lockReward:=make([]FlowRecord,0)
	if len(headerExtra.LockReward)>0 {
		for _, item := range headerExtra.LockReward {
			if(item.IsReward==sscEnumFlwReward){
				lockReward=append(lockReward,FlowRecord{
					Target: item.Target,
					Amount: item.Amount,
					FlowValue1: item.FlowValue1,
					FlowValue2: item.FlowValue2,
				})
			}
		}
	}
	snapshotFlow := &SnapshotFlow{
		LockReward: lockReward,
	}
	return snapshotFlow, err
}

type SnapshotFlow struct {
	LockReward  []FlowRecord `json:"flowrecords"`
}

type FlowRecord struct {
	Target   common.Address
	Amount   *big.Int
	FlowValue1 uint64 `json:"realFlowvalue"`
	FlowValue2 uint64 `json:"validFlowvalue"`
}

func (api *API) GetSnapshotFlowMinerAtNumber(number uint64) (*SnapshotFlowMiner, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
	if err != nil {
		log.Warn("Fail to GetSnapshotFlowMinerAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	flowMiner := &SnapshotFlowMiner{
		DayStartTime:snapshot.FlowMiner.DayStartTime,
		FlowMinerPrevTotal: snapshot.FlowMiner.FlowMinerPrevTotal,
		FlowMiner: snapshot.FlowMiner.FlowMiner,
		FlowMinerPrev:snapshot.FlowMiner.FlowMinerPrev,
		FlowMinerReport:[]*FlowMinerReport{},
		FlowMinerPrevReport:[]*FlowMinerReport{},
	}
	fMiner:=snapshot.FlowMiner
	db:=api.alien.db
	items:=flowMiner.loadFlowMinerCache(fMiner,fMiner.FlowMinerCache,db)
	flowMiner.FlowMinerReport=append(flowMiner.FlowMinerReport,items...)
	items=flowMiner.loadFlowMinerCache(fMiner,fMiner.FlowMinerPrevCache,db)
	flowMiner.FlowMinerPrevReport=append(flowMiner.FlowMinerPrevReport,items...)
	return flowMiner, err
}


type SnapshotFlowMiner struct {
	DayStartTime       uint64                                              `json:"dayStartTime"`
	FlowMinerPrevTotal uint64                                              `json:"flowminerPrevTotal"`
	FlowMiner          map[common.Address]map[common.Hash]*FlowMinerReport `json:"flowminerCurr"`
	FlowMinerReport    []*FlowMinerReport `json:"flowminerReport"`
	FlowMinerPrev      map[common.Address]map[common.Hash]*FlowMinerReport `json:"flowminerPrev"`
	FlowMinerPrevReport    []*FlowMinerReport `json:"flowminerPrevReport"`
}

func (sf *SnapshotFlowMiner) loadFlowMinerCache(fMiner *FlowMinerSnap,flowMinerCache []string,db ethdb.Database) ([]*FlowMinerReport) {
	item:=[]*FlowMinerReport{}
	for _, key := range flowMinerCache {
		flows, err := fMiner.load(db, key)
		if err != nil {
			log.Warn("appendFlowMinerCache load cache error", "key", key, "err", err)
			continue
		}
		item=append(item,flows...)
	}
	return item
}



func (api *API) GetSnapshotFlowReportAtNumber(number uint64) (*SnapshotFlowReport, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	flowReport:=make([]MinerFlowReportRecord,0)
	if len(headerExtra.FlowReport)>0 {
		flowReport=append(flowReport,headerExtra.FlowReport...)
	}
	snapshotFlowReport := &SnapshotFlowReport{
		FlowReport: flowReport,
	}
	return snapshotFlowReport, err
}

type SnapshotFlowReport struct {
	FlowReport []MinerFlowReportRecord `json:"flowreport"`
}


func (api *API) GetLockRewardAtNumber(number uint64) ([]LockRewardRecord, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	LockReward:=make([]LockRewardRecord,0)
	if len(headerExtra.LockReward)>0 {
		LockReward=append(LockReward,headerExtra.LockReward...)
	}
	return LockReward, err
}
