package v2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhashkevych/trinity/internal/dex"
	"github.com/zhashkevych/trinity/pkg/web3"
)

const (
	UniswapV2Fee int64 = 3 // 0.003
)

var DexData = map[web3.Blockchain]map[string]string{
	web3.ETHEREUM: {
		"UniswapV2Pair": "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
	},
}

// TODO: parse from config.
var PoolsV2 = map[string]*dex.PoolPair{
	"USDC / ETH": {
		Pair:             "USDC / ETH",
		Addr:             common.HexToAddress("0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc"),
		TokenOne:         web3.USDC,
		TokenTwo:         web3.ETH,
		TokenOneAddr:     dex.Tokens[web3.ETHEREUM][web3.USDC],
		TokenTwoAddr:     dex.Tokens[web3.ETHEREUM][web3.ETH],
		TokenOneAmountIn: big.NewInt(1000),
		TokenTwoAmountIn: big.NewInt(1),
		Fee:              big.NewInt(500), // 0.05%
	},
	"WBTC / ETH": {
		Pair:         "WBTC / ETH",
		Addr:         common.HexToAddress("0xcbcdf9626bc03e24f779434178a73a0b4bad62ed"),
		TokenOneAddr: dex.Tokens[web3.ETHEREUM][web3.WBTC],
		TokenTwoAddr: dex.Tokens[web3.ETHEREUM][web3.ETH],
		Fee:          big.NewInt(300), // 0.03%
	},
}
