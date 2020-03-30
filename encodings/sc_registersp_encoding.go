package encodings

import (
	"fmt"
	"strings"

	"github.com/archoncloud/archoncloud-ethereum/wallet"

	"github.com/pariz/gountries"
)

var MaxSLALevel int = 8 // this can be anything

type SPParams struct {
	Wallet         wallet.EthereumKeyset
	SLALevel       int
	PledgedStorage uint64
	Bandwidth      uint64

	CountryCode gountries.Codes // must contain A2 field
	MinAskPrice uint64          // Wei per MByte

	Stake         uint64
	HardwareProof [32]byte

	NodeID string
}

func EncodeParams(params SPParams) (res [32]byte, err error) {
	var ret [32]byte
	// SLALevel int (1..3 or so) // reserve 3 bits for up to 8 levels
	if params.SLALevel > 0 && params.SLALevel <= MaxSLALevel {
		ret[0] = byte(params.SLALevel)
	} else {
		var empty [32]byte
		return empty, fmt.Errorf(
			"registersp.encodeParams error: Invalid SLALevel! Range is [1, 8]")
	}
	// PledgedStorage uint64
	ret[1] = byte((params.PledgedStorage & 0xFF00000000000000) >> 56)
	ret[2] = byte((params.PledgedStorage & 0x00FF000000000000) >> 48)
	ret[3] = byte((params.PledgedStorage & 0x0000FF0000000000) >> 40)
	ret[4] = byte((params.PledgedStorage & 0x000000FF00000000) >> 32)
	ret[5] = byte((params.PledgedStorage & 0x00000000FF000000) >> 24)
	ret[6] = byte((params.PledgedStorage & 0x0000000000FF0000) >> 16)
	ret[7] = byte((params.PledgedStorage & 0x000000000000FF00) >> 8)
	ret[8] = byte((params.PledgedStorage & 0x00000000000000FF) >> 0)

	// Bandwidth uint64
	ret[9] = byte((params.Bandwidth & 0xFF00000000000000) >> 56)
	ret[10] = byte((params.Bandwidth & 0x00FF000000000000) >> 48)
	ret[11] = byte((params.Bandwidth & 0x0000FF0000000000) >> 40)
	ret[12] = byte((params.Bandwidth & 0x000000FF00000000) >> 32)
	ret[13] = byte((params.Bandwidth & 0x00000000FF000000) >> 24)
	ret[14] = byte((params.Bandwidth & 0x0000000000FF0000) >> 16)
	ret[15] = byte((params.Bandwidth & 0x000000000000FF00) >> 8)
	ret[16] = byte((params.Bandwidth & 0x00000000000000FF) >> 0)

	// MinAskPrice uint64
	ret[17] = byte((params.MinAskPrice & 0xFF00000000000000) >> 56)
	ret[18] = byte((params.MinAskPrice & 0x00FF000000000000) >> 48)
	ret[19] = byte((params.MinAskPrice & 0x0000FF0000000000) >> 40)
	ret[20] = byte((params.MinAskPrice & 0x000000FF00000000) >> 32)
	ret[21] = byte((params.MinAskPrice & 0x00000000FF000000) >> 24)
	ret[22] = byte((params.MinAskPrice & 0x0000000000FF0000) >> 16)
	ret[23] = byte((params.MinAskPrice & 0x000000000000FF00) >> 8)
	ret[24] = byte((params.MinAskPrice & 0x00000000000000FF) >> 0)

	encCC, err := encodeCountryCode(params.CountryCode)
	if err != nil {
		var empty [32]byte
		return empty, err
	}
	ret[25] |= encCC[0]
	ret[26] |= encCC[1]
	// NOTE THERE ARE SIX BITS AVAILABLE IN IDX 26 SINCE encCC is only 10 bits
	// NOTE IDX's 27 to 31 are available

	return ret, nil
}

func DecodeParams(params [32]byte) (res *SPParams) {
	ret := new(SPParams)
	// SLALEVEL
	ret.SLALevel = int(params[0])
	// AVAILABLE STORAGE
	ret.PledgedStorage = 0
	for i := 0; i < 8; i++ {
		ret.PledgedStorage += uint64(params[i+1]) << uint((7-i)*8)
	}
	// BANDWIDTH
	ret.Bandwidth = 0
	for i := 0; i < 8; i++ {
		ret.Bandwidth += uint64(params[i+9]) << uint((7-i)*8)
	}
	// MINASKPRICE
	ret.MinAskPrice = 0
	for i := 0; i < 8; i++ {
		ret.MinAskPrice += uint64(params[i+17]) << uint((7-i)*8)
	}
	var bCC [2]byte
	copy(bCC[:], params[25:27])
	decCC := decodeCountryCode(bCC)
	ret.CountryCode = decCC

	return ret
}

func encodeCountryCode(alpha2cc gountries.Codes) ([2]byte, error) {
	// encodes using 10 bits
	var ret [2]byte
	if alpha2cc.Alpha2 == "" {
		return ret, fmt.Errorf(
			"registersp.encodeCountryCode error: CountryCode must be ISO Alpha2")
	}
	ALPHA := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes := []rune(alpha2cc.Alpha2)
	a := strings.Index(ALPHA, string(runes[0]))
	b := strings.Index(ALPHA, string(runes[1]))
	if a < 0 || b < 0 {
		return ret, fmt.Errorf(
			"registersp.encodeCountryCode error: CountryCode must be ISO Alpha2")
	}
	ret[0] = byte(a)
	ret[0] |= byte(b) << 5
	ret[1] |= byte(b) >> 3
	return ret, nil
}

func decodeCountryCode(tenBits [2]byte) gountries.Codes {
	var ret gountries.Codes
	idx0 := int(tenBits[0] & 0x1F)
	idx1 := int((tenBits[0]&0xE0)>>5) + int((tenBits[1]&0x03)<<3)
	ALPHA := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes := []rune(ALPHA)
	a0 := string(runes[idx0])
	a1 := string(runes[idx1])
	cc := a0 + a1
	ret.Alpha2 = string(cc)
	return ret
}
