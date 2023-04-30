package v3

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zhashkevych/dex-pools-aggregator/pkg/web3"
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
	TokenInAddr      common.Address
	TokenOutAddr     common.Address
	TokenInDecimals  int64
	TokenOutDecimals int64
	AmountIn         *big.Int
	Fee              *big.Int
}

func (lp LiquidityPoolParser) CalculateEffectivePrice(inp CalculateEffectivePriceInput) (*big.Float, error) {
	amountIn := web3.ToTokenUnits(inp.AmountIn, inp.TokenInDecimals)
	amountInF := big.NewFloat(0).SetInt(amountIn)

	res := make([]interface{}, 0)

	err := lp.quoterv2.CallQuoteExactInputSingle(
		&bind.CallOpts{},
		&res,
		IQuoterV2QuoteExactInputSingleParams{
			TokenIn:           inp.TokenInAddr,
			TokenOut:          inp.TokenOutAddr,
			AmountIn:          amountIn,
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

	pricePerToken := big.NewFloat(0).Quo(amountInF, amountOut)

	return pricePerToken, nil
}

// EXTEND LOGIC OF QuoterV2 ABI
func (_QuoterV2 *QuoterV2Caller) CallQuoteExactInputSingle(opts *bind.CallOpts, res *[]interface{}, params IQuoterV2QuoteExactInputSingleParams) error {
	return _QuoterV2.contract.Call(opts, res, "quoteExactInputSingle", params)
}
