package alien

import (
	"errors"
	"fmt"
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/consensus"
	"math/big"
	"strconv"
)

const (
	lr_s="LockReward"
	en_s="ExchangeNFC"
	db_s="DeviceBind"
	cpl_s="CandidatePledge"
	cp_s="CandidatePunish"
	ms_s="MinerStake"
	cb_s="ClaimedBandwidth"
	bp_s="BandwidthPunish"
	cd_s="ConfigDeposit"
	ci_s="ConfigISPQOS"
	lp_s="LockParameters"
	ma_s="ManagerAddress"
	gp_s="GrantProfit"
	fr_s="FlowReport"
	mfrt_s="MinerFlowReportItem"
)
func verifyHeaderExtern(currentExtra *HeaderExtra, verifyExtra *HeaderExtra) error {

	//ExchangeNFC               []ExchangeNFCRecord
	err := verifyExchangeNFC(currentExtra.ExchangeNFC, verifyExtra.ExchangeNFC)
	if err != nil {
		return err
	}
	//LockReward                []LockRewardRecord
	err = verifyLockReward(currentExtra.LockReward, verifyExtra.LockReward)
	if err != nil {
		return err
	}

	//DeviceBind                []DeviceBindRecord
	err = verifyDeviceBind(currentExtra.DeviceBind, verifyExtra.DeviceBind)
	if err != nil {
		return err
	}

	//CandidatePledge           []CandidatePledgeRecord
	err = verifyCandidatePledge(currentExtra.CandidatePledge, verifyExtra.CandidatePledge)
	if err != nil {
		return err
	}
	//CandidatePunish           []CandidatePunishRecord
	err = verifyCandidatePunish(currentExtra.CandidatePunish, verifyExtra.CandidatePunish)
	if err != nil {
		return err
	}
	//MinerStake                []MinerStakeRecord
	err = verifyMinerStake(currentExtra.MinerStake, verifyExtra.MinerStake)
	if err != nil {
		return err
	}

	//CandidateExit             []common.Address
	err = verifyExit(currentExtra.CandidateExit, verifyExtra.CandidateExit,"CandidateExit")
	if err != nil {
		return err
	}

	//ClaimedBandwidth          []ClaimedBandwidthRecord
	err = verifyClaimedBandwidth(currentExtra.ClaimedBandwidth, verifyExtra.ClaimedBandwidth)
	if err != nil {
		return err
	}

	//FlowMinerExit             []common.Address
	err = verifyExit(currentExtra.FlowMinerExit, verifyExtra.FlowMinerExit,"FlowMinerExit")
	if err != nil {
		return err
	}

	//BandwidthPunish           []BandwidthPunishRecord
	err = verifyBandwidthPunish(currentExtra.BandwidthPunish, verifyExtra.BandwidthPunish)
	if err != nil {
		return err
	}

	//ConfigExchRate            uint32
	err = verifyUint32Config(currentExtra.ConfigExchRate, verifyExtra.ConfigExchRate,"ConfigExchRate")
	if err != nil {
		return err
	}
	//ConfigOffLine             uint32
	err = verifyUint32Config(currentExtra.ConfigOffLine, verifyExtra.ConfigOffLine,"ConfigOffLine")
	if err != nil {
		return err
	}

	//ConfigDeposit             []ConfigDepositRecord
	err = verifyConfigDeposit(currentExtra.ConfigDeposit, verifyExtra.ConfigDeposit)
	if err != nil {
		return err
	}

	//ConfigISPQOS              []ISPQOSRecord
	err = verifyConfigISPQOS(currentExtra.ConfigISPQOS, verifyExtra.ConfigISPQOS)
	if err != nil {
		return err
	}

	//LockParameters            []LockParameterRecord
	err = verifyLockParameters(currentExtra.LockParameters, verifyExtra.LockParameters)
	if err != nil {
		return err
	}

	//ManagerAddress            []ManagerAddressRecord
	err = verifyManagerAddress(currentExtra.ManagerAddress, verifyExtra.ManagerAddress)
	if err != nil {
		return err
	}
	//FlowHarvest               *big.Int
	err = verifyFlowHarvest(currentExtra.FlowHarvest, verifyExtra.FlowHarvest)
	if err != nil {
		return err
	}
	//GrantProfit               []consensus.GrantProfitRecord
	err = verifyGrantProfit(currentExtra.GrantProfit, verifyExtra.GrantProfit)
	if err != nil {
		return err
	}

	//FlowReport                []MinerFlowReportRecord
	err = verifyFlowReport(currentExtra.FlowReport, verifyExtra.FlowReport)
	if err != nil {
		return err
	}
	return nil
}


func verifyUint32Config(current uint32, verify uint32,name string) error {
	if current!=verify{
		s:=strconv.FormatUint(uint64(current), 10)
		s2:=strconv.FormatUint(uint64(verify), 10)
		return errors.New("Compare "+name+", current is "+s+". but verify is "+s2)
	}
	return nil
}


