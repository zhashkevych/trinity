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
	"github.com/zhashkevych/dex-arbitrage/screener/internal/dex/data"
	"github.com/zhashkevych/dex-arbitrage/screener/pkg/erc20"
)

func Calculate(rpcURL string) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		fmt.Println("Error connecting to Ethereum")
		return
	}
	defer client.Close()

	// Get token addresses
	token0 := data.Tokens["ETHEREUM"]["USDC"]
	token1 := data.Tokens["ETHEREUM"]["WETH"]
	fee := big.NewInt(500)

	fmt.Println(token0.Hex())
	fmt.Println(token1.Hex())

	// Get Uniswap V3 pool address
	uniswapV3FactoryAddress := common.HexToAddress(data.DexData["ETHEREUM"]["UniswapV3Factory"])

	uniswapV3Factory, err := NewUniswapV3Factory(uniswapV3FactoryAddress, client)
	if err != nil {
		fmt.Println("Error parsing UniswapV3Factory ABI")
		return
	}

	uniswapV3PoolAddress, err := uniswapV3Factory.GetPool(&bind.CallOpts{}, token0, token1, fee)
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

	token0Address := data.Tokens["ETHEREUM"]["WETH"]
	token0Contract, err := erc20.NewERC20(token0Address, client)
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

	token1Address := data.Tokens["ETHEREUM"]["USDC"]
	token1Contract, err := erc20.NewERC20(token1Address, client)
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
}
