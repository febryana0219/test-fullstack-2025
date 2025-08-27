package main

import (
	"fmt"
	"math/big"
)

func hitungFaktorial(n int64) *big.Int {
	if n < 0 {
		return big.NewInt(0)
	}

	// faktorial
	f := big.NewInt(1)
	// pangkat
	p := big.NewInt(2)
	// sisa
	sisa := big.NewInt(0)

	// hitung n faktorial secara manual
	for i := int64(2); i <= n; i++ {
		f.Mul(f, big.NewInt(i))
	}

	// pangkatkab 2 sebanyak n
	p.Exp(p, big.NewInt(n), nil)

	hasilBagi, _ := new(big.Int).DivMod(f, p, sisa)

	// jika ada sisa bulatkan ke atas
	if sisa.Cmp(big.NewInt(0)) > 0 {
		hasilBagi.Add(hasilBagi, big.NewInt(1))
	}

	return hasilBagi
}

func main() {
	for i := int64(0); i <= 10; i++ {
		fmt.Printf("f(%d) = %s\n", i, hitungFaktorial(i).String())
	}
}
