package web3

import "math/big"

// ToTokenUnits converts a raw balance in the smallest token units to the full balance
func ToTokenUnits(rawBalance *big.Int, decimals int64) *big.Int {
	// return big.NewInt(0).Div(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	return big.NewInt(0).Mul(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}

func FromTokenUnits(rawBalance *big.Int, decimals int64) *big.Int {
	return big.NewInt(0).Div(rawBalance, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}
