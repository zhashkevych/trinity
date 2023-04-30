package v2

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhashkevych/trinity/pkg/web3"
)

type LiquidityPoolParser struct {
	client        *ethclient.Client
	uniswapV2Pair *UniswapV2Pair
}

func NewLiquidityPoolParser(client *ethclient.Client) (*LiquidityPoolParser, error) {
	uniswapV2Pair, err := NewUniswapV2Pair(common.HexToAddress(DexData[web3.ETHEREUM]["UniswapV2Pair"]), client)
	if err != nil {
		return nil, err
	}

	return &LiquidityPoolParser{client, uniswapV2Pair}, nil
}

// TODO: maybe move to shared lib
type CalculateEffectivePriceInput struct {
	PoolName         string
	TokenInReserve   *big.Int
	TokenOutReserve  *big.Int
	TokenInDecimals  int64
	TokenOutDecimals int64
	AmountIn         *big.Int
}

func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) {
	reserves, err := lp.uniswapV2Pair.GetReserves(&bind.CallOpts{})
	if err != nil {
		return
	}

	inp.TokenInReserve, inp.TokenOutReserve = reserves.Reserve0, reserves.Reserve1

	effectivePrice1, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		return
	}

	inp.TokenInReserve, inp.TokenOutReserve = reserves.Reserve1, reserves.Reserve0
	inp.TokenInDecimals, inp.TokenOutDecimals = inp.TokenOutDecimals, inp.TokenInDecimals

	effectivePrice2, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		return
	}

	fmt.Println(inp.PoolName)
	fmt.Println("effective price 1:", effectivePrice1)
	fmt.Println("effective price 2:", effectivePrice2)
}

func (lp LiquidityPoolParser) calculateEffectivePrice(inp CalculateEffectivePriceInput) (*big.Float, error) {
	tokenInReserves := web3.FromTokenUnits(inp.TokenInReserve, inp.TokenInDecimals)
	tokenOutReserves := web3.FromTokenUnits(inp.TokenOutReserve, inp.TokenOutDecimals)

	amountInF := big.NewFloat(0).SetInt(inp.AmountIn)
	netAmount := big.NewFloat(0).Mul(amountInF, big.NewFloat(1-0.003))

	tokenInReservesF := big.NewFloat(0).SetInt(tokenInReserves)
	tokenOutReservesF := big.NewFloat(0).SetInt(tokenOutReserves)

	newTokenABalanceF := big.NewFloat(0).Add(tokenInReservesF, netAmount)
	newTokenBBalanceF := big.NewFloat(0).Mul(tokenInReservesF, tokenOutReservesF)
	newTokenBBalanceF = newTokenBBalanceF.Quo(newTokenBBalanceF, newTokenABalanceF)

	tokenBReceivedF := tokenOutReservesF.Sub(tokenOutReservesF, newTokenBBalanceF)

	effectivePriceF := amountInF.Quo(amountInF, tokenBReceivedF)
	if effectivePriceF.IsInf() {
		return nil, errors.New("division by zero")
	}

	return effectivePriceF, nil
}
