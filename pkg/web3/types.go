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
