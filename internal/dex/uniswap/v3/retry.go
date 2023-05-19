package v3

import (
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zhashkevych/trinity/pkg/web3"
)

// Define max retries and base delay
const maxRetries = 5
const baseDelay = 1 * time.Second

func callQuoteExactInputSingleWithRetry(quoterv2 *QuoterV2, inp CalculateEffectivePriceInput) ([]interface{}, error) {
	amountIn := web3.ToTokenUnitsF(inp.amountIn, inp.TokenInDecimals)
	amountInI, _ := amountIn.Int(nil)

	res := make([]interface{}, 0)

	var err error
	for i := 0; i < maxRetries; i++ {
		err = quoterv2.CallQuoteExactInputSingle(
			&bind.CallOpts{},
			&res,
			IQuoterV2QuoteExactInputSingleParams{
				TokenIn:           inp.TokenInAddr,
				TokenOut:          inp.TokenOutAddr,
				AmountIn:          amountInI,
				Fee:               inp.Fee,
				SqrtPriceLimitX96: big.NewInt(0),
			})

		// If no error break the loop
		if err == nil {
			// USED FOR DEBUGING
			// if i > 0 {
			// 	log.WithFields(log.Fields{
			// 		"source": "uniswap/v3/parser.go",
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
	return res, err
}

func isRateLimitError(err error) bool {
	return strings.Contains(err.Error(), "429")
}
