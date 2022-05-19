package alien

import (
	"encoding/json"
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/ethdb"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"math/big"
	"strconv"
)

const (
	utgSRTExch = "Exch"
    utgSRTIndex ="srt-index-%d"
	utgBalKey ="srt-bal-%d"
)

type ExchangeSRTRecord struct {
	Target common.Address `json:"target"`
	Amount *big.Int `json:"amount"`
}

type SRTIndex struct {

}
func (srtI *SRTIndex) updateExchangeSRT(exchangeSRT []ExchangeSRTRecord, number uint64, db ethdb.Database) {
		indexNum:=srtI.loadSRTIndex(number-1,db)
	    if len(exchangeSRT)==0||exchangeSRT==nil{
			if indexNum==uint64(0){
				srtI.deleteSRTIndex(number,db)
				srtI.deleteSRTBal(number,db)
				return
			}else{
				srtI.storeSRTIndex(number,indexNum,db)
			}
		}else{
			srtBalance:=make(map[common.Address]*big.Int)
			if indexNum!=uint64(0){
				srtBalance,_=srtI.loadSRTBal(indexNum,db)
			}
			for _, item := range exchangeSRT {
				if balance, ok := srtBalance[item.Target]; !ok {
					srtBalance[item.Target] = new(big.Int).Set(item.Amount)
				} else {
					srtBalance[item.Target] = new(big.Int).Add(balance, item.Amount)
				}
			}
			srtI.storeSRTIndex(number,number,db)
			srtI.storeSRTBal(srtBalance,number,db)
		}

	return
}

func (srtI *SRTIndex) loadSRTBal(number uint64,db ethdb.Database) (map[common.Address]*big.Int, error) {
	key := fmt.Sprintf(utgBalKey, number)
	blob, err := db.Get([]byte(key))
	if err != nil {
		log.Info("loadSRTBal Get", "err", err)
		return nil, err
	}
	Balances := make(map[common.Address]*big.Int)
	if err := json.Unmarshal(blob, &Balances); err != nil {
		log.Info("loadSRTBal Unmarshal", "err", err)
		return nil, err
	}
	return Balances,nil
}


func (srtI *SRTIndex) storeSRTBal(srtBalance map[common.Address]*big.Int,number uint64,db ethdb.Database) (error) {
	key := fmt.Sprintf(utgBalKey, number)
	blob, err := json.Marshal(srtBalance)
	if err != nil {
		return err
	}
	err = db.Put([]byte(key), blob)
	if err != nil {
		return err
	}
	return nil
}

