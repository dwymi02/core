package actions

import (
	"bytes"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

// btc 转账 （amt 单位 聪）
func DoSimpleSatoshiTransferFromChainState(state interfaces.ChainStateOperation, addr1 fields.Address, addr2 fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return fmt.Errorf("Satoshi transfer amount is empty") // 不允许转账数量为0
	}
	bls1 := state.Balance(addr1)
	if bls1 == nil {
		return fmt.Errorf("Satoshi need %d but empty.", sat)
	}
	sat1 := bls1.Satoshi
	// 检查余额
	if uint64(sat1) < uint64(sat) {
		return fmt.Errorf("Address %s satoshi %d not enough, need at least %d.", addr1.ToReadable(), sat1, sat)
	}
	// 检查自己转给自己
	if bytes.Compare(addr1, addr2) == 0 {
		return nil // 可以自己转给自己，不改变状态，先检查余额充足，白费手续费
	}
	bls2 := state.Balance(addr2)
	if bls2 == nil {
		bls2 = stores.NewEmptyBalance() // create satoshi store
	}
	sat2 := bls2.Satoshi
	bls1.Satoshi = fields.Satoshi(uint64(sat1) - uint64(sat)) // 扣除
	bse1 := state.BalanceSet(addr1, bls1)
	if bse1 != nil {
		return bse1
	}
	bls2.Satoshi = fields.Satoshi(uint64(sat2) + uint64(sat)) // 增加
	bse2 := state.BalanceSet(addr2, bls2)
	if bse2 != nil {
		return bse2
	}
	// return ok
	return nil
}

// 单纯增加 BTC 余额 （amt 单位 聪）
func DoAddSatoshiFromChainState(state interfaces.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // 数量为0，直接成功
	}
	blssto := state.Balance(addr)
	if blssto == nil {
		blssto = stores.NewEmptyBalance() // first create account
	}
	basesat := blssto.Satoshi
	newsat := uint64(basesat) + uint64(sat) // 增加
	// 新余额
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}

// 单纯扣除 BTC 余额 （amt 单位 聪）
func DoSubSatoshiFromChainState(state interfaces.ChainStateOperation, addr fields.Address, sat fields.Satoshi) error {
	if sat == 0 {
		return nil // 数量为0，直接成功
	}
	blssto := state.Balance(addr)
	if blssto == nil {
		return fmt.Errorf("address %s satoshi need %d but empty.", addr.ToReadable(), sat)
	}
	basesat := blssto.Satoshi
	// 检查余额
	if uint64(basesat) < uint64(sat) {
		return fmt.Errorf("address %s satoshi %d not enough, need more %d.", addr.ToReadable(), basesat, sat)
	}
	newsat := uint64(basesat) - uint64(sat) // 扣除
	blssto.Satoshi = fields.Satoshi(newsat)
	bserr := state.BalanceSet(addr, blssto)
	if bserr != nil {
		return bserr
	}
	return nil
}
