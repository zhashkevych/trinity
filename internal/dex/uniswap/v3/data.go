package v3

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zhashkevych/trinity/internal/dex"
	"github.com/zhashkevych/trinity/pkg/web3"
)

var DexData = map[web3.Blockchain]map[string]string{
	web3.ETHEREUM: {
		"UniswapV3Factory": "0x1F98431c8aD98523631AE4a59f267346ea31F984",
		"QuoterV2":         "0x61ffe014ba17989e743c5f6cb21bf9697530b21e",
	},
}

// TODO: parse from config.
var PoolsV3 = map[string]*dex.PoolPair{
	"USDC / ETH": {
		Pair:             "USDC / ETH",
		Addr:             common.HexToAddress("0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8"),
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
