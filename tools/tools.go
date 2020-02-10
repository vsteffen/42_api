package tools

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

const (
	WagnerFischerInsertionCost    = 1
	WagnerFischerDeletionCost     = 1
	WagnerFischerSubstitutionCost = 1
)

// ReadAndHideData read on stdin and hide user input
func ReadAndHideData() string {
	byteRead, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal().Err(err)
	}
	return string(byteRead)
}

// StringToRuneSlice convert string to Slice Rune
func StringToRuneSlice(s string) []rune {
	var r []rune
	for _, runeValue := range s {
		r = append(r, runeValue)
	}
	return r
}

// MinInt return minimumInt
func MinInt(firstArg int, args ...int) int {
	min := firstArg
	for _, nb := range args {
		if min > nb {
			min = nb
		}
	}
	return min
}

// EditDistance return the edit distance based on the algorithm of Wagner-Fischer of the given two strings
func EditDistance(s, t string) int {
	ss := StringToRuneSlice(s)
	m := len(ss) + 1
	tt := StringToRuneSlice(t)
	n := len(tt) + 1

	d := make([][]int, m)
	for i := range d {
		d[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		d[i][0] = i
	}

	for j := 0; j < n; j++ {
		d[0][j] = j
	}

	for j := 1; j < n; j++ {
		for i := 1; i < m; i++ {
			var substitutionCost int
			if ss[i-1] == tt[j-1] {
				substitutionCost = 0
			} else {
				substitutionCost = WagnerFischerSubstitutionCost
			}
			d[i][j] = MinInt(
				d[i-1][j]+WagnerFischerInsertionCost,
				d[i][j-1]+WagnerFischerDeletionCost,
				d[i-1][j-1]+substitutionCost,
			)
		}
	}

	return d[m-1][n-1]
}
