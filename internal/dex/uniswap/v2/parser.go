/*
	THIS CODE SHOULD ALSO BE USED FOR SUSHISWAP INTERACTION
*/

package v2

import (
	"errors"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhashkevych/trinity/internal/dex"
	"github.com/zhashkevych/trinity/pkg/web3"
)

type ClientPool interface {
	GetClient() (*ethclient.Client, error)
}

type LiquidityPoolParser struct {
	clientPool ClientPool
}

// This module can be used both for Uni V2 and Sushiswap
func NewLiquidityPoolParser(pool ClientPool) *LiquidityPoolParser {
	return &LiquidityPoolParser{pool}
}

// TODO: maybe move to shared lib
type CalculateEffectivePriceInput struct {
	PoolName         string
	PoolID           string
	TokenInDecimals  int64
	TokenOutDecimals int64

	TradeAmount0 *big.Float
	TradeAmount1 *big.Float
}

// CalculateEffectivePrice requires PoolID, token decimals and AmountIn
func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) (*dex.EffectivePrice, error) {
	// todo: pass pool
	client, err := lp.clientPool.GetClient()
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v2/parser.go",
		}).Error("failed to get client from client pool", err)

		return nil, err
	}

	uniswapV2Pair, err := NewUniswapV2Pair(common.HexToAddress(inp.PoolID), client)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v2/parser.go",
		}).Error("failed to init uniswapv2 pair client", err)

		return nil, err
	}

	reserves, err := getReservesWithRetry(uniswapV2Pair)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v2/parser.go",
		}).Error("failed to get reserves via uniswapv2 pair client", err)

		return nil, err
	}

	// log.WithFields(log.Fields{
	// 	"source": "uniswap/v2/parser.go",
	// }).Info("reserves", reserves)

	calcPriceInp := calculatePriceInput{
		TokenInDecimals:  inp.TokenInDecimals,
		TokenOutDecimals: inp.TokenOutDecimals,
		AmountIn:         inp.TradeAmount0,
	}

	calcPriceInp.TokenInReserve, calcPriceInp.TokenOutReserve = reserves.Reserve0, reserves.Reserve1

	effectivePrice0, err := lp.calculateEffectivePrice(calcPriceInp)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v2/parser.go",
		}).Error("failed to calculate effective price 0", err)

		return nil, err
	}

	calcPriceInp.TokenInReserve, calcPriceInp.TokenOutReserve = reserves.Reserve1, reserves.Reserve0
	calcPriceInp.TokenInDecimals, calcPriceInp.TokenOutDecimals = inp.TokenOutDecimals, inp.TokenInDecimals
	calcPriceInp.AmountIn = inp.TradeAmount1

	effectivePrice1, err := lp.calculateEffectivePrice(calcPriceInp)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v2/parser.go",
		}).Error("failed to calculate effective price 1", err)

		return nil, err
	}

	return &dex.EffectivePrice{
		DexID:           dex.UNISWAP_V2,
		PoolID:          inp.PoolID,
		Reserve0:        reserves.Reserve0,
		Reserve1:        reserves.Reserve1,
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

	effectivePrice := inp.AmountIn.Quo(tokenBReceived, inp.AmountIn)
	if effectivePrice.IsInf() {
		return nil, errors.New("division by zero")
	}

	return effectivePrice, nil
}
