package interfaces

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// chain state 操作

type ChainStateOperation interface {
	// 数据库升级模式
	IsDatabaseVersionRebuildMode() bool

	// status
	IsInMemTxPool() bool // 否在交易池
	SetInMemTxPool(bool)

	// status
	GetPendingBlockHeight() uint64
	SetPendingBlockHeight(uint64)
	GetPendingBlockHash() fields.Hash
	SetPendingBlockHash(fields.Hash)

	GetPendingSubmitStoreDiamond() (*stores.DiamondSmelt, error)
	SetPendingSubmitStoreDiamond(*stores.DiamondSmelt) error

	// status
	SetLastestBlockHeadAndMeta(Block) error
	ReadLastestBlockHeadAndMeta() (Block, error)
	SetLastestDiamond(*stores.DiamondSmelt) error
	ReadLastestDiamond() (*stores.DiamondSmelt, error)

	UpdateSetTotalSupply(totalobj *stores.TotalSupply) error
	ReadTotalSupply() (*stores.TotalSupply, error)

	// store
	BlockStore() BlockStore
	SetBlockStore(BlockStore) error

	// state

	// query

	Balance(fields.Address) *stores.Balance
	//Satoshi(fields.Address) *stores.Satoshi
	Lockbls(fields.LockblsId) *stores.Lockbls
	Channel(fields.ChannelId) *stores.Channel
	Diamond(fields.DiamondName) *stores.Diamond
	DiamondSystemLending(fields.DiamondSyslendId) *stores.DiamondSystemLending
	BitcoinSystemLending(fields.BitcoinSyslendId) *stores.BitcoinSystemLending
	UserLending(fields.UserLendingId) *stores.UserLending
	Chaswap(fields.HashHalfChecker) *stores.Chaswap

	// operate

	BalanceSet(fields.Address, *stores.Balance) error
	BalanceDel(fields.Address) error

	//SatoshiSet(fields.Address, *stores.Satoshi) error
	//SatoshiDel(fields.Address) error

	LockblsCreate(fields.LockblsId, *stores.Lockbls) error // 创建线性锁仓
	LockblsUpdate(fields.LockblsId, *stores.Lockbls) error // 更新：释放（取出部分任意可取额度）
	LockblsDelete(fields.LockblsId) error                  // 释放完毕后自动删除

	ChannelCreate(fields.ChannelId, *stores.Channel) error
	ChannelUpdate(fields.ChannelId, *stores.Channel) error
	ChannelDelete(fields.ChannelId) error

	DiamondSet(fields.DiamondName, *stores.Diamond) error
	DiamondDel(fields.DiamondName) error

	DiamondLendingCreate(fields.DiamondSyslendId, *stores.DiamondSystemLending) error
	DiamondLendingUpdate(fields.DiamondSyslendId, *stores.DiamondSystemLending) error
	DiamondLendingDelete(fields.DiamondSyslendId) error

	BitcoinLendingCreate(fields.BitcoinSyslendId, *stores.BitcoinSystemLending) error
	BitcoinLendingUpdate(fields.BitcoinSyslendId, *stores.BitcoinSystemLending) error
	BitcoinLendingDelete(fields.BitcoinSyslendId) error

	UserLendingCreate(fields.UserLendingId, *stores.UserLending) error
	UserLendingUpdate(fields.UserLendingId, *stores.UserLending) error
	UserLendingDelete(fields.UserLendingId) error

	ChaswapCreate(fields.HashHalfChecker, *stores.Chaswap) error
	ChaswapUpdate(fields.HashHalfChecker, *stores.Chaswap) error
	ChaswapDelete(fields.HashHalfChecker) error

	// movebtc
	SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error
	ReadMoveBTCTxHashByNumber(trsno uint32) ([]byte, error)
	LoadValidatedSatoshiGenesis(int64) (*stores.SatoshiGenesis, bool) // 获取已验证的BTC转移日志 & 是否需要验证

}
