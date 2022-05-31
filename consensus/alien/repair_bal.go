package alien

import (
	"github.com/UltronGlow/UltronGlow-Origin/common"
	"github.com/UltronGlow/UltronGlow-Origin/core/state"
	"github.com/UltronGlow/UltronGlow-Origin/log"
	"math/big"
)

const (
	RepairBalanceNumber=1170523
)
func (a *Alien) RepairBal(state *state.StateDB,number uint64){
	if number==RepairBalanceNumber{
		illegalAccounts:=[]string{
			"uxec5bfd0c33e25c09ffd8e4720e73976e60f99c4a",
			"ux4b971f450d430ee72e4c38d0798ac5a99b0bfd8d",
			"uxd82585d44e2b3499fbc16d908f86d0da7c68e08b",
			"uxeb5a23bd9b388b0d5b6325936b5369daafed2d34",
			"ux02a919271e24e7fdfb418960ce1f43f1e022421c",
			"ux91ffdfb80b115ce2d20ed35bce1f1dd52e288858",
			"uxfd38dec028af4832d89a8b99df1e84f4dcb365b4",
			"uxa8d0405e5556cac9100c24e69ec257075113ea8f",
			"ux81c93d4f1818d15575a02e79cc8ea8b0baec4094",
			"uxc5d95357f16070b931e4978d83a95dd20d462519",
			"uxa94cf9c2fb44df0e6015a16063f3ca688294801a",
			"uxd410b62e8e7d0cc1c6f7454318f98111a3a6b55a",
			"uxe9e98ef1baf846edbfe46c69c8a2ad4941557bcc",
			"ux48a6cd019da4cd2ad8fdc9f45caa04d69a31c417",
			"uxb8cb2783f1e7f4003d15bf773ad9a28093ebc244",
		}
		illTolBal:=common.Big0
		for _,illAcc:=range illegalAccounts{
			illBal:=state.GetBalance(common.HexToAddress(illAcc))
			illTolBal=new(big.Int).Add(illTolBal,illBal)
		}
		for _,illAcc:=range illegalAccounts{
			state.SetBalance(common.HexToAddress(illAcc),common.Big0)
		}
		targetAccount:=common.BigToAddress(big.NewInt(0))
		state.AddBalance(targetAccount,illTolBal)
		log.Info("RepairBal", "number", number, "targetAccount", targetAccount,"addBal", illTolBal)

	}
}
