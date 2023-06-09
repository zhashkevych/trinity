package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

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
- Implement proper logging
- Implement graceful shutdown
*/

const PROCESSING_INTERVAL = time.Second * 10

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to load env file", err)
		return
	}

	// Init client pool
	urls, err := getAlchemyURLs()
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to load alchemy urls", err)
		return
	}

	clients := make([]*ethclient.Client, len(urls))

	for i, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			log.WithFields(log.Fields{
				"source": "main.go",
			}).Error("failed to init eth client with url", url)
			return
		}
		defer client.Close()

		clients[i] = client
	}

	clientPool := clientpool.NewPool(clients)

	for i := 0; i < 100; i++ {
		clientPool.GetClient()
	}

	natsURL := os.Getenv("NATS_URL")

	// Connect to NATS MQ
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to connect to nats", err)
		return
	}
	defer nc.Close()

	file, err := ioutil.ReadFile("pools/uniswapV2_Pools.json")
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to read uniswap v2 json", err)
		return
	}

	uniswapV2Pools := make([]*dex.PoolPair, 0)
	err = json.Unmarshal(file, &uniswapV2Pools)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to unmarshal uniswap v2 json", err)
		return
	}

	for i := range uniswapV2Pools {
		uniswapV2Pools[i].DexID = dex.UNISWAP_V2
	}

	file, err = ioutil.ReadFile("pools/uniswapV3_Pools.json")
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to read uniswap v3 json", err)
		return
	}

	uniswapV3Pools := make([]*dex.PoolPair, 0)
	err = json.Unmarshal(file, &uniswapV3Pools)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to unmarshal uniswap v3 json", err)
		return
	}

	for i := range uniswapV3Pools {
		uniswapV3Pools[i].DexID = dex.UNISWAP_V3
	}

	file, err = ioutil.ReadFile("pools/sushiswap_Pools.json")
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to read uniswap v3 json", err)
		return
	}

	sushiswapPools := make([]*dex.PoolPair, 0)
	err = json.Unmarshal(file, &sushiswapPools)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("failed to unmarshal uniswap v3 json", err)
		return
	}

	for i := range sushiswapPools {
		sushiswapPools[i].DexID = dex.SUSHISWAP
	}

	uniswapV3Parser := v3.NewLiquidityPoolParser(clientPool)
	uniswapV2Parser := v2.NewLiquidityPoolParser(clientPool)

	pools := make([]*dex.PoolPair, 0)

	pools = append(pools, uniswapV2Pools[0:200]...)
	pools = append(pools, uniswapV3Pools[0:200]...)
	pools = append(pools, sushiswapPools[0:200]...)

	p := processor.NewDexPoolProcessor(uniswapV2Parser, uniswapV3Parser, nc)

	// Create a channel to listen for termination signals
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	go func(doneCh chan struct{}) {
		timer := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-timer.C:
				log.WithFields(log.Fields{
					"source": "main.go",
					"time":   time.Now(),
				}).Info("Started processing")
				p.StartProcessing(pools)
			case <-doneCh:
				log.WithFields(log.Fields{
					"source": "main.go",
					"time":   time.Now(),
				}).Info("Finished processing")
			}
		}
	}(done)

	// Wait for a termination signal
	<-termChan
	done <- struct{}{}

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.WithFields(log.Fields{
			"source": "main.go",
		}).Error("NATS error", err)
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
