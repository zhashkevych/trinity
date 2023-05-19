package v2

import (
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const maxRetries = 5
const baseDelay = 1 * time.Second

type reserves struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}

func getReservesWithRetry(uniswapV2Pair *UniswapV2Pair) (reserves, error) {
	var err error
	var reserves reserves

	for i := 0; i < maxRetries; i++ {
		reserves, err = uniswapV2Pair.GetReserves(&bind.CallOpts{})
		if err == nil {
			// USED FOR DEBUGING
			// if i > 0 {
			// 	log.WithFields(log.Fields{
			// 		"source": "uniswap/v2/parser.go",
			// 		"retry":  i,
			// 	}).Info("request successfully sent")
			// }
			break
		}

		// If the error is related to rate limiting (you need to detect this error), retry with delay
		// Here, you should replace `isRateLimitError(err)` with your actual logic to detect the rate limit error
		if isRateLimitError(err) && i < maxRetries-1 {
			delay := time.Duration(math.Pow(2, float64(i))) * baseDelay
			time.Sleep(delay)
			continue
		} else {
			break
		}
	}
	return reserves, err
}

// A placeholder function for rate limit error detection
func isRateLimitError(err error) bool {
	return strings.Contains(err.Error(), "429")
}
