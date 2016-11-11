package main

import (
	"math"

	"gopkg.in/distil.v1"
)


//This is our distillate algorithm
type addTwoDistiller struct {
	distil.DistillateTools
	basefreq int64
}

func (d *addTwoDistiller) LeadNanos() int64 {
	return 1000000000
}


// We will now show off all the tunables that a distillate MAY implement if
// it requires them. The first is "lead time". Some algorithms require some
// data outside of the changed range in order to compute. Frequency is one
// of those types of algorithms. By implementing LeadNanos() you can tell
// the DISTIL engine to load some extra data for you, which will be available
// at negative indices in the Process() method. The default implementation
// (in DistillateTools) returns 0

// PadSnap is a rebase stage that will adjust the incoming data to strictly
// appear on a timebase of the given frequency (hence the 'rebase'). Any
// values that do not appear with exactly the right time are snapped to
// the nearest time, and any duplicates are dropped. In addition, any missing
// values are replaced by NaN (hence pad). The advantage of this is that it
// simplifies calculations that refer to values across time, you can rest
// assured that a value 1s ago is exactly basefreq samples away, even if
// there were holes in the data or if there were duplicates. Note that in
// the presence of duplicate data there is ZERO GUARANTEE as to WHICH of the
// multiple duplicate values you receive. In general this makes algorithms
// that compare across time quite useless, as the real time difference
// between the points a fixed interval apart will experience extreme jitter.
// The default implementation (in DistillateTools) returns RebasePassthrough
func (d *addTwoDistiller) Rebase() distil.Rebaser {
	return distil.RebasePadSnap(d.basefreq)
}


// This is our main algorithm. It will automatically be called with chunks
// of data that require processing by the engine.
func (d *addTwoDistiller) Process(in *distil.InputSet, out *distil.OutputSet) {
	var ns int = in.NumSamples(0)
	var i int
	for i = 0; i < ns; i++ {
		point1 := in.Get(0, i).V
		point2 := in.Get(1, i).V
		sum := point1 + point2

		var time int64 = in.Get(0, i).T

		if !math.IsNaN(sum){
			out.Add(0, time, sum)
		}
	}
}
