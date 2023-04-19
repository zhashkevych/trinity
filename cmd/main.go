package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/zhashkevych/dex-arbitrage/screener/internal/dex/uniswap"
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

	uniswap.Calculate(rpcURL)
}
