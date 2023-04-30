package dex

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zhashkevych/trinity/pkg/web3"
)

// TODO: parse from config
var Tokens = map[web3.Blockchain]map[web3.Cryptocurrency]common.Address{
	web3.ETHEREUM: {
		web3.USDC: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // Token address for specific dex / pool ? Or general ?
		web3.ETH:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		web3.WBTC: common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	},
}
