package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	v3 "github.com/zhashkevych/dex-pools-aggregator/internal/dex/uniswap/v3"
)

/*
Ideas on development:
- Pause V3, implement V2
- Leave protobuf for now, but it won't be used
- Think on paralel dex processing & then calculation on arbitrage opportunities
*/

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

	uniswapV3Parser, err := v3.NewLiquidityPoolParser(client)
	if err != nil {
		log.Fatal(err)
	}

	_, err = uniswapV3Parser.CalculateEffectivePrice(v3.CalculateEffectivePriceInput{
		TokenInAddr:      v3.PoolsV3["USDC / ETH"].TokenOneAddr,
		TokenOutAddr:     v3.PoolsV3["USDC / ETH"].TokenTwoAddr,
		TokenInDecimals:  v3.PoolsV3["USDC / ETH"].TokenOne.GetMultiplicator(),
		TokenOutDecimals: v3.PoolsV3["USDC / ETH"].TokenTwo.GetMultiplicator(),
		AmountIn:         v3.PoolsV3["USDC / ETH"].TokenOneAmountIn,
		Fee:              v3.PoolsV3["USDC / ETH"].Fee,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- REVERSE ---")

	_, err = uniswapV3Parser.CalculateEffectivePrice(v3.CalculateEffectivePriceInput{
		TokenInAddr:      v3.PoolsV3["USDC / ETH"].TokenTwoAddr,
		TokenOutAddr:     v3.PoolsV3["USDC / ETH"].TokenOneAddr,
		TokenInDecimals:  v3.PoolsV3["USDC / ETH"].TokenTwo.GetMultiplicator(),
		TokenOutDecimals: v3.PoolsV3["USDC / ETH"].TokenOne.GetMultiplicator(),
		AmountIn:         v3.PoolsV3["USDC / ETH"].TokenTwoAmountIn,
		Fee:              v3.PoolsV3["USDC / ETH"].Fee,
	})
	if err != nil {
		log.Fatal(err)
	}
}
