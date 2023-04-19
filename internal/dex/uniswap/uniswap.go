package uniswap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zhashkevych/dex-arbitrage/screener/pkg/erc20"
)

type LiquidityPool struct {
	client *ethclient.Client
}

func NewLiquidityPool(client *ethclient.Client) *LiquidityPool {
	return &LiquidityPool{client}
}

func (lp LiquidityPool) Fetch() {
	// Get token addresses
	token0 := Tokens[ETHEREUM][USDC]
	token1 := Tokens[ETHEREUM][WETH]
	fee := big.NewInt(500)

	// Get Uniswap V3 pool address
	uniswapV3FactoryAddress := common.HexToAddress(DexData[ETHEREUM]["UniswapV3Factory"])

	uniswapV3Factory, err := NewUniswapV3Factory(uniswapV3FactoryAddress, lp.client)
	if err != nil {
		fmt.Println("Error parsing UniswapV3Factory ABI")
		return
	}

	uniswapV3PoolAddress, err := uniswapV3Factory.GetPool(&bind.CallOpts{}, token0, token1, fee)
	if err != nil {
		fmt.Println("Error getting UniswapV3 pool address")
		return
	}

	uniswapV3Pool, err := NewUniswapV3Pool(uniswapV3PoolAddress, lp.client)
	if err != nil {
		fmt.Println("Error getting UniswapV3 pool address")
		return
	}

	poolFee, err := uniswapV3Pool.Fee(&bind.CallOpts{})
	if err != nil {
		fmt.Println("Error getting UniswapV3 pool address")
		return
	}

	// Get token balances in the pool
	tokenAbiBytes, err := ioutil.ReadFile("abi/ERC20.json")
	if err != nil {
		fmt.Println("Error reading ERC20 ABI")
		return
	}
	var tokenAbi abi.ABI
	if err := json.Unmarshal(tokenAbiBytes, &tokenAbi); err != nil {
		fmt.Println("Error parsing ERC20 ABI")
		return
	}

	token0Address := Tokens[ETHEREUM][WETH]
	token0Contract, err := erc20.NewERC20(token0Address, lp.client)
	if err != nil {
		fmt.Println("Error creating token0 contract instance")
		return
	}

	token0BalancePool, err := token0Contract.BalanceOf(&bind.CallOpts{}, uniswapV3PoolAddress)
	if err != nil {
		fmt.Println("Error getting token0 balance in the pool")
		return
	}
	fmt.Printf("Token0 balance in the pool: %v\n", token0BalancePool)

	token1Address := Tokens["ETHEREUM"]["USDC"]
	token1Contract, err := erc20.NewERC20(token1Address, lp.client)
	if err != nil {
		fmt.Println("Error creating token1 contract instance")
		return
	}

	token1BalancePool, err := token1Contract.BalanceOf(&bind.CallOpts{}, uniswapV3PoolAddress)
	if err != nil {
		fmt.Println("Error getting token1 balance in the pool")
		return
	}
	fmt.Printf("Token1 balance in the pool: %v\n", token1BalancePool)

	fmt.Printf("Fee: %v\n", poolFee)
}
