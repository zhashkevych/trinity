package dex

import (
	"math/big"

	"github.com/zhashkevych/trinity/internal/models"
)

// TODO: store smart contract addresses in JSON config and parse into global config

type DEX string

const (
	UNISWAP_V2 DEX = "uniswap v2"
	UNISWAP_V3 DEX = "uniswap v3"
	SUSHISWAP  DEX = "sushi swap"
)

func (d DEX) GetProto() models.DEX {
	switch d {
	case UNISWAP_V2:
		return models.DEX_UNISWAP_V2
	case UNISWAP_V3:
		return models.DEX_UNISWAP_V3
	case SUSHISWAP:
		return models.DEX_SUSHISWAP
	default:
		return -1
	}
}

type Token struct {
	Symbol     string `json:"symbol"`
	ID         string `json:"id"`
	DerivedETH string `json:"derivedETH"`
	Decimals   string `json:"decimals"`
}

type PoolPair struct {
	DexID               DEX
	ID                  string  `json:"id"`
	Token0              Token   `json:"token0"`
	Token1              Token   `json:"token1"`
	Liquidity           string  `json:"liquidity"`
	SqrtPrice           string  `json:"sqrtPrice"`
	TotalValueLockedUSD string  `json:"totalValueLockedUSD"`
	Reserve0            float64 `json:"reserve0"`
	Reserve1            float64 `json:"reserve1"`
	Reserve0USD         float64 `json:"reserve0USD"`
	Reserve1USD         float64 `json:"reserve1USD"`
	TradeAmount0        float64 `json:"tradeAmount0"`
	TradeAmount1        float64 `json:"tradeAmount1"`
	FeeTier             string  `json:"feeTier"`

	// Calculated inside the app
	EffectivePrice0 *big.Float
	EffectivePrice1 *big.Float
}

type Reserves struct {
	Reserve0 string `json:"reserve0"`
	Reserve1 string `json:"reserve1"`
}
