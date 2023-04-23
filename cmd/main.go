package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	v3 "github.com/zhashkevych/dex-pools-aggregator/internal/dex/uniswap/v3"
)

type Token struct {
	Name    string
	Address common.Address
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Connect to Ethereum
	rpcURL := os.Getenv("ALCHEMY_RPC_URL_MAINNET")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		fmt.Println("Error connecting to Ethereum")
		return
	}
	defer client.Close()

	uniswapLPClient := v3.NewLiquidityPoolClient(client)
	poolsData := uniswapLPClient.ParseAllEthereumPools()

	for _, pd := range poolsData {
		fmt.Printf("%+v\n", pd)
	}
}
