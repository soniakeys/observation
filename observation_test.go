// Public domain.

package observation_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/soniakeys/coord"
	"github.com/soniakeys/mpcformat"
)

func TestSiteObs_EarthObserverVect(t *testing.T) {
	m, err := mpcformat.ReadObscodeDat(bytes.NewBufferString(`F51 203.744090.936242+0.351541Pan-STARRS 1, Haleakala`))
	if err != nil {
		t.Fatal(err)
		return
	}
	_, o, err := mpcformat.ParseObs80("     K14G49E* C2014 04 09.45004 16 29 34.386+18 18 53.97         19.3 iL~133CF51", m)
	if err != nil {
		t.Fatal(err)
		return
	}
	// (desig, VMeas are tested in testObs80)

	eoWant := coord.Cart{
		X: -3.6643985193887653e-05,
		Y: -1.5829620927425095e-05,
		Z: 1.4988032341235873e-05,
	}
	eoGot := o.EarthObserverVect()
	if math.Abs(eoGot.X-eoWant.X) > 1e-15 ||
		math.Abs(eoGot.Y-eoWant.Y) > 1e-15 ||
		math.Abs(eoGot.Z-eoWant.Z) > 1e-15 {
		t.Fatalf("VObs.EarthObserverVect() = %+v, want %+v", eoGot, eoWant)
		return
	}
}
