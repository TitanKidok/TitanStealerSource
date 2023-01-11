package random

import (
	crypto "crypto/rand"
	"fmt"
	"math/big"
)

func Random() string {
	safeNum, _ := crypto.Int(crypto.Reader, big.NewInt(999))
	return fmt.Sprint(safeNum.Int64())
}
