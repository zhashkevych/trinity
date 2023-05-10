package v3

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
	client   *ethclient.Client
	quoterv2 *QuoterV2
}

func NewLiquidityPoolParser(client *ethclient.Client) (*LiquidityPoolParser, error) {
	quoterv2client, err := NewQuoterV2(common.HexToAddress(DexData[web3.ETHEREUM]["QuoterV2"]), client)
	if err != nil {
		return nil, err
	}

	return &LiquidityPoolParser{client, quoterv2client}, nil
}

type CalculateEffectivePriceInput struct {
	PoolID           string
	TokenInAddr      common.Address
	TokenOutAddr     common.Address
	TokenInDecimals  int64
	TokenOutDecimals int64
	TradeAmount0     *big.Float
	TradeAmount1     *big.Float
	Fee              *big.Int

	amountIn *big.Float
}

func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) (*dex.EffectivePrice, error) {
	// tradeAmount0, _ := inp.TradeAmount0.Int64()
	// inp.amountIn = big.NewInt(tradeAmount0)
	inp.amountIn = inp.TradeAmount0

	effectivePrice0, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		return nil, err
	}

	// reverse
	inp.TokenInAddr, inp.TokenOutAddr = inp.TokenOutAddr, inp.TokenInAddr
	inp.TokenInDecimals, inp.TokenOutDecimals = inp.TokenOutDecimals, inp.TokenInDecimals
	// tradeAmount1, _ := inp.TradeAmount1.Int64()
	// inp.amountIn = big.NewInt(tradeAmount1)
	inp.amountIn = inp.TradeAmount1

	effectivePrice1, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		return nil, err
	}

	// fmt.Println("V3, PoolID ", inp.PoolID)
	// fmt.Println("effective price 0:", effectivePrice0)
	// fmt.Println("effective price 1:", effectivePrice1)

	return &dex.EffectivePrice{
		DexID:           dex.UNISWAP_V3,
		PoolID:          inp.PoolID,
		EffectivePrice0: effectivePrice0,
		EffectivePrice1: effectivePrice1,
		Timestamp:       time.Now(),
	}, nil
}

func (lp LiquidityPoolParser) calculateEffectivePrice(inp CalculateEffectivePriceInput) (*big.Float, error) {
	amountIn := web3.ToTokenUnitsF(inp.amountIn, inp.TokenInDecimals)
	// amountInF := big.NewFloat(0).SetInt(amountIn)

	res := make([]interface{}, 0)
	amountInI, _ := amountIn.Int(nil)

	err := lp.quoterv2.CallQuoteExactInputSingle(
		&bind.CallOpts{},
		&res,
		IQuoterV2QuoteExactInputSingleParams{
			TokenIn:           inp.TokenInAddr,
			TokenOut:          inp.TokenOutAddr,
			AmountIn:          amountInI,
			Fee:               inp.Fee,
			SqrtPriceLimitX96: big.NewInt(0),
		})
	if err != nil {
		fmt.Println("quoterv2 err: ", err)
		return nil, err
	}

	if res == nil {
		return nil, errors.New("didn't receive data from Quoter V2")
	}

	amountOut := big.NewFloat(0).SetInt(res[0].(*big.Int))

	// This logic in python looks like that
	// amount_out = amount_out * (10 ** (token_in_decimals - token_out_decimals))

	var pow big.Int                                                        // Here we calculate the value for 10 ^ (x). x can be negative.
	pow.Abs(big.NewInt(int64(inp.TokenInDecimals - inp.TokenOutDecimals))) // USDC (6) - ETH (18) = -12
	exp := big.NewInt(10).Exp(big.NewInt(10), &pow, nil)

	if inp.TokenInDecimals < inp.TokenOutDecimals {
		amountOut = amountOut.Quo(amountOut, big.NewFloat(0).SetInt(exp))
	} else {
		amountOut = amountOut.Mul(amountOut, big.NewFloat(0).SetInt(exp))
	}

	pricePerToken := big.NewFloat(0).Quo(amountIn, amountOut)

	return pricePerToken, nil
}

// EXTEND LOGIC OF QuoterV2 ABI
func (_QuoterV2 *QuoterV2Caller) CallQuoteExactInputSingle(opts *bind.CallOpts, res *[]interface{}, params IQuoterV2QuoteExactInputSingleParams) error {
	return _QuoterV2.contract.Call(opts, res, "quoteExactInputSingle", params)
}
