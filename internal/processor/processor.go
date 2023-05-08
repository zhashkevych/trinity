package processor

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/zhashkevych/trinity/internal/dex"
	v2 "github.com/zhashkevych/trinity/internal/dex/uniswap/v2"
	v3 "github.com/zhashkevych/trinity/internal/dex/uniswap/v3"
)

/*
	Processor should go through the list of all DEX pools and calculate effective price for each of them.
	When all pools are parsed, the data is aggregated to single array of all effective prices.
	Then it shoud be transported to the module, that searches for arbitrage opportunities.
*/

type DexPoolProcessor struct {
	uniV2Client v2.LiquidityPoolParser
	uniV3Client v3.LiquidityPoolParser
}

func NewDexPoolProcessor() *DexPoolProcessor {
	return &DexPoolProcessor{}
}

func (p *DexPoolProcessor) StartProcessing(pools []*dex.PoolPair) {
	wg := &sync.WaitGroup{}
	effectivePriceCh := make(chan *dex.EffectivePrice)

	for _, pool := range pools {
		wg.Add(1)
		go p.calculateEffectivePrice(wg, effectivePriceCh, pool)
	}
}

func (p *DexPoolProcessor) calculateEffectivePrice(wg *sync.WaitGroup, effectivePriceCh chan<- *dex.EffectivePrice, pool *dex.PoolPair) {
	defer wg.Done()

	tokenInDecimals, err := strconv.Atoi(pool.Token0.Decimals)
	if err != nil {
		// todo
	}

	tokenOutDecimals, err := strconv.Atoi(pool.Token1.Decimals)
	if err != nil {
		// todo
	}

	switch pool.DexID {
	case dex.UNISWAP_V2:
		effectivePrice, err := p.uniV2Client.CalculateEffectivePrice(v2.CalculateEffectivePriceInput{
			PoolID:           pool.ID,
			TokenInDecimals:  int64(tokenInDecimals),
			TokenOutDecimals: int64(tokenOutDecimals),
		})
		if err != nil {
			// todo
		}

		effectivePriceCh <- effectivePrice
	case dex.UNISWAP_V3:
		p.uniV3Client.CalculateEffectivePrice()
	case dex.SUSHISWAP:
		// TODO
	default:
		fmt.Println("Unknown DEX")
	}
}
