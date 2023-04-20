package uniswap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zhashkevych/dex-pools-aggregator/internal/models"
	"github.com/zhashkevych/dex-pools-aggregator/pkg/erc20"
)

// TODO: in which datatype should we exchange balances?

type LiquidityPoolClient struct {
	client *ethclient.Client
}

func NewLiquidityPoolClient(client *ethclient.Client) *LiquidityPoolClient {
	return &LiquidityPoolClient{client}
}

func (lp LiquidityPoolClient) ParseAllEthereumPools() []*models.PoolData {
	out := make([]*models.PoolData, 0)

	for pair, pool := range Pools {
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

func (lp LiquidityPoolClient) ParsePool(pool Pool) (*models.PoolData, error) {
	// Get token balances in the pool
	tokenAbiBytes, err := ioutil.ReadFile("abi/ERC20.json")
	if err != nil {
		fmt.Println("Error reading ERC20 ABI")
		return nil, err
	}
	var tokenAbi abi.ABI
	if err := json.Unmarshal(tokenAbiBytes, &tokenAbi); err != nil {
		fmt.Println("Error parsing ERC20 ABI")
		return nil, err
	}

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
