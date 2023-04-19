package data

import "github.com/ethereum/go-ethereum/common"

var Tokens = map[string]map[string]common.Address{
	"ETHEREUM": {
		"USDC": common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		"WETH": common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
	},
}

var DexData = map[string]map[string]string{
	"ETHEREUM": {
		"UniswapV3Factory": "0x1F98431c8aD98523631AE4a59f267346ea31F984",
	},
}
