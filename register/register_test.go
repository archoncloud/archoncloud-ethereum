package register

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pariz/gountries"
)

func TestGetEthPair(t *testing.T) {
	{
		// Example setting country code
		JP := "JP"
		countryCode := gountries.Codes{Alpha2: JP}
		spParams := SPParams{CountryCode: countryCode}

		query := gountries.New()
		japan, _ := query.FindCountryByAlpha(spParams.CountryCode.Alpha2)
		assert.Equal(t, JP, japan.Alpha2, "incorrect alpha2 assignment")
	}
	{
		// testing func encodeCountryCode(alpha2cc gountries.Codes) (int, error)

		// Example setting country code
		JP := "JP" // 9 X 15 = 135
		countryCode := gountries.Codes{Alpha2: JP}
		encCC, _ := encodeCountryCode(countryCode)
		var expectedBool = true
		var areEqual = true
		var expectedBytes [2]byte
		// 9  =b 00001001 to5bits 01001
		// 15 =b 00001111 "" ""   01111
		// {1: 00000001 0: 11101001} => {1: 1, 0: 233}
		expectedBytes[0] = byte(233)
		expectedBytes[1] = byte(1)
		for i := 0; i < len(expectedBytes); i++ {
			if expectedBytes[i] != encCC[i] {
				areEqual = false
			}
		}
		assert.Equal(t, expectedBool, areEqual, "incorrect countryCode numeric assignment")
	}
	{
		// testing func decodeCountryCode(tenBits [2]byte) gountries.Codes
		// Example setting country code
		JP := "JP" // 9 X 15 = 135
		countryCode := gountries.Codes{Alpha2: JP}
		encCC, _ := encodeCountryCode(countryCode)
		decCC := decodeCountryCode(encCC)
		assert.Equal(t, countryCode.Alpha2, decCC.Alpha2, "decodeCountryCode failed")
	}
	{
		var slaLevel = rand.Intn(8) // 0 has been explicitly tested
		var availableStorage = rand.Uint64()
		var bandwidth = rand.Uint64()
		var minAskPrice = rand.Uint64()
		JP := "JP" // 9 X 15 = 135
		countryCode := gountries.Codes{Alpha2: JP}

		p := SPParams{
			SLALevel:       slaLevel,
			PledgedStorage: uint64(availableStorage), // GB
			Bandwidth:      uint64(bandwidth),        // GB
			CountryCode:    countryCode,
			MinAskPrice:    uint64(minAskPrice)}

		encParams, _ := encodeParams(p)
		decParams := DecodeParams(encParams)

		var same bool = true
		if decParams.SLALevel != p.SLALevel {
			same = false
		}
		if decParams.PledgedStorage != p.PledgedStorage {
			same = false
		}
		if decParams.Bandwidth != p.Bandwidth {
			same = false
		}
		if decParams.CountryCode.Alpha2 != p.CountryCode.Alpha2 {
			same = false
		}
		if decParams.MinAskPrice != p.MinAskPrice {
			same = false
		}
		assert.Equal(t, true, same, "encodeParams/DecodeParams failed")
	}
}

func TestCheckIfPmtIsEnoughtForRegTx(t *testing.T) {
	//func checkIfPmtIsEnoughForRegTx(params SPParams) bool

	registerCost := uint64(100000000000) // goal 1e15
	registerCostScalar := uint64(10000)  // goal 1e15
	goodStake := registerCost * registerCostScalar
	// note: if sc scalars change, this test needs to be revised

	params := new(SPParams)
	params.Stake = goodStake
	enough := checkIfPmtIsEnoughForRegTx(*params)
	assert.Equal(t, true, enough, "checkIfPmtIsEnoughForRegTx failed")

	badStake := goodStake - 1
	params.Stake = badStake
	enough = checkIfPmtIsEnoughForRegTx(*params)
	assert.Equal(t, false, enough, "checkIfPmtIsEnoughForRegTx failed")
}
