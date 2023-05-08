package dex

import (
	"math/big"
	"time"
)

type EffectivePrice struct {
	DexID           DEX
	PoolID          string
	EffectivePrice0 *big.Float
	EffectivePrice1 *big.Float
	Timestamp       time.Time
}
