package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/zhashkevych/trinity/internal/dex"
	v2 "github.com/zhashkevych/trinity/internal/dex/uniswap/v2"
	v3 "github.com/zhashkevych/trinity/internal/dex/uniswap/v3"
	"github.com/zhashkevych/trinity/internal/processor"
)

/*
TODO:
- Parse pools from JSON
- Implement SushiSwap
- Implement async processing
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

	file, err := ioutil.ReadFile("pools/uniswapV2_Pools.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	uniswapV2Pools := make([]*dex.PoolPair, 0)
	err = json.Unmarshal(file, &uniswapV2Pools)
	if err != nil {
		fmt.Println("Error unmarshaling V2 JSON data:", err)
		return
	}

	for i := range uniswapV2Pools {
		uniswapV2Pools[i].DexID = dex.UNISWAP_V2
	}

	file, err = ioutil.ReadFile("pools/uniswapV3_Pools.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	uniswapV3Pools := make([]*dex.PoolPair, 0)
	err = json.Unmarshal(file, &uniswapV3Pools)
	if err != nil {
		fmt.Println("Error unmarshaling V3 JSON data:", err)
		return
	}

	for i := range uniswapV3Pools {
		uniswapV3Pools[i].DexID = dex.UNISWAP_V3
	}

	fmt.Println(len(uniswapV2Pools))
	fmt.Println(len(uniswapV3Pools))

	// for _, pool := range uniswapV3Pools {
	// 	fmt.Printf("%+v\n", pool)
	// }

	uniswapV3Parser, err := v3.NewLiquidityPoolParser(client)
	if err != nil {
		log.Fatal(err)
	}

	// effectivePrice1, err := uniswapV3Parser.CalculateEffectivePrice(v3.CalculateEffectivePriceInput{
	// 	TokenInAddr:      uniswapV3Pools[0].Token0.ID,
	// 	TokenOutAddr:     uniswapV3Pools[0].Token1.ID,
	// 	TokenInDecimals:  uniswapV3Pools[0].Token0.Decimals,
	// 	TokenOutDecimals: uniswapV3Pools[0].Token0.Decimals,
	// 	// AmountIn:         v3.PoolsV3["USDC / ETH"].TokenOneAmountIn,
	// 	// Fee:             uniswapV3Pools[0],
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// effectivePrice2, err := uniswapV3Parser.CalculateEffectivePrice(v3.CalculateEffectivePriceInput{
	// 	TokenInAddr:      v3.PoolsV3["USDC / ETH"].TokenTwoAddr,
	// 	TokenOutAddr:     v3.PoolsV3["USDC / ETH"].TokenOneAddr,
	// 	TokenInDecimals:  v3.PoolsV3["USDC / ETH"].TokenTwo.GetMultiplicator(),
	// 	TokenOutDecimals: v3.PoolsV3["USDC / ETH"].TokenOne.GetMultiplicator(),
	// 	AmountIn:         v3.PoolsV3["USDC / ETH"].TokenTwoAmountIn,
	// 	Fee:              v3.PoolsV3["USDC / ETH"].Fee,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("USDC / ETH")
	// fmt.Println("effective price 1:", effectivePrice1)
	// fmt.Println("effective price 2:", effectivePrice2)

	// fmt.Println("--- UNISWAP v2 ---")

	uniswapV2Parser := v2.NewLiquidityPoolParser(client)
	if err != nil {
		log.Fatal(err)
	}

	pools := make([]*dex.PoolPair, 0)

	pools = append(pools, uniswapV2Pools...)
	pools = append(pools, uniswapV3Pools...)

	p := processor.NewDexPoolProcessor(uniswapV2Parser, uniswapV3Parser)
	p.StartProcessing(pools)

	// uniswapV2Parser.CalculateEffectivePrice(v2.CalculateEffectivePriceInput{
	// 	PoolName:         "USDC / ETH",
	// 	TokenInDecimals:  uniswapV2Pools[0].Token0.Decimals,
	// 	TokenOutDecimals: uniswapV2Pools[0].Token1.Decimals,
	// 	AmountIn:         v2.PoolsV2["USDC / ETH"].TokenOneAmountIn,
	// })
}
