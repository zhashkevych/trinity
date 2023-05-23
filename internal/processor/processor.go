package processor

import (
	"math/big"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	"github.com/zhashkevych/trinity/internal/dex"
	v2 "github.com/zhashkevych/trinity/internal/dex/uniswap/v2"
	v3 "github.com/zhashkevych/trinity/internal/dex/uniswap/v3"
	"github.com/zhashkevych/trinity/internal/models"
	"google.golang.org/protobuf/proto"
)

const (
	NATS_SUBJECT = "calculated-prices"
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

	mq *nats.Conn // move to separate transport layer
}

func NewDexPoolProcessor(uniV2Client UniV2Parser, uniV3Client UniV3Parser, mq *nats.Conn) *DexPoolProcessor {
	return &DexPoolProcessor{uniV2Client, uniV3Client, mq}
}

func (p *DexPoolProcessor) StartProcessing(pools []*dex.PoolPair) {
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

	log.WithFields(log.Fields{
		"source": "processor.go",
	}).Info("time spent:", time.Since(ts))

	// Send to "Arbitrage Opportunity Finder" via Message Queue
	// 1. Convert to Protobuf
	poolPairs := toProtoModel(calculatedPoolPrices)
	data, err := proto.Marshal(poolPairs)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "processor.go",
		}).Error("failed to marshal proto data:", err)

		return
	}

	// 2. Send to NATS
	if err := p.mq.Publish(NATS_SUBJECT, data); err != nil {
		log.WithFields(log.Fields{
			"source": "processor.go",
		}).Error("failed to sent proto data to NATS:", err)

		return
	}

	log.WithFields(log.Fields{
		"source": "processor.go",
	}).Info("sent effective prices to NATS")
}

// todo handle errors
func (p *DexPoolProcessor) calculateEffectivePrice(wg *sync.WaitGroup, effectivePriceCh chan<- *dex.PoolPair, pool *dex.PoolPair) {
	defer wg.Done()

	tokenInDecimals, err := strconv.Atoi(pool.Token0.Decimals)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "processor.go",
			"method": "calculateEffectivePrice",
			"poolID": pool.ID,
			"dex":    pool.DexID,
		}).Error("failed to convert token0.Decimals to int:", err)

		return
	}

	tokenOutDecimals, err := strconv.Atoi(pool.Token1.Decimals)
	if err != nil {
		log.WithFields(log.Fields{
			"source": "processor.go",
			"method": "calculateEffectivePrice",
			"poolID": pool.ID,
			"dex":    pool.DexID,
		}).Error("failed to convert token1.Decimals to int:", err)

		return
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
			log.WithFields(log.Fields{
				"source": "processor.go",
				"method": "calculateEffectivePrice",
				"poolID": pool.ID,
				"dex":    pool.DexID,
			}).Error("failed to calculate effective price:", err)
		}

		if effectivePrice != nil {
			pool.EffectivePrice0 = effectivePrice.EffectivePrice0
			pool.EffectivePrice1 = effectivePrice.EffectivePrice1
		}

		effectivePriceCh <- pool
	case dex.UNISWAP_V3:
		feeI, err := strconv.Atoi(pool.FeeTier)
		if err != nil {
			log.WithFields(log.Fields{
				"source": "processor.go",
				"method": "calculateEffectivePrice",
				"poolID": pool.ID,
				"dex":    pool.DexID,
			}).Error("failed to convert feeTier to int:", err)

			return
		}

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
			log.WithFields(log.Fields{
				"source": "processor.go",
				"method": "calculateEffectivePrice",
				"poolID": pool.ID,
				"dex":    pool.DexID,
			}).Error("failed to calculate effective price:", err)
		}

		if effectivePrice != nil {
			pool.EffectivePrice0 = effectivePrice.EffectivePrice0
			pool.EffectivePrice1 = effectivePrice.EffectivePrice1
		}

		effectivePriceCh <- pool
	case dex.SUSHISWAP:
		effectivePrice, err := p.uniV2Client.CalculateEffectivePrice(v2.CalculateEffectivePriceInput{
			PoolID:           pool.ID,
			TokenInDecimals:  int64(tokenInDecimals),
			TokenOutDecimals: int64(tokenOutDecimals),
			TradeAmount0:     big.NewFloat(pool.TradeAmount0),
			TradeAmount1:     big.NewFloat(pool.TradeAmount1),
		})
		if err != nil {
			log.WithFields(log.Fields{
				"source": "processor.go",
				"method": "calculateEffectivePrice",
				"poolID": pool.ID,
				"dex":    pool.DexID,
			}).Error("failed to calculate effective price:", err)
		}

		if effectivePrice != nil {
			pool.EffectivePrice0 = effectivePrice.EffectivePrice0
			pool.EffectivePrice1 = effectivePrice.EffectivePrice1
		}

		effectivePriceCh <- pool
	default:
	}
}

func toProtoModel(poolPairs []*dex.PoolPair) *models.PoolPairList {
	pairs := make([]*models.PoolPair, 0)

	for _, pair := range poolPairs {
		if pair == nil {
			continue
		}
		pairs = append(pairs, toProtoPoolPair(pair))
	}

	return &models.PoolPairList{
		PoolPairs: pairs,
	}
}

func toProtoPoolPair(p *dex.PoolPair) *models.PoolPair {
	out := &models.PoolPair{
		DexId: p.DexID.GetProto(),
		Id:    p.ID,
		Token0: &models.Token{
			Id:         p.Token0.ID,
			Symbol:     p.Token0.Symbol,
			DerivedEth: p.Token0.DerivedETH,
			Decimals:   p.Token0.Decimals,
		},
		Token1: &models.Token{
			Id:         p.Token1.ID,
			Symbol:     p.Token1.Symbol,
			DerivedEth: p.Token1.DerivedETH,
			Decimals:   p.Token1.Decimals,
		},
		Reserve0:     p.Reserve0,
		Reserve1:     p.Reserve1,
		Reserve0Usd:  p.Reserve0USD,
		Reserve1Usd:  p.Reserve1USD,
		TradeAmount0: p.TradeAmount0,
		TradeAmount1: p.TradeAmount1,
	}

	if p.EffectivePrice0 != nil {
		out.EffectivePrice0 = p.EffectivePrice0.String()
	}

	if p.EffectivePrice1 != nil {
		out.EffectivePrice1 = p.EffectivePrice1.String()
	}

	return out
}