func (srtI *SRTIndex) deleteSRTBal(number uint64,db ethdb.Database) (error) {
	key := fmt.Sprintf(utgBalKey, number)
	err := db.Delete([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

func (srtI *SRTIndex) loadSRTIndex(number uint64,db ethdb.Database) (uint64) {
	key := fmt.Sprintf(utgSRTIndex, number)
	blob, err := db.Get([]byte(key))
	if err != nil {
		log.Info("loadSRTIndex Get", "err", err)
		return 0
	}
	var indexNumber uint64
	if indexNumber,_ = strconv.ParseUint(string(blob), 10, 64); err != nil {
		log.Info("loadSRTIndex ParseUint", "err", err)
		return 0
	}
	return indexNumber
}

func (srtI *SRTIndex) storeSRTIndex(number uint64,indexNum uint64,db ethdb.Database) (error) {
	key := fmt.Sprintf(utgSRTIndex, number)
	err := db.Put([]byte(key), []byte(strconv.FormatUint(indexNum,10)))
	if err != nil {
		return err
	}
	return nil
}

func (srtI *SRTIndex) deleteSRTIndex(number uint64,db ethdb.Database) (error) {
	key := fmt.Sprintf(utgSRTIndex, number)
	err := db.Delete([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

func (srtI *SRTIndex) getSRTBalAtNumber(number uint64,db ethdb.Database) (map[common.Address]*big.Int, error) {
	indexNumber:= srtI.loadSRTIndex(number,db)
	return srtI.loadSRTBal(indexNumber,db)
}

func (srtI *SRTIndex) checkEnoughSRT(sRent []LeaseRequestRecord, rent LeaseRequestRecord, number uint64, db ethdb.Database) bool {
	srtAmount:=new(big.Int).Mul(rent.Duration,rent.Price)
	srtAmount=new(big.Int).Mul(srtAmount,rent.Capacity)
	srtAmount=new(big.Int).Div(srtAmount,gbTob)
	for _, item := range sRent {
		if item.Tenant==rent.Tenant{
			itemSrtAmount:=new(big.Int).Mul(item.Duration,item.Price)
			itemSrtAmount=new(big.Int).Mul(itemSrtAmount,item.Capacity)
			itemSrtAmount=new(big.Int).Div(itemSrtAmount,gbTob)
			srtAmount=new(big.Int).Add(srtAmount,itemSrtAmount)
		}
	}
	srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
	if balance, ok := srtBalance[rent.Tenant]; ok {
		if balance.Cmp(srtAmount)>=0 {
			return true
		}
	}
	return false
}


func (srtI *SRTIndex) checkEnoughSRTPg(sRentPg []LeasePledgeRecord, rent LeasePledgeRecord,number uint64, db ethdb.Database) bool {
	srtAmount:=rent.BurnSRTAmount
	for _, item := range sRentPg {
		if item.BurnSRTAddress==rent.BurnSRTAddress{
			srtAmount=new(big.Int).Add(srtAmount,item.BurnSRTAmount)
		}
	}
	srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
	if balance, ok := srtBalance[rent.BurnSRTAddress]; ok {
		if balance.Cmp(srtAmount)>=0 {
			return true
		}
	}
	return false
}

func (srtI *SRTIndex) checkEnoughSRTReNewPg(sRentPg []LeaseRenewalPledgeRecord, rent LeaseRenewalPledgeRecord,number uint64, db ethdb.Database) bool {
	srtAmount:=rent.BurnSRTAmount
	for _, item := range sRentPg {
		if item.BurnSRTAddress==rent.BurnSRTAddress{
			srtAmount=new(big.Int).Add(srtAmount,item.BurnSRTAmount)
		}
	}
	srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
	if balance, ok := srtBalance[rent.BurnSRTAddress]; ok {
		if balance.Cmp(srtAmount)>=0 {
			return true
		}
	}
	return false
}

func (srtI *SRTIndex) burnSRTAmount(pg []LeasePledgeRecord, number uint64, db ethdb.Database) {
	if pg!=nil&&len(pg)>0 {
		srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
		for _, item := range pg {
			if balance, ok := srtBalance[item.BurnSRTAddress]; ok {
				if balance.Cmp(item.BurnSRTAmount)>0 {
					srtBalance[item.BurnSRTAddress]=new(big.Int).Sub(balance,item.BurnSRTAmount)
				}else{
					delete(srtBalance,item.BurnSRTAddress)
				}
			}
		}
		srtI.storeSRTIndex(number,number,db)
		srtI.storeSRTBal(srtBalance,number,db)

	}
}

func (srtI *SRTIndex) checkEnoughSRTReNew(currentSRentReNew []LeaseRenewalRecord, sRentReNew LeaseRenewalRecord, number uint64, db ethdb.Database) bool {
	srtAmount:=new(big.Int).Mul(sRentReNew.Duration,sRentReNew.Price)
	srtAmount=new(big.Int).Mul(srtAmount,sRentReNew.Capacity)
	srtAmount=new(big.Int).Div(srtAmount,gbTob)
	for _, item := range currentSRentReNew {
		if item.Tenant==sRentReNew.Tenant{
			itemSrtAmount:=new(big.Int).Mul(item.Duration,item.Price)
			itemSrtAmount=new(big.Int).Mul(itemSrtAmount,item.Capacity)
			itemSrtAmount=new(big.Int).Div(itemSrtAmount,gbTob)
			srtAmount=new(big.Int).Add(srtAmount,itemSrtAmount)
		}
	}
	srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
	if balance, ok := srtBalance[sRentReNew.Tenant]; ok {
		if balance.Cmp(srtAmount)>=0 {
			return true
		}
	}
	return false

}
func (srtI *SRTIndex) burnSRTAmountReNew(pg []LeaseRenewalPledgeRecord, number uint64, db ethdb.Database) {
	if pg!=nil&&len(pg)>0 {
		srtBalance,_:=srtI.getSRTBalAtNumber(number,db)
		for _, item := range pg {
			if balance, ok := srtBalance[item.BurnSRTAddress]; ok {
				if balance.Cmp(item.BurnSRTAmount)>0 {
					srtBalance[item.BurnSRTAddress]=new(big.Int).Sub(balance,item.BurnSRTAmount)
				}else{
					delete(srtBalance,item.BurnSRTAddress)
				}
			}
		}
		srtI.storeSRTIndex(number,number,db)
		srtI.storeSRTBal(srtBalance,number,db)
	}
}