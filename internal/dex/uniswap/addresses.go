package uniswap

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TODO: store smart contract addresses in JSON config and parse into global config

const (
	ETHEREUM = "ETHEREUM"
	USDC     = "USDC"
	ETH      = "ETH"
	WBTC     = "WBTC"
)

type Pool struct {
	Pair         string
	Addr         common.Address
	Fee          *big.Int
	TokenOneAddr common.Address
	TokenTwoAddr common.Address
}

var Pools = map[string]Pool{
	"USDC / ETH": {
		Pair:         "USDC / ETH",
		Addr:         common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"),
		TokenOneAddr: Tokens[ETHEREUM][USDC],
		TokenTwoAddr: Tokens[ETHEREUM][ETH],
		Fee:          big.NewInt(500), // 0.05%
	},
	"WBTC / ETH": {
		Pair:         "WBTC / ETH",
		Addr:         common.HexToAddress("0xcbcdf9626bc03e24f779434178a73a0b4bad62ed"),
		TokenOneAddr: Tokens[ETHEREUM][WBTC],
		TokenTwoAddr: Tokens[ETHEREUM][ETH],
		Fee:          big.NewInt(300), // 0.03%
	},
}

var Tokens = map[string]map[string]common.Address{
	ETHEREUM: {
		USDC: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // Token address for specific dex / pool ? Or general ?
		ETH:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		WBTC: common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	},
}

var DexData = map[string]map[string]string{
	ETHEREUM: {
		"UniswapV3Factory": "0x1F98431c8aD98523631AE4a59f267346ea31F984",
	},
}
