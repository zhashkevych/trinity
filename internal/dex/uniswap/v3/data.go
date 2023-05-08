package v3

import (
	"github.com/zhashkevych/trinity/pkg/web3"
)

var DexData = map[web3.Blockchain]map[string]string{
	web3.ETHEREUM: {
		"UniswapV3Factory": "0x1F98431c8aD98523631AE4a59f267346ea31F984",
		"QuoterV2":         "0x61ffe014ba17989e743c5f6cb21bf9697530b21e",
	},
}
