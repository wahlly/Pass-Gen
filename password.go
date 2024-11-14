package main

import (
	"math"
	"math/rand"
	"slices"
	"strings"
)

const (
	half = .5
	onethird = .3
	onefourth = .25
)

var (
	randlowers = randFromSeed(lowers())
	randuppers = randFromSeed(uppers())
	randdigits = randFromSeed(digits())
	randsymbols = randFromSeed(symbols())
)

var basicPassword = randlowers

func mediumPassword(n int) string {
	frac := math.Round(float64(n) * half)
	pwd := basicPassword(n)
	return pwd[:n-int(frac)] + randuppers(int(frac))
}

func hardPassword(n int) string {
	pwd := mediumPassword(n)
	frac := math.Round(float64(n) * onefourth)
	return pwd[:n - int(frac)] + randsymbols(int(frac))
}

func xhardPassword(n int) string {
	pwd := hardPassword(n)
	frac := math.Round(float64(n) * onefourth)
	return pwd[:n-int(frac)] + randsymbols(int(frac))
}

func randFromSeed(seed string) func(int) string{
	return func(n int) string {
		var b strings.Builder
		for range n {
			b.WriteByte(seed[rand.Intn(len(seed))])
		}
		return b.String()
	}
}

func lowers() string {
	var b strings.Builder
	for i := 'a'; i < 'a'+26; i++ {
		b.WriteRune(i)
	}
	return b.String()
}

func uppers() string {
	var b strings.Builder
	for i := 'A'; i<'A'+26; i++ {
		b.WriteRune(i)
	}
	return b.String()
}

func symbols() string {
	var b strings.Builder
	for i := '!'; i < '!'+14; i++ {
		b.WriteRune(i)
	}
	for i := ':'; i < ':'+6; i++ {
		b.WriteRune(i)
	}
	for i := '['; i < '['+5; i++ {
		b.WriteRune(i)
	}
	for i := '{'; i < '{'+3; i++ {
		b.WriteRune(i)
	}

	return b.String()
}

func digits() string {
	var b strings.Builder
	for i := '0'; i < '0'+9; i++ {
		b.WriteRune(i)
	}

	return b.String()
}

func shuffle[T any](ts []T) []T {
	cloned := slices.Clone(ts)
	rand.Shuffle(len(cloned), func(i, j int) {
		cloned[i], cloned[j] = cloned[j], cloned[i]
	})

	return cloned
}

func shuffleStr(s string) string {
	return strings.Join(shuffle(strings.Split(s, "")), "")
}