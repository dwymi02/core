package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/channel"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/core/sys"
)

// 无任何单据单方面关闭通道，进入挑战期
// 资金分配按初始存入计算
type Action_22_UnilateralClosePaymentChannelByNothing struct {
	// 通道 ID
	ChannelId          fields.ChannelId // 通道id
	AssertCloseAddress fields.Address   // 单方面主张关闭的提议地址

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Kind() uint16 {
	return 22
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Size() uint32 {
	return 2 + elm.ChannelId.Size() + elm.AssertCloseAddress.Size()
}

// json api
func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var bt2, _ = elm.AssertCloseAddress.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.AssertCloseAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) RequestSignAddresses() []fields.Address {
	// 提议者必须签名
	return []fields.Address{
		elm.AssertCloseAddress,
	}
}

func (act *Action_22_UnilateralClosePaymentChannelByNothing) WriteinChainState(state interfaces.ChainStateOperation) error {
	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 检查状态（必须为开启状态）
	if paychan.IsOpening() == false {
		return fmt.Errorf("Payment Channel status is not on opening.")
	}
	// 检查两个账户地址看地址是否匹配
	addrIsLeft := paychan.LeftAddress.Equal(act.AssertCloseAddress)
	addrIsRight := paychan.RightAddress.Equal(act.AssertCloseAddress)
	if !addrIsLeft && !addrIsRight {
		return fmt.Errorf("Payment Channel <%s> address signature verify fail.", hex.EncodeToString(act.ChannelId))
	}
	// 挑战者状态
	clghei := state.GetPendingBlockHeight()
	var clgamt = fields.Amount{}
	var clgsat = fields.Satoshi(0)
	if addrIsLeft {
		clgamt = paychan.LeftAmount
		clgsat = paychan.LeftSatoshi.GetRealSatoshi()
	} else {
		clgamt = paychan.RightAmount
		clgsat = paychan.RightSatoshi.GetRealSatoshi()
	}
	// 更新至挑战期，没有账单编号
	paychan.SetChallenging(clghei, addrIsLeft, &clgamt, clgsat, 0)
	// 写入状态
	e = state.ChannelUpdate(act.ChannelId, paychan)
	if e != nil {
		return e
	}
	return nil
}

func (act *Action_22_UnilateralClosePaymentChannelByNothing) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 回退状态
	paychan.SetOpening()
	state.ChannelUpdate(act.ChannelId, paychan)
	return nil
}

func (elm *Action_22_UnilateralClosePaymentChannelByNothing) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_22_UnilateralClosePaymentChannelByNothing) IsBurning90PersentTxFees() bool {
	return false
}

/////////////////////////////////////////////////////////////

// 1. 通过中间实时对账单方面关闭通道，进入挑战期
// 2. 提供实时对账单，回应挑战，夺取对方全部金额
type Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation struct {
	// 主张者地址
	AssertAddress fields.Address
	// 对账单
	Reconciliation channel.OnChainArbitrationBasisReconciliation

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Kind() uint16 {
	return 23
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Size() uint32 {
	return 2 + elm.AssertAddress.Size() + elm.Reconciliation.Size()
}

// json api
func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.Reconciliation.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	return buffer.Bytes(), nil
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.Reconciliation.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) RequestSignAddresses() []fields.Address {
	// 检查签名
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// cid
	channelId := act.Reconciliation.GetChannelId()

	// 查询通道
	paychan := state.Channel(channelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(channelId))
	}
	// 检查两个账户地址签名，双方都检查
	// 进入挑战期还是夺取资金
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.Reconciliation)
}

func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) RecoverChainState(state interfaces.ChainStateOperation) error {

	channelId := act.Reconciliation.GetChannelId()

	// 查询通道
	paychan := state.Channel(channelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(channelId))
	}

	// 回退
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.Reconciliation)
}

func (elm *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_23_UnilateralCloseOrRespondChallengePaymentChannelByRealtimeReconciliation) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////

