package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
