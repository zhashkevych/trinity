/*
	THIS CODE SHOULD ALSO BE USED FOR SUSHISWAP INTERACTION
*/

package v2

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhashkevych/trinity/internal/dex"
	"github.com/zhashkevych/trinity/pkg/web3"
)

type LiquidityPoolParser struct {
	client *ethclient.Client
}

// This module can be used both for Uni V2 and Sushiswap
func NewLiquidityPoolParser(client *ethclient.Client) *LiquidityPoolParser {
	return &LiquidityPoolParser{client}
}

// TODO: maybe move to shared lib
type CalculateEffectivePriceInput struct {
	PoolName string
	PoolID   string
	// TokenInReserve   *big.Int
	// TokenOutReserve  *big.Int
	TokenInDecimals  int64
	TokenOutDecimals int64
	// AmountIn         *big.Int

	TradeAmount0 *big.Float
	TradeAmount1 *big.Float
}

// CalculateEffectivePrice requires PoolID, token decimals and AmountIn
func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) (*dex.EffectivePrice, error) {
	// todo: pass pool
	uniswapV2Pair, err := NewUniswapV2Pair(common.HexToAddress(inp.PoolID), lp.client)
	if err != nil {
		return nil, err
	}

	reserves, err := uniswapV2Pair.GetReserves(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	calcPriceInp := calculatePriceInput{
		TokenInDecimals:  inp.TokenInDecimals,
		TokenOutDecimals: inp.TokenOutDecimals,
		AmountIn:         inp.TradeAmount0,
	}

	calcPriceInp.TokenInReserve, calcPriceInp.TokenOutReserve = reserves.Reserve0, reserves.Reserve1

	effectivePrice0, err := lp.calculateEffectivePrice(calcPriceInp)
	if err != nil {
		return nil, err
	}

	calcPriceInp.TokenInReserve, calcPriceInp.TokenOutReserve = reserves.Reserve1, reserves.Reserve0
	calcPriceInp.TokenInDecimals, calcPriceInp.TokenOutDecimals = inp.TokenOutDecimals, inp.TokenInDecimals
	calcPriceInp.AmountIn = inp.TradeAmount1

	effectivePrice1, err := lp.calculateEffectivePrice(calcPriceInp)
	if err != nil {
		return nil, err
	}

	fmt.Println(inp.PoolName)
	fmt.Println("effective price 0:", effectivePrice0)
	fmt.Println("effective price 1:", effectivePrice1)

	return &dex.EffectivePrice{
		DexID:           dex.UNISWAP_V2,
		PoolID:          inp.PoolID,
		EffectivePrice0: effectivePrice0,
		EffectivePrice1: effectivePrice1,
		Timestamp:       time.Now(),
	}, nil
}

type calculatePriceInput struct {
	TokenInReserve   *big.Int
	TokenOutReserve  *big.Int
	TokenInDecimals  int64
	TokenOutDecimals int64
	AmountIn         *big.Float
}

func (lp LiquidityPoolParser) calculateEffectivePrice(inp calculatePriceInput) (*big.Float, error) {
	tokenInReserves := web3.FromTokenUnits(inp.TokenInReserve, inp.TokenInDecimals)
	tokenOutReserves := web3.FromTokenUnits(inp.TokenOutReserve, inp.TokenOutDecimals)

	netAmount := big.NewFloat(0).Mul(inp.AmountIn, big.NewFloat(1-0.003))

	tokenInReservesF := big.NewFloat(0).SetInt(tokenInReserves)
	tokenOutReservesF := big.NewFloat(0).SetInt(tokenOutReserves)

	newTokenABalance := big.NewFloat(0).Add(tokenInReservesF, netAmount)
	newTokenBBalance := big.NewFloat(0).Mul(tokenInReservesF, tokenOutReservesF)
	newTokenBBalance = newTokenBBalance.Quo(newTokenBBalance, newTokenABalance)

	tokenBReceived := tokenOutReservesF.Sub(tokenOutReservesF, newTokenBBalance)

	effectivePrice := inp.AmountIn.Quo(inp.AmountIn, tokenBReceived)
	if effectivePrice.IsInf() {
		return nil, errors.New("division by zero")
	}

	return effectivePrice, nil
}
