package flaggio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorkohl/flaggio/internal/flaggio"
)

func TestDistributionList_Distribute(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	vrnt1 := &flaggio.Variant{ID: "1"}
	vrnt2 := &flaggio.Variant{ID: "2"}
	vrnt3 := &flaggio.Variant{ID: "3"}
	dstrbtn := flaggio.DistributionList{
		{Variant: vrnt1, Percentage: 20},
		{Variant: vrnt2, Percentage: 70},
		{Variant: vrnt3, Percentage: 10},
	}

	totalDistributions := 500000
	var vrnt1Count, vrnt2Count, vrnt3Count int
	for i := 0; i < totalDistributions; i++ {
		vrnt := dstrbtn.Distribute()
		switch vrnt {
		case vrnt1:
			vrnt1Count++
		case vrnt2:
			vrnt2Count++
		case vrnt3:
			vrnt3Count++
		}
	}

	assert.InDelta(t, dstrbtn[0].Percentage, float32(vrnt1Count)/float32(totalDistributions)*100, 0.2)
	assert.InDelta(t, dstrbtn[1].Percentage, float32(vrnt2Count)/float32(totalDistributions)*100, 0.2)
	assert.InDelta(t, dstrbtn[2].Percentage, float32(vrnt3Count)/float32(totalDistributions)*100, 0.2)
}

func BenchmarkDistributionList_Distribute(b *testing.B) {
	vrnt1 := &flaggio.Variant{ID: "1"}
	vrnt2 := &flaggio.Variant{ID: "2"}
	vrnt3 := &flaggio.Variant{ID: "3"}
	dstrbtn := flaggio.DistributionList{
		{Variant: vrnt1, Percentage: 20},
		{Variant: vrnt2, Percentage: 70},
		{Variant: vrnt3, Percentage: 10},
	}

	for n := 0; n < b.N; n++ {
		_ = dstrbtn.Distribute()
	}
}