// 单方面结束
// 1. 通过 通道链支付 单方面关闭通道，进入挑战期
// 2. 提供通道链支付对账单，回应挑战，夺取对方全部金额
type Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody struct {
	// 主张者地址
	AssertAddress fields.Address
	// 通道整体支付数据
	ChannelChainTransferData channel.OffChainFormPaymentChannelTransfer
	// 本通道支付体数据
	ChannelChainTransferTargetProveBody channel.ChannelChainTransferProveBodyInfo

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Kind() uint16 {
	return 24
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Size() uint32 {
	return 2 + elm.AssertAddress.Size() +
		elm.ChannelChainTransferData.Size() +
		elm.ChannelChainTransferTargetProveBody.Size()
}

// json api
func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.ChannelChainTransferData.Serialize()
	var bt3, _ = elm.ChannelChainTransferTargetProveBody.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	return buffer.Bytes(), nil
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferData.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) RequestSignAddresses() []fields.Address {
	// 检查签名
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) WriteinChainState(state interfaces.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 查询通道
	paychan := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// 检查通道哈希是否正确
	hxhalf := act.ChannelChainTransferTargetProveBody.GetSignStuffHashHalfChecker()
	// 检查哈希值是否包含在列表内
	var isHashCheckOk = false
	for _, hxckr := range act.ChannelChainTransferData.ChannelTransferProveHashHalfCheckers {
		if hxhalf.Equal(hxckr) {
			isHashCheckOk = true
			break
		}
	}
	if !isHashCheckOk {
		return fmt.Errorf("ChannelChainTransferTargetProveBody hash <%s> not find.", hxhalf.ToHex())
	}

	// 检查双方通道地址是否包含在签名列表内
	lsgok := false
	rsgok := false
	for _, v := range act.ChannelChainTransferData.MustSignAddresses {
		if v.Equal(paychan.LeftAddress) {
			lsgok = true
		} else if v.Equal(paychan.RightAddress) {
			rsgok = true
		}
	}
	if !lsgok || !rsgok {
		return fmt.Errorf("Channel signature address is missing.")
	}

	// 检查所有签名是否完整和正确
	e = act.ChannelChainTransferData.CheckMustAddressAndSigns()
	if e != nil {
		return e
	}

	// 检查 进入挑战期，还是最终夺取
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// 回退
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (elm *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_24_UnilateralCloseOrRespondChallengePaymentChannelByChannelChainTransferBody) IsBurning90PersentTxFees() bool {
	return false
}

///////////////////////////////////////////////

// 单方面结束
// 1. 通过 通道&链上原子交换 单方面关闭通道，进入挑战期
// 2. 提供通道链支付对账单，回应挑战，夺取对方全部金额
type Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange struct {
	// 主张者地址
	AssertAddress fields.Address
	// 凭据
	ProveBodyHashChecker fields.HashHalfChecker
	// 对账数据
	ChannelChainTransferTargetProveBody channel.ChannelChainTransferProveBodyInfo

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Kind() uint16 {
	return 26
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Size() uint32 {
	return 2 + elm.AssertAddress.Size() +
		elm.ProveBodyHashChecker.Size() +
		elm.ChannelChainTransferTargetProveBody.Size()
}

// json api
func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.AssertAddress.Serialize()
	var bt2, _ = elm.ProveBodyHashChecker.Serialize()
	var bt3, _ = elm.ChannelChainTransferTargetProveBody.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	buffer.Write(bt2)
	buffer.Write(bt3)
	return buffer.Bytes(), nil
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.AssertAddress.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ProveBodyHashChecker.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	seek, e = elm.ChannelChainTransferTargetProveBody.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) RequestSignAddresses() []fields.Address {
	// 检查签名
	return []fields.Address{
		elm.AssertAddress,
	}
}

func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) WriteinChainState(state interfaces.ChainStateOperation) error {

	var e error

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}

	// 快速通道模式不能用来发起挑战和仲裁，只有普通模式才可以

	// 查询通道
	paychan := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// 查询互换交易
	swapex := state.Chaswap(act.ProveBodyHashChecker)
	if swapex == nil {
		return fmt.Errorf("Chaswap tranfer <%s> not find.", act.ProveBodyHashChecker.ToHex())
	}
	// 是否已经使用过
	if swapex.IsBeUsed.Check() {
		return fmt.Errorf("Chaswap tranfer <%s> already be used.", act.ProveBodyHashChecker.ToHex())
	}

	// 检查必须签名的地址是否完整和正确
	addrsmap := make(map[string]bool)
	for _, addr := range swapex.OnchainTransferFromAndMustSignAddresses {
		addrsmap[string(addr)] = true
	}
	_, hasleft := addrsmap[string(paychan.LeftAddress)]
	_, hasright := addrsmap[string(paychan.RightAddress)]
	if !hasleft || !hasright {
		return fmt.Errorf("Chaswap tranfer signature error.")
	}

	// 标记票据已使用
	swapex.IsBeUsed.Set(true)
	e = state.ChaswapUpdate(act.ProveBodyHashChecker, swapex)
	if e != nil {
		return e
	}

	// 检查 进入挑战期，还是最终夺取
	return checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.ChannelChainTransferTargetProveBody.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelChainTransferTargetProveBody.ChannelId))
	}

	// 回退使用状态
	swapex := state.Chaswap(act.ProveBodyHashChecker)
	swapex.IsBeUsed.Set(false)
	state.ChaswapCreate(act.ProveBodyHashChecker, swapex)

	// 回退状态
	return checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state, act.AssertAddress, paychan, &act.ChannelChainTransferTargetProveBody)
}