func verifyLockReward(current []LockRewardRecord, verify []LockRewardRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(lr_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(lr_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(lr_s,len(current),len(verify) )
	}
	if len(current)==0{
		return nil
	}
	err:=compareLockReward(current,verify)
	if err!=nil{
		return err
	}
	err=compareLockReward(verify,current)
	if err!=nil{
		return err
	}
	return nil
}
func compareLockReward(a []LockRewardRecord, b []LockRewardRecord) error{
	b2:= make([]LockRewardRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Amount.Cmp(v.Amount) == 0  && c.FlowValue1 == v.FlowValue1 && c.FlowValue2 == v.FlowValue2 && c.IsReward == v.IsReward && c.Target == v.Target  {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(lr_s,c)
		}
	}
	return nil
}


func verifyExchangeNFC(current []ExchangeNFCRecord, verify []ExchangeNFCRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(en_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(en_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(en_s,len(current),len(verify))
	}
	if len(current)== 0{
		return nil
	}
	err:=compareExchangeNFC(current,verify)
	if err!=nil{
		return err
	}
	err=compareExchangeNFC(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareExchangeNFC(a []ExchangeNFCRecord, b []ExchangeNFCRecord) error{
	b2:= make([]ExchangeNFCRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Amount.Cmp(v.Amount)==0 {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(en_s,c)
		}
	}
	return nil
}

func verifyDeviceBind(current []DeviceBindRecord, verify []DeviceBindRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(db_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(db_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(db_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareDeviceBind(current,verify)
	if err!=nil{
		return err
	}
	err=compareDeviceBind(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareDeviceBind(a []DeviceBindRecord, b []DeviceBindRecord) error{
	b2:= make([]DeviceBindRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Device == v.Device  && c.Revenue == v.Revenue && c.Contract == v.Contract && c.MultiSign == v.MultiSign && c.Type == v.Type  && c.Bind == v.Bind {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(db_s,c)
		}
	}
	return nil
}

func verifyCandidatePledge(current []CandidatePledgeRecord, verify []CandidatePledgeRecord) error {

	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(cpl_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(cpl_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(cpl_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareCandidatePledge(current,verify)
	if err!=nil{
		return err
	}
	err=compareCandidatePledge(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareCandidatePledge(a []CandidatePledgeRecord, b []CandidatePledgeRecord) error{
	b2:= make([]CandidatePledgeRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target  && c.Amount.Cmp(v.Amount)==0 {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cpl_s,c)
		}
	}
	return nil
}


func verifyCandidatePunish(current []CandidatePunishRecord, verify []CandidatePunishRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(cp_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(cp_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(cp_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareCandidatePunish(current,verify)
	if err!=nil{
		return err
	}
	err=compareCandidatePunish(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareCandidatePunish(a []CandidatePunishRecord, b []CandidatePunishRecord) error{
	b2:= make([]CandidatePunishRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target  && c.Amount.Cmp(v.Amount)==0  && c.Credit==v.Credit{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cp_s,c)
		}
	}
	return nil
}

func verifyMinerStake(current []MinerStakeRecord, verify []MinerStakeRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(ms_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(ms_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(ms_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareMinerStake(current,verify)
	if err!=nil{
		return err
	}
	err=compareMinerStake(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareMinerStake(a []MinerStakeRecord, b []MinerStakeRecord) error{
	b2:= make([]MinerStakeRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target  && c.Stake.Cmp(v.Stake)==0{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ms_s,c)
		}
	}
	return nil
}

func verifyExit(current []common.Address, verify []common.Address,name string) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(name)
	}
	if current != nil && verify == nil {
		return errorsMsg2(name)
	}
	if len(current) != len(verify) {
		return errorsMsg3(name,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareExit(current,verify,name)
	if err!=nil{
		return err
	}
	err=compareExit(verify,current,name)
	if err!=nil{
		return err
	}
	return nil
}

func compareExit(a []common.Address, b []common.Address,name string) error {
	b2:= make([]common.Address, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c == v {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(name,c)
		}
	}
	return nil
}

func verifyClaimedBandwidth(current []ClaimedBandwidthRecord, verify []ClaimedBandwidthRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(cb_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(cb_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(cb_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareClaimedBandwidth(current,verify)
	if err!=nil{
		return err
	}
	err=compareClaimedBandwidth(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareClaimedBandwidth(a []ClaimedBandwidthRecord, b []ClaimedBandwidthRecord) error {
	b2:= make([]ClaimedBandwidthRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.Amount.Cmp(v.Amount)==0&&c.ISPQosID==v.ISPQosID&&c.Bandwidth==v.Bandwidth {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cb_s,c)
		}
	}
	return nil
}


func verifyBandwidthPunish(current []BandwidthPunishRecord, verify []BandwidthPunishRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(bp_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(bp_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(bp_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareBandwidthPunish(current,verify)
	if err!=nil{
		return err
	}
	err=compareBandwidthPunish(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareBandwidthPunish(a []BandwidthPunishRecord, b []BandwidthPunishRecord) error {
	b2:= make([]BandwidthPunishRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.WdthPnsh==v.WdthPnsh {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(bp_s,c)
		}
	}
	return nil
}

func verifyConfigDeposit(current []ConfigDepositRecord, verify []ConfigDepositRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(cd_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(cd_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(cd_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareConfigDeposit(current,verify)
	if err!=nil{
		return err
	}
	err=compareConfigDeposit(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareConfigDeposit(a []ConfigDepositRecord, b []ConfigDepositRecord) error {
	b2:= make([]ConfigDepositRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Who == v.Who&&c.Amount.Cmp(v.Amount)==0 {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cd_s,c)
		}
	}
	return nil
}

func verifyConfigISPQOS(current []ISPQOSRecord, verify []ISPQOSRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(ci_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(ci_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(ci_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareConfigISPQOS(current,verify)
	if err!=nil{
		return err
	}
	err=compareConfigISPQOS(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareConfigISPQOS(a []ISPQOSRecord, b []ISPQOSRecord) error {
	b2:= make([]ISPQOSRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.ISPID == v.ISPID&&c.QOS==v.QOS {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ci_s,c)
		}
	}
	return nil
}

func verifyLockParameters(current []LockParameterRecord, verify []LockParameterRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(lp_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(lp_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(lp_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareLockParameters(current,verify)
	if err!=nil{
		return err
	}
	err=compareLockParameters(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareLockParameters(a []LockParameterRecord, b []LockParameterRecord) error {
	b2:= make([]LockParameterRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.LockPeriod == v.LockPeriod&&c.RlsPeriod==v.RlsPeriod &&c.Interval==v.Interval&&c.Who==v.Who{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(lp_s,c)
		}
	}
	return nil
}

func verifyManagerAddress(current []ManagerAddressRecord, verify []ManagerAddressRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(ma_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(ma_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(ma_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareManagerAddress(current,verify)
	if err!=nil{
		return err
	}
	err=compareManagerAddress(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareManagerAddress(a []ManagerAddressRecord, b []ManagerAddressRecord) error {
	b2:= make([]ManagerAddressRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.Who==v.Who{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ma_s,c)
		}
	}
	return nil
}

func verifyFlowHarvest(current *big.Int, verify *big.Int) error {
	fh_s:="FlowHarvest"
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(fh_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(fh_s)
	}
	if current != nil && verify != nil && current.Cmp(verify)!=0 {
		return errors.New("Compare "+fh_s+", current is "+current.String()+". but verify is "+verify.String())
	}
	return nil
}

func verifyGrantProfit(current []consensus.GrantProfitRecord, verify []consensus.GrantProfitRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(gp_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(gp_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(gp_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareGrantProfit(current,verify)
	if err!=nil{
		return err
	}
	err=compareGrantProfit(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareGrantProfit(a []consensus.GrantProfitRecord, b []consensus.GrantProfitRecord)error {
	b2:= make([]consensus.GrantProfitRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Which == v.Which&&c.MinerAddress==v.MinerAddress&&c.BlockNumber==v.BlockNumber&&c.Amount.Cmp(v.Amount)==0&&c.RevenueAddress==v.RevenueAddress&&c.RevenueContract==v.RevenueContract&&c.MultiSignature==v.MultiSignature{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(gp_s,c)
		}
	}
	return nil
}


func verifyFlowReport(current []MinerFlowReportRecord, verify []MinerFlowReportRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(fr_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(fr_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(fr_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareFlowReport(current,verify)
	if err!=nil{
		return err
	}
	err=compareFlowReport(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareFlowReport(a []MinerFlowReportRecord, b []MinerFlowReportRecord)error {
	b2:= make([]MinerFlowReportRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.ChainHash == v.ChainHash&&c.ReportTime==v.ReportTime{
				if err:=verifyMinerFlowReportItem(c.ReportContent,v.ReportContent);err==nil{
					find = true
					b2=append(b2[:i],b2[i+1:]...)
					break
				}
			}
		}
		if !find {
			return errorsMsg4(fr_s,c)
		}
	}
	return nil
}

func verifyMinerFlowReportItem(current []MinerFlowReportItem, verify []MinerFlowReportItem) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(mfrt_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(mfrt_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(mfrt_s,len(current),len(verify))
	}
	if len(current)==0{
		return nil
	}
	err:=compareMinerFlowReportItem(current,verify)
	if err!=nil{
		return err
	}
	err=compareMinerFlowReportItem(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareMinerFlowReportItem(a []MinerFlowReportItem, b []MinerFlowReportItem)error {
	b2:= make([]MinerFlowReportItem, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.ReportNumber==v.ReportNumber&&c.FlowValue1==v.FlowValue1&&c.FlowValue2==v.FlowValue2{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(mfrt_s,c)
		}
	}
	return nil
}

func errorsMsg1(name string) error {
	return errors.New("Compare "+name+" , current is nil. but verify is not nil")
}
func errorsMsg2(name string) error {
	return errors.New("Compare "+name+" , current is not nil. but verify is nil")
}
func errorsMsg3(name string,lenc int,lenv int) error {
	return errors.New(fmt.Sprintf("Compare "+name+", The array length is not equals. the current length is %d. the verify length is %d", lenc, lenv))
}
func errorsMsg4(name string,c interface{}) error {
	return errors.New(fmt.Sprintf("Compare "+name+", can't find %v in verify data", c))
}

