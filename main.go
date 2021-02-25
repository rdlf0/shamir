package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type share struct {
	x int
	y big.Int
}

func main() {
	secret := big.NewInt(1234567890)

	prime, err := rand.Prime(rand.Reader, secret.BitLen()+1)
	if err != nil {
		panic(err)
	}
	shares := split(secret, 3, 6, prime)

	fmt.Printf("Secret: %s\nPrime: %s\n", secret.String(), prime.String())
	for i, s := range *shares {
		fmt.Printf("[%d] %+v\n", i, s)
	}

	fmt.Printf("Combine %d: %s\n", 1, combine([]share{(*shares)[1]}, prime).String())
	fmt.Printf("Combine %d and %d: %s\n", 0, 1, combine([]share{(*shares)[0], (*shares)[1]}, prime).String())
	fmt.Printf("Combine %d and %d: %s\n", 3, 4, combine([]share{(*shares)[3], (*shares)[4]}, prime).String())
	fmt.Printf("Combine %d and %d: %s\n", 1, 5, combine([]share{(*shares)[1], (*shares)[5]}, prime).String())
	fmt.Printf("Combine %d, %d and %d: %s\n", 0, 1, 2, combine([]share{(*shares)[0], (*shares)[1], (*shares)[2]}, prime).String())
	fmt.Printf("Combine %d, %d and %d: %s\n", 2, 4, 5, combine([]share{(*shares)[2], (*shares)[4], (*shares)[5]}, prime).String())
	fmt.Printf("Combine %d, %d and %d: %s\n", 1, 3, 4, combine([]share{(*shares)[1], (*shares)[3], (*shares)[4]}, prime).String())
	fmt.Printf("Combine %d, %d and %d: %s\n", 0, 2, 3, combine([]share{(*shares)[0], (*shares)[2], (*shares)[3]}, prime).String())
	fmt.Printf("Combine %d, %d, %d and %d: %s\n", 0, 2, 3, 5, combine([]share{(*shares)[0], (*shares)[2], (*shares)[3], (*shares)[5]}, prime).String())
	fmt.Printf("Combine %d, %d, %d, %d, %d and %d: %s\n", 0, 1, 2, 3, 4, 5, combine([]share{(*shares)[0], (*shares)[1], (*shares)[2], (*shares)[3], (*shares)[4], (*shares)[5]}, prime).String())
}

func split(secret *big.Int, needed, available int, prime *big.Int) *[]share {
	coef := make([]big.Int, needed)
	coef[0] = *secret
	for i := 1; i < needed; i++ {
		c, _ := rand.Int(rand.Reader, prime)
		coef[i] = *c
	}

	shares := make([]share, available)
	for x := 1; x <= available; x++ {
		acc := *secret
		for exp := 1; exp < needed; exp++ {
			xBig := big.NewInt(int64(x))
			eBig := big.NewInt(int64(exp))
			c := coef[exp]

			xBig.Exp(xBig, eBig, nil)
			c.Mul(&c, xBig)
			acc.Add(&acc, &c)
		}
		acc.Mod(&acc, prime)
		shares[x-1] = share{x, acc}
	}

	return &shares

}

func combine(shares []share, prime *big.Int) *big.Int {
	var secret big.Int

	for s := 0; s < len(shares); s++ {
		num := big.NewInt(1)
		den := big.NewInt(1)

		for c := 0; c < len(shares); c++ {
			if s == c {
				continue
			}

			start := shares[s].x
			next := shares[c].x
			diff := big.NewInt(int64(start - next))
			nextBig := big.NewInt(int64(next))
			nextBig.Neg(nextBig)
			num.Mul(num, nextBig)
			num.Mod(num, prime)

			den.Mul(den, diff)
			den.Mod(den, prime)
		}
		den.ModInverse(den, prime)

		var val big.Int
		val.Set(&shares[s].y)
		val.Mul(&val, num)
		val.Mul(&val, den)

		secret.Add(&secret, prime)
		secret.Add(&secret, &val)
		secret.Mod(&secret, prime)
	}

	return &secret
}