func (elm *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_26_UnilateralCloseOrRespondChallengePaymentChannelByChannelOnchainAtomicExchange) IsBurning90PersentTxFees() bool {
	return false
}

//////////////////////////////////////////////////////////////

// 检查通道进入挑战期或者最终仲裁
func checkChannelGotoChallegingOrFinalDistributionWriteinChainState(state interfaces.ChainStateOperation, assertAddress fields.Address, paychan *stores.Channel, obj channel.OnChainChannelPaymentArbitrationReconciliationBasis) error {

	channelId := obj.GetChannelId()

	// 通道不能已经关闭
	if paychan.IsClosed() {
		return fmt.Errorf("Payment Channel <%s> is closed.", hex.EncodeToString(channelId))
	}
	// 检查地址匹配你
	var assertAddressIsLeft = paychan.LeftAddress.Equal(assertAddress)
	var assertAddressIsRight = paychan.RightAddress.Equal(assertAddress)
	if !assertAddressIsLeft && !assertAddressIsRight {
		return fmt.Errorf("Payment Channel AssertAddress is not match left or right.")
	}
	// 检查两个账户地址签名，双方都检查
	e20 := obj.CheckAddressAndSign(paychan.LeftAddress, paychan.RightAddress)
	if e20 != nil {
		return e20
	}
	// 检查对账单资金数额和重用版本
	channelReuseVersion := obj.GetReuseVersion()
	billAutoNumber := obj.GetAutoNumber()
	if channelReuseVersion != uint32(paychan.ReuseVersion) {
		return fmt.Errorf("Payment Channel ReuseVersion is not match, need <%d> but got <%d>.",
			paychan.ReuseVersion, channelReuseVersion)
	}
	// 检查对账单资金数额和重用版本
	objlamt := obj.GetLeftBalance()
	objramt := obj.GetRightBalance()
	billTotalAmt, e21 := objlamt.Add(&objramt)
	if e21 != nil {
		return e21
	}
	paychanTotalAmt, e22 := paychan.LeftAmount.Add(&paychan.RightAmount)
	if e22 != nil {
		return e22
	}
	if billTotalAmt.NotEqual(paychanTotalAmt) {
		return fmt.Errorf("Payment Channel Total Amount is not match, need %s but got %s.",
			paychanTotalAmt.ToFinString(), billTotalAmt.ToFinString())
	}
	// 仲裁金额
	var assertTargetAmount = objlamt // 左侧
	var assertTargetSAT = obj.GetLeftSatoshi()
	if assertAddressIsRight {
		assertTargetAmount = objramt // 右侧
		assertTargetSAT = obj.GetRightSatoshi()
	}

	// 判断通道状态，是进入挑战期，还是最终夺取
	if paychan.IsOpening() {

		// 进入挑战期
		blkhei := state.GetPendingBlockHeight()
		// 改变状态
		paychan.SetChallenging(blkhei, assertAddressIsLeft,
			&assertTargetAmount, assertTargetSAT, uint64(billAutoNumber))
		// 写入状态
		return state.ChannelUpdate(channelId, paychan)

	} else if paychan.IsChallenging() {

		// 只能夺取对方，不能既自己提出仲裁，然后又自己回应挑战
		if paychan.AssertAddressIsLeftOrRight.Check() == assertAddressIsLeft {
			return fmt.Errorf("The arbitration request and the response cannot be the same address")
		}

		// 判断仲裁，是否夺取对方资金
		if billAutoNumber <= uint64(paychan.AssertBillAutoNumber) {
			// 账单流水号不满足（必须大于等待挑战的流水号）
			return fmt.Errorf("Payment Channel BillAutoNumber must more than %d.", paychan.AssertBillAutoNumber)
		}
		// 更高的流水号
		// 夺取全部资金，关闭通道
		var lamt = fields.NewEmptyAmount()
		var ramt = fields.NewEmptyAmount()
		var lsat = fields.Satoshi(0)
		var rsat = fields.Satoshi(0)
		ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
		if assertAddressIsLeft {
			lamt = paychanTotalAmt // 左侧账户夺取全部资金，包括HAC和SAT
			lsat = ttsat
		} else {
			ramt = paychanTotalAmt // 右侧账户夺取全部资金，包括HAC和SAT
			rsat = ttsat
		}
		// 关闭通道，夺取全部资金和利息
		isFinalClosed := true // 最终仲裁永久关闭
		return closePaymentChannelWriteinChainState(state, channelId, paychan, lamt, ramt, lsat, rsat, isFinalClosed)

	} else {
		return fmt.Errorf("Payment Channel <%s> status error.", hex.EncodeToString(channelId))
	}
}

