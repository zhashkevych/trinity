package dex

import (
	"math/big"
	"time"
)

type EffectivePrice struct {
	DexID           DEX
	PoolID          string
	Reserve0        *big.Int
	Reserve1        *big.Int
	EffectivePrice0 *big.Float
	EffectivePrice1 *big.Float
	Timestamp       time.Time
}
