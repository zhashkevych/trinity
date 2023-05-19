package v3

import (
	"errors"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func NewLiquidityPoolParser(pool ClientPool) *LiquidityPoolParser {
	return &LiquidityPoolParser{pool}
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
	inp.amountIn = inp.TradeAmount0

	effectivePrice0, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v3/parser.go",
		}).Error("failed to calculate effective price 0", err)

		return nil, err
	}

	// reverse
	inp.TokenInAddr, inp.TokenOutAddr = inp.TokenOutAddr, inp.TokenInAddr
	inp.TokenInDecimals, inp.TokenOutDecimals = inp.TokenOutDecimals, inp.TokenInDecimals
	inp.amountIn = inp.TradeAmount1

	effectivePrice1, err := lp.calculateEffectivePrice(inp)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "uniswap/v3/parser.go",
		}).Error("failed to calculate effective price 1", err)

		return nil, err
	}

	return &dex.EffectivePrice{
		DexID:           dex.UNISWAP_V3,
		PoolID:          inp.PoolID,
		EffectivePrice0: effectivePrice0,
		EffectivePrice1: effectivePrice1,
		Timestamp:       time.Now(),
	}, nil
}

func (lp LiquidityPoolParser) calculateEffectivePrice(inp CalculateEffectivePriceInput) (*big.Float, error) {
	client, err := lp.clientPool.GetClient()
	if err != nil {
		return nil, err
	}

	quoterv2client, err := NewQuoterV2(common.HexToAddress(DexData[web3.ETHEREUM]["QuoterV2"]), client)
	if err != nil {
		return nil, err
	}

	amountIn := web3.ToTokenUnitsF(inp.amountIn, inp.TokenInDecimals)

	res, err := callQuoteExactInputSingleWithRetry(quoterv2client, inp)
	if err != nil {
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