// 挑战期或最终仲裁回退
func checkChannelGotoChallegingOrFinalDistributionRecoverChainState(state interfaces.ChainStateOperation, assertAddress fields.Address, paychan *stores.Channel, obj channel.OnChainChannelPaymentArbitrationReconciliationBasis) error {

	panic("RecoverChainState() func is deleted.")

	channelId := obj.GetChannelId()

	// 判断通道状态，是进入挑战期，还是最终夺取
	if paychan.IsFinalDistributionClosed() {
		if !paychan.IsHaveChallengeLog.Check() {
			return fmt.Errorf("IsHaveChallengeLog is not find.")
		}
		// 计算回退数额
		paychanTotalAmt, _ := paychan.LeftAmount.Add(&paychan.RightAmount)
		var lamt = fields.NewEmptyAmount()
		var ramt = fields.NewEmptyAmount()
		if paychan.LeftAddress.Equal(assertAddress) {
			lamt = paychanTotalAmt // 左侧账户夺取全部资金
		} else {
			ramt = paychanTotalAmt // 右侧账户夺取全部资金
		}
		// 状态回到挑战期
		isBackToChalleging := true
		// 回退账户余额
		return closePaymentChannelRecoverChainState_deprecated(state, channelId, lamt, ramt, isBackToChalleging)

	} else if paychan.IsChallenging() {

		// 回退到开启状态，清除挑战期数据
		paychan.SetOpening()
		paychan.CleanChallengingLog()
		return state.ChannelUpdate(channelId, paychan)

	} else {
		return fmt.Errorf("Payment Channel <%s> status error.", hex.EncodeToString(channelId))
	}
}

/////////////////////////////////////////////////////

