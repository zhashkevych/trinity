package v3

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhashkevych/dex-pools-aggregator/internal/dex"
	"github.com/zhashkevych/dex-pools-aggregator/pkg/web3"
)

var DexData = map[web3.Blockchain]map[string]string{
	web3.ETHEREUM: {
		"UniswapV3Factory": "0x1F98431c8aD98523631AE4a59f267346ea31F984",
	},
}

// TODO: parse from config.
var PoolsV3 = map[string]dex.Pool{
	"USDC / ETH": {
		Pair:         "USDC / ETH",
		Addr:         common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"),
		TokenOneAddr: TokensV3[web3.ETHEREUM][web3.USDC],
		TokenTwoAddr: TokensV3[web3.ETHEREUM][web3.ETH],
		Fee:          big.NewInt(500), // 0.05%
	},
	"WBTC / ETH": {
		Pair:         "WBTC / ETH",
		Addr:         common.HexToAddress("0xcbcdf9626bc03e24f779434178a73a0b4bad62ed"),
		TokenOneAddr: TokensV3[web3.ETHEREUM][web3.WBTC],
		TokenTwoAddr: TokensV3[web3.ETHEREUM][web3.ETH],
		Fee:          big.NewInt(300), // 0.03%
	},
}

// TODO: parse from config
var TokensV3 = map[web3.Blockchain]map[web3.Cryptocurrency]common.Address{
	web3.ETHEREUM: {
		web3.USDC: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // Token address for specific dex / pool ? Or general ?
		web3.ETH:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		web3.WBTC: common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	},
}
