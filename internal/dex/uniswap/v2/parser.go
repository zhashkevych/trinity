package v2

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zhashkevych/dex-pools-aggregator/pkg/web3"
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
	TokenInAddr      common.Address
	TokenOutAddr     common.Address
	TokenInDecimals  int64
	TokenOutDecimals int64
	AmountIn         *big.Int
	Fee              *big.Int
}

func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) (*big.Int, error) {
	reserves, err := lp.uniswapV2Pair.GetReserves(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", inp)
	fmt.Printf("reserves %+v\n", reserves)

	// Some math right here
	tokenInReserves := FromTokenUnits(reserves.Reserve0, inp.TokenInDecimals)
	tokenOutReserves := FromTokenUnits(reserves.Reserve1, inp.TokenOutDecimals)

	netAmount := big.NewInt(0).Mul(inp.AmountIn, big.NewInt(1000-UniswapV2Fee))

	newTokenABalance := tokenInReserves.Add(tokenInReserves, netAmount)
	newTokenBBalance := big.NewInt(0).Mul(tokenInReserves, big.NewInt(0).Div(tokenOutReserves, newTokenABalance))

	tokenBReceived := tokenOutReserves.Sub(tokenOutReserves, newTokenBBalance)

	effectivePrice := inp.AmountIn.Div(inp.AmountIn, tokenBReceived)

	fmt.Println("effective price:", effectivePrice)

	return effectivePrice, nil
}

func FromTokenUnits(rawBalance *big.Int, decimals int64) *big.Int {
	return big.NewInt(0).Div(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}