// 挑战期结束，最终按主张分配通道资金
type Action_27_ClosePaymentChannelByClaimDistribution struct {
	// 通道 ID
	ChannelId fields.ChannelId // 通道id

	// data ptr
	belong_trs interfaces.Transaction
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Kind() uint16 {
	return 27
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Size() uint32 {
	return 2 + elm.ChannelId.Size()
}

// json api
func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Describe() map[string]interface{} {
	var data = map[string]interface{}{}
	return data
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Serialize() ([]byte, error) {
	var kindByte = make([]byte, 2)
	binary.BigEndian.PutUint16(kindByte, elm.Kind())
	var bt1, _ = elm.ChannelId.Serialize()
	var buffer bytes.Buffer
	buffer.Write(kindByte)
	buffer.Write(bt1)
	return buffer.Bytes(), nil
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error
	seek, e = elm.ChannelId.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	return seek, nil
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) RequestSignAddresses() []fields.Address {
	// 不需要任何人签名
	return []fields.Address{}
}

func (act *Action_27_ClosePaymentChannelByClaimDistribution) WriteinChainState(state interfaces.ChainStateOperation) error {

	if !sys.TestDebugLocalDevelopmentMark {
		return fmt.Errorf("mainnet not yet") // 暂未启用等待review
	}

	if act.belong_trs == nil {
		panic("Action belong to transaction not be nil !")
	}
	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 检查状态（必须为挑战期状态）
	if paychan.IsChallenging() == false {
		return fmt.Errorf("Payment Channel status is not on challenging.")
	}
	// 检查挑战期限
	clghei := state.GetPendingBlockHeight()
	expireHei := uint64(paychan.ChallengeLaunchHeight) + uint64(paychan.ArbitrationLockBlock)
	if clghei <= expireHei {
		// 挑战期还没过
		return fmt.Errorf("Payment Channel Challenging expire is %d.", expireHei)
	}
	// 按主张分配资金，结束通道
	var lamt = fields.NewEmptyAmount()
	var ramt = fields.NewEmptyAmount()
	var ttamt, e = paychan.LeftAmount.Add(&paychan.RightAmount)
	if e != nil {
		return e
	}
	var lsat = fields.Satoshi(0)
	var rsat = fields.Satoshi(0)
	ttsat := paychan.LeftSatoshi.GetRealSatoshi() + paychan.RightSatoshi.GetRealSatoshi()
	if paychan.AssertAddressIsLeftOrRight.Check() {
		lamt = &paychan.AssertAmount // 左侧为主张者
		ramt, _ = ttamt.Sub(lamt)
		lsat = paychan.AssertSatoshi.GetRealSatoshi()
		rsat = ttsat - lsat // 右侧自动分得剩余资金
	} else {
		ramt = &paychan.AssertAmount // 右侧为主张者
		lamt, _ = ttamt.Sub(ramt)
		rsat = paychan.AssertSatoshi.GetRealSatoshi()
		lsat = ttsat - rsat // 左侧自动分得剩余资金
	}

	// 永久关闭
	isFinnalClosed := true
	return closePaymentChannelWriteinChainState(state, act.ChannelId, paychan, lamt, ramt, lsat, rsat, isFinnalClosed)
}

func (act *Action_27_ClosePaymentChannelByClaimDistribution) RecoverChainState(state interfaces.ChainStateOperation) error {

	// 查询通道
	paychan := state.Channel(act.ChannelId)
	if paychan == nil {
		return fmt.Errorf("Payment Channel Id <%s> not find.", hex.EncodeToString(act.ChannelId))
	}
	// 回退状态
	// 按主张分配资金，结束通道
	var lamt = fields.NewEmptyAmount()
	var ramt = fields.NewEmptyAmount()
	if paychan.AssertAddressIsLeftOrRight.Check() {
		lamt = &paychan.AssertAmount // 左侧为主张者
	} else {
		ramt = &paychan.AssertAmount // 右侧为主张者
	}

	// 关闭
	isFinnalClosed := true
	return closePaymentChannelRecoverChainState_deprecated(state, act.ChannelId, lamt, ramt, isFinnalClosed)
}

func (elm *Action_27_ClosePaymentChannelByClaimDistribution) SetBelongTransaction(t interfaces.Transaction) {
	elm.belong_trs = t
}

// burning fees  // 是否销毁本笔交易的 90% 的交易费用
func (act *Action_27_ClosePaymentChannelByClaimDistribution) IsBurning90PersentTxFees() bool {
	return false
}
