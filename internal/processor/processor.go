package processor

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhashkevych/trinity/internal/dex"
	v2 "github.com/zhashkevych/trinity/internal/dex/uniswap/v2"
	v3 "github.com/zhashkevych/trinity/internal/dex/uniswap/v3"
)

/*
	Processor should go through the list of all DEX pools and calculate effective price for each of them.
	When all pools are parsed, the data is aggregated to single array of all effective prices.
	Then it shoud be transported to the module, that searches for arbitrage opportunities.
*/

type UniV2Parser interface {
	CalculateEffectivePrice(inp v2.CalculateEffectivePriceInput) (*dex.EffectivePrice, error)
}

type UniV3Parser interface {
	CalculateEffectivePrice(inp v3.CalculateEffectivePriceInput) (*dex.EffectivePrice, error)
}

type DexPoolProcessor struct {
	uniV2Client UniV2Parser
	uniV3Client UniV3Parser
}

func NewDexPoolProcessor(uniV2Client UniV2Parser, uniV3Client UniV3Parser) *DexPoolProcessor {
	return &DexPoolProcessor{uniV2Client, uniV3Client}
}

func (p *DexPoolProcessor) StartProcessing(pools []*dex.PoolPair) {
	fmt.Println(len(pools))

	ts := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(len(pools))

	calculatedPoolPricesCh := make(chan *dex.PoolPair, len(pools))

	// calculate effective price for each pool
	calculatedPoolPrices := make([]*dex.PoolPair, len(pools))

	// Start a goroutine to receive data
	receiveWg := &sync.WaitGroup{}
	receiveWg.Add(1)
	go func() {
		defer receiveWg.Done()
		counter := 0
		for p := range calculatedPoolPricesCh {
			calculatedPoolPrices[counter] = p
			counter++
		}
	}()

	for _, pool := range pools {
		go p.calculateEffectivePrice(wg, calculatedPoolPricesCh, pool)
	}

	// Wait for all calculateEffectivePrice goroutines to finish
	wg.Wait()
	// Now we can safely close the channel
	close(calculatedPoolPricesCh)

	// Wait for the receiving goroutine to finish
	receiveWg.Wait()

	fmt.Println("time spent:", time.Since(ts))

	// Send to "Arbitrage Opportunity Finder" via Message Queue
	fmt.Println("sending effective prices to message queue")

	for i, p := range calculatedPoolPrices {
		fmt.Printf("%d. %+v\n", i+1, p)
	}
}

// todo handle errors
func (p *DexPoolProcessor) calculateEffectivePrice(wg *sync.WaitGroup, effectivePriceCh chan<- *dex.PoolPair, pool *dex.PoolPair) {
	defer wg.Done()

	tokenInDecimals, err := strconv.Atoi(pool.Token0.Decimals)
	if err != nil {
		// todo
	}

	tokenOutDecimals, err := strconv.Atoi(pool.Token1.Decimals)
	if err != nil {
		// todo
	}

	feeI, err := strconv.Atoi(pool.FeeTier)
	if err != nil {
		// todo
	}

	switch pool.DexID {
	case dex.UNISWAP_V2:
		// todo: pass trade amount
		effectivePrice, err := p.uniV2Client.CalculateEffectivePrice(v2.CalculateEffectivePriceInput{
			PoolID:           pool.ID,
			TokenInDecimals:  int64(tokenInDecimals),
			TokenOutDecimals: int64(tokenOutDecimals),
			TradeAmount0:     big.NewFloat(pool.TradeAmount0),
			TradeAmount1:     big.NewFloat(pool.TradeAmount1),
		})
		if err != nil {
			// todo
			fmt.Printf("UNI v2: error calculating effective price for %s: %s\n", pool.ID, err)
		}

		pool.EffectivePrice0 = effectivePrice.EffectivePrice0
		pool.EffectivePrice1 = effectivePrice.EffectivePrice1

		effectivePriceCh <- pool
	case dex.UNISWAP_V3:
		effectivePrice, err := p.uniV3Client.CalculateEffectivePrice(v3.CalculateEffectivePriceInput{
			PoolID:           pool.ID,
			TokenInAddr:      common.HexToAddress(pool.Token0.ID),
			TokenOutAddr:     common.HexToAddress(pool.Token1.ID),
			TokenInDecimals:  int64(tokenInDecimals),
			TokenOutDecimals: int64(tokenOutDecimals),
			TradeAmount0:     big.NewFloat(pool.TradeAmount0),
			TradeAmount1:     big.NewFloat(pool.TradeAmount1),
			Fee:              big.NewInt(int64(feeI)),
		})
		if err != nil {
			// todo
		}

		pool.EffectivePrice0 = effectivePrice.EffectivePrice0
		pool.EffectivePrice1 = effectivePrice.EffectivePrice1

		effectivePriceCh <- pool
	case dex.SUSHISWAP:
		// todo: pass trade amount
		effectivePrice, err := p.uniV2Client.CalculateEffectivePrice(v2.CalculateEffectivePriceInput{
			PoolID:           pool.ID,
			TokenInDecimals:  int64(tokenInDecimals),
			TokenOutDecimals: int64(tokenOutDecimals),
			TradeAmount0:     big.NewFloat(pool.TradeAmount0),
			TradeAmount1:     big.NewFloat(pool.TradeAmount1),
		})
		if err != nil {
			// todo
		}

		pool.EffectivePrice0 = effectivePrice.EffectivePrice0
		pool.EffectivePrice1 = effectivePrice.EffectivePrice1

		effectivePriceCh <- pool
	default:
		fmt.Println("Unknown DEX")
	}
}
