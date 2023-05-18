package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/zhashkevych/trinity/internal/clientpool"
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

	// Init client pool
	urls, err := getAlchemyURLs()
	if err != nil {
		fmt.Println("Error parsing Alchemy URLs")
		return
	}

	clients := make([]*ethclient.Client, len(urls))

	for i, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			fmt.Println("Error connecting to Ethereum")
			return
		}
		defer client.Close()

		clients[i] = client
	}

	clientPool := clientpool.NewPool(clients)

	for i := 0; i < 100; i++ {
		clientPool.GetClient()
	}

	// Connect to NATS MQ
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

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

	// fmt.Println(len(uniswapV2Pools))
	// fmt.Println(len(uniswapV3Pools))

	// for _, pool := range uniswapV3Pools {
	// 	fmt.Printf("%+v\n", pool)
	// }

	uniswapV3Parser := v3.NewLiquidityPoolParser(clientPool)
	uniswapV2Parser := v2.NewLiquidityPoolParser(clientPool)

	pools := make([]*dex.PoolPair, 0)

	pools = append(pools, uniswapV2Pools[0:50]...)
	pools = append(pools, uniswapV3Pools[0:50]...)

	p := processor.NewDexPoolProcessor(uniswapV2Parser, uniswapV3Parser, nc)
	p.StartProcessing(pools)

	// Make sure the data is sent before we close the connection
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE!")
}

func getAlchemyURLs() ([]string, error) {
	urlsCount, err := strconv.Atoi(os.Getenv("URLS_COUNT"))
	if err != nil {
		return nil, err
	}

	urls := make([]string, urlsCount)
	for i := 0; i < urlsCount; i++ {
		urls[i] = os.Getenv(fmt.Sprintf("ALCHEMY_RPC_URL_MAINNET_%d", i))
	}

	return urls, nil
}
