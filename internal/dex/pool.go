package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TODO: store smart contract addresses in JSON config and parse into global config

type Pool struct {
	Pair             string
	Addr             common.Address
	TokenOneAddr     common.Address
	TokenTwoAddr     common.Address
	TokenOneAmountIn int64
	TokenTwoAmountIn int64
	Fee              *big.Int
}

// var PoolsV3 = map[string]Pool{
// 	"USDC / ETH": {
// 		Pair:         "USDC / ETH",
// 		Addr:         common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"),
// 		TokenOneAddr: TokensV3[ETHEREUM][USDC],
// 		TokenTwoAddr: TokensV3[ETHEREUM][ETH],
// 		Fee:          big.NewInt(500), // 0.05%
// 	},
// 	"WBTC / ETH": {
// 		Pair:         "WBTC / ETH",
// 		Addr:         common.HexToAddress("0xcbcdf9626bc03e24f779434178a73a0b4bad62ed"),
// 		TokenOneAddr: TokensV3[ETHEREUM][WBTC],
// 		TokenTwoAddr: TokensV3[ETHEREUM][ETH],
// 		Fee:          big.NewInt(300), // 0.03%
// 	},
// }

// var TokensV3 = map[string]map[string]common.Address{
// 	ETHEREUM: {
// 		USDC: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // Token address for specific dex / pool ? Or general ?
// 		ETH:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
// 		WBTC: common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
// 	},
// }
