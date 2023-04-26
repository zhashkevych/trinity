package v3

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zhashkevych/dex-pools-aggregator/internal/dex"
	"github.com/zhashkevych/dex-pools-aggregator/internal/models"
	"github.com/zhashkevych/dex-pools-aggregator/pkg/erc20"
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

func (lp LiquidityPoolParser) CalculateEffectivePrice(poolPair *dex.PoolPair) (*big.Int, error) {
	// to token units
	fmt.Println("token 1 amount in", poolPair.TokenOneAmountIn)
	fmt.Println("token 1 multiplicator", poolPair.TokenOne.GetMultiplicator())

	amountIn := ToTokenUnits(poolPair.TokenOneAmountIn, poolPair.TokenOne.GetMultiplicator())

	fmt.Println(amountIn)
	fmt.Printf("%+v\n", poolPair)

	res := make([]interface{}, 0)

	err := lp.quoterv2.CallQuoteExactInputSingle(
		&bind.CallOpts{},
		&res,
		IQuoterV2QuoteExactInputSingleParams{
			TokenIn:           poolPair.TokenOneAddr,
			TokenOut:          poolPair.TokenTwoAddr,
			AmountIn:          amountIn,
			Fee:               poolPair.Fee,
			SqrtPriceLimitX96: big.NewInt(0),
		})
	if err != nil {
		fmt.Println("quoterv2 err: ", err)
		return nil, err
	}

	fmt.Println(res)

	if res == nil {
		return nil, errors.New("didn't receive data from Quoter V2")
	}

	amountOut := res[0].(*big.Int)
	fmt.Println("amount out", amountOut)

	amountOut.Mul(amountOut, big.NewInt(10^(poolPair.TokenOne.GetMultiplicator()-poolPair.TokenTwo.GetMultiplicator()))) // I don't know what's going on but it works. Just copied logic from Oleg's python code.

	pricePerToken := big.NewInt(0).Div(amountIn, amountOut)

	fmt.Println("amount in", amountIn)
	fmt.Println("amount out", amountOut)
	fmt.Println("price per token", pricePerToken)

	return nil, nil
}

func (lp LiquidityPoolParser) ParseAllEthereumPools() []*models.PoolData {
	out := make([]*models.PoolData, 0)

	for pair, pool := range PoolsV3 {
		fmt.Printf("-> Parsing %s <-\n", pair)

		poolData, err := lp.ParsePool(pool)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		out = append(out, poolData)
	}

	return out
}

func (lp LiquidityPoolParser) ParsePool(pool *dex.PoolPair) (*models.PoolData, error) {
	// Get token balances in the pool
	token0Contract, err := erc20.NewERC20(pool.TokenOneAddr, lp.client)
	if err != nil {
		fmt.Println("Error creating token0 contract instance")
		return nil, err
	}

	token0BalancePool, err := token0Contract.BalanceOf(&bind.CallOpts{}, pool.Addr)
	if err != nil {
		fmt.Println("Error getting token0 balance in the pool")
		return nil, err
	}

	token1Contract, err := erc20.NewERC20(pool.TokenTwoAddr, lp.client)
	if err != nil {
		fmt.Println("Error creating token1 contract instance")
		return nil, err
	}

	token1BalancePool, err := token1Contract.BalanceOf(&bind.CallOpts{}, pool.Addr)
	if err != nil {
		fmt.Println("Error getting token1 balance in the pool")
		return nil, err
	}

	return &models.PoolData{
		Dex:             "UNISWAP",
		DexId:           "0",
		Pool:            pool.Pair,
		PoolAddress:     pool.Addr.String(),
		TokenOne:        pool.TokenOneAddr.String(),
		TokenTwo:        pool.TokenTwoAddr.String(),
		TokenOneBalance: token0BalancePool.Text(10),
		TokenTwoBalance: token1BalancePool.Text(10),
		Fee:             pool.Fee.String(),
	}, nil
}

// FromTokenUnits converts a raw balance in the smallest token units to the full balance
func ToTokenUnits(rawBalance *big.Int, decimals int64) *big.Int {
	// return big.NewInt(0).Div(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	return big.NewInt(0).Mul(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}
