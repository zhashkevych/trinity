package web3

type Blockchain string

const (
	ETHEREUM Blockchain = "ETHEREUM"
)

type Cryptocurrency string

const (
	USDC Cryptocurrency = "USDC"
	ETH  Cryptocurrency = "ETH"
	WBTC Cryptocurrency = "WBTC"
)

func (c Cryptocurrency) GetMultiplicator() int64 {
	return cryptocurencyMultiplicator[c]
}

var cryptocurencyMultiplicator = map[Cryptocurrency]int64{
	USDC: 6,
	ETH:  18,
}

// var baseUnitsMultiplicators = map[int]int64{
// 	0:  1,
// 	2:  100,
// 	3:  1000,
// 	4:  10000,
// 	5:  100000,
// 	6:  1000000,
// 	8:  100000000,
// 	18: 1000000000000000000,
// }
