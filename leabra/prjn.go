// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leabra

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"cogentcore.org/core/gox/indent"
	"cogentcore.org/core/math32"
	"github.com/emer/emergent/v2/emer"
	"github.com/emer/emergent/v2/params"
	"github.com/emer/emergent/v2/path"
	"github.com/emer/emergent/v2/weights"
	"github.com/emer/etable/v2/etensor"
)

// Path is a basic Leabra pathway with synaptic learning parameters
type Path struct {
	PathBase

	// initial random weight distribution
	WtInit WtInitParams `view:"inline"`

	// weight scaling parameters: modulates overall strength of pathway, using both absolute and relative factors
	WtScale WtScaleParams `view:"inline"`

	// synaptic-level learning parameters
	Learn LearnSynParams `view:"add-fields"`

	// synaptic state values, ordered by the sending layer units which owns them -- one-to-one with SConIndex array
	Syns []Synapse

	// scaling factor for integrating synaptic input conductances (G's) -- computed in AlphaCycInit, incorporates running-average activity levels
	GScale float32

	// local per-recv unit increment accumulator for synaptic conductance from sending units -- goes to either GeRaw or GiRaw on neuron depending on pathway type -- this will be thread-safe
	GInc []float32

	// weight balance state variables for this pathway, one per recv neuron
	WbRecv []WtBalRecvPath
}

// AsLeabra returns this path as a leabra.Path -- all derived paths must redefine
// this to return the base Path type, so that the LeabraPath interface does not
// need to include accessors to all the basic stuff.
func (pj *Path) AsLeabra() *Path {
	return pj
}

func (pj *Path) Defaults() {
	pj.WtInit.Defaults()
	pj.WtScale.Defaults()
	pj.Learn.Defaults()
	pj.GScale = 1
}

// UpdateParams updates all params given any changes that might have been made to individual values
func (pj *Path) UpdateParams() {
	pj.WtScale.Update()
	pj.Learn.Update()
	pj.Learn.LrateInit = pj.Learn.Lrate
}

func (pj *Path) SetClass(cls string) emer.Path         { pj.Cls = cls; return pj }
func (pj *Path) SetPattern(pat path.Pattern) emer.Path { pj.Pat = pat; return pj }
func (pj *Path) SetType(typ emer.PathType) emer.Path   { pj.Typ = typ; return pj }

// AllParams returns a listing of all parameters in the Layer
func (pj *Path) AllParams() string {
	str := "///////////////////////////////////////////////////\nPath: " + pj.Name() + "\n"
	b, _ := json.MarshalIndent(&pj.WtInit, "", " ")
	str += "WtInit: {\n " + JsonToParams(b)
	b, _ = json.MarshalIndent(&pj.WtScale, "", " ")
	str += "WtScale: {\n " + JsonToParams(b)
	b, _ = json.MarshalIndent(&pj.Learn, "", " ")
	str += "Learn: {\n " + strings.Replace(JsonToParams(b), " XCal: {", "\n  XCal: {", -1)
	return str
}

// SetParam sets parameter at given path to given value.
// returns error if path not found or value cannot be set.
func (pj *Path) SetParam(path, val string) error {
	return params.SetParam(pj, path, val)
}

func (pj *Path) SynVarNames() []string {
	return SynapseVars
}

// SynVarProps returns properties for variables
func (pj *Path) SynVarProps() map[string]string {
	return SynapseVarProps
}

// SynIndex returns the index of the synapse between given send, recv unit indexes
// (1D, flat indexes). Returns -1 if synapse not found between these two neurons.
// Requires searching within connections for receiving unit.
func (pj *Path) SynIndex(sidx, ridx int) int {
	nc := int(pj.SConN[sidx])
	st := int(pj.SConIndexSt[sidx])
	for ci := 0; ci < nc; ci++ {
		ri := int(pj.SConIndex[st+ci])
		if ri != ridx {
			continue
		}
		return int(st + ci)
	}
	return -1
}

// SynVarIndex returns the index of given variable within the synapse,
// according to *this path's* SynVarNames() list (using a map to lookup index),
// or -1 and error message if not found.
func (pj *Path) SynVarIndex(varNm string) (int, error) {
	return SynapseVarByName(varNm)
}

// SynVarNum returns the number of synapse-level variables
// for this path.  This is needed for extending indexes in derived types.
func (pj *Path) SynVarNum() int {
	return len(SynapseVars)
}

// Syn1DNum returns the number of synapses for this path as a 1D array.
// This is the max idx for SynVal1D and the number of vals set by SynValues.
func (pj *Path) Syn1DNum() int {
	return len(pj.Syns)
}

// SynVal1D returns value of given variable index (from SynVarIndex) on given SynIndex.
// Returns NaN on invalid index.
// This is the core synapse var access method used by other methods,
// so it is the only one that needs to be updated for derived layer types.
func (pj *Path) SynVal1D(varIndex int, synIndex int) float32 {
	if synIndex < 0 || synIndex >= len(pj.Syns) {
		return math32.NaN()
	}
	if varIndex < 0 || varIndex >= pj.SynVarNum() {
		return math32.NaN()
	}
	sy := &pj.Syns[synIndex]
	return sy.VarByIndex(varIndex)
}

// SynValues sets values of given variable name for each synapse, using the natural ordering
// of the synapses (sender based for Leabra),
// into given float32 slice (only resized if not big enough).
// Returns error on invalid var name.
func (pj *Path) SynValues(vals *[]float32, varNm string) error {
	vidx, err := pj.LeabraPrj.SynVarIndex(varNm)
	if err != nil {
		return err
	}
	ns := len(pj.Syns)
	if *vals == nil || cap(*vals) < ns {
		*vals = make([]float32, ns)
	} else if len(*vals) < ns {
		*vals = (*vals)[0:ns]
	}
	for i := range pj.Syns {
		(*vals)[i] = pj.LeabraPrj.SynVal1D(vidx, i)
	}
	return nil
}

// SynVal returns value of given variable name on the synapse
// between given send, recv unit indexes (1D, flat indexes).
// Returns math32.NaN() for access errors (see SynValTry for error message)
func (pj *Path) SynValue(varNm string, sidx, ridx int) float32 {
	vidx, err := pj.LeabraPrj.SynVarIndex(varNm)
	if err != nil {
		return math32.NaN()
	}
	synIndex := pj.SynIndex(sidx, ridx)
	return pj.LeabraPrj.SynVal1D(vidx, synIndex)
}

// SetSynVal sets value of given variable name on the synapse
// between given send, recv unit indexes (1D, flat indexes)
// returns error for access errors.
func (pj *Path) SetSynValue(varNm string, sidx, ridx int, val float32) error {
	vidx, err := pj.LeabraPrj.SynVarIndex(varNm)
	if err != nil {
		return err
	}
	synIndex := pj.SynIndex(sidx, ridx)
	if synIndex < 0 || synIndex >= len(pj.Syns) {
		return err
	}
	sy := &pj.Syns[synIndex]
	sy.SetVarByIndex(vidx, val)
	if varNm == "Wt" {
		pj.Learn.LWtFmWt(sy)
	}
	return nil
}

///////////////////////////////////////////////////////////////////////
//  Weights File

// WriteWtsJSON writes the weights from this pathway from the receiver-side perspective
// in a JSON text format.  We build in the indentation logic to make it much faster and
// more efficient.
func (pj *Path) WriteWtsJSON(w io.Writer, depth int) {
	slay := pj.Send.(LeabraLayer).AsLeabra()
	rlay := pj.Recv.(LeabraLayer).AsLeabra()
	nr := len(rlay.Neurons)
	w.Write(indent.TabBytes(depth))
	w.Write([]byte("{\n"))
	depth++
	w.Write(indent.TabBytes(depth))
	w.Write([]byte(fmt.Sprintf("\"From\": %q,\n", slay.Name())))
	w.Write(indent.TabBytes(depth))
	w.Write([]byte(fmt.Sprintf("\"MetaData\": {\n")))
	depth++
	w.Write(indent.TabBytes(depth))
	w.Write([]byte(fmt.Sprintf("\"GScale\": \"%g\"\n", pj.GScale)))
	depth--
	w.Write(indent.TabBytes(depth))
	w.Write([]byte("},\n"))
	w.Write(indent.TabBytes(depth))
	w.Write([]byte(fmt.Sprintf("\"Rs\": [\n")))
	depth++
	for ri := 0; ri < nr; ri++ {
		nc := int(pj.RConN[ri])
		st := int(pj.RConIndexSt[ri])
		w.Write(indent.TabBytes(depth))
		w.Write([]byte("{\n"))
		depth++
		w.Write(indent.TabBytes(depth))
		w.Write([]byte(fmt.Sprintf("\"Ri\": %v,\n", ri)))
		w.Write(indent.TabBytes(depth))
		w.Write([]byte(fmt.Sprintf("\"N\": %v,\n", nc)))
		w.Write(indent.TabBytes(depth))
		w.Write([]byte("\"Si\": [ "))
		for ci := 0; ci < nc; ci++ {
			si := pj.RConIndex[st+ci]
			w.Write([]byte(fmt.Sprintf("%v", si)))
			if ci == nc-1 {
				w.Write([]byte(" "))
			} else {
				w.Write([]byte(", "))
			}
		}
		w.Write([]byte("],\n"))
		w.Write(indent.TabBytes(depth))
		w.Write([]byte("\"Wt\": [ "))
		for ci := 0; ci < nc; ci++ {
			rsi := pj.RSynIndex[st+ci]
			sy := &pj.Syns[rsi]
			w.Write([]byte(strconv.FormatFloat(float64(sy.Wt), 'g', weights.Prec, 32)))
			if ci == nc-1 {
				w.Write([]byte(" "))
			} else {
				w.Write([]byte(", "))
			}
		}
		w.Write([]byte("]\n"))
		depth--
		w.Write(indent.TabBytes(depth))
		if ri == nr-1 {
			w.Write([]byte("}\n"))
		} else {
			w.Write([]byte("},\n"))
		}
	}
	depth--
	w.Write(indent.TabBytes(depth))
	w.Write([]byte("]\n"))
	depth--
	w.Write(indent.TabBytes(depth))
	w.Write([]byte("}")) // note: leave unterminated as outer loop needs to add , or just \n depending
}

// ReadWtsJSON reads the weights from this pathway from the receiver-side perspective
// in a JSON text format.  This is for a set of weights that were saved *for one path only*
// and is not used for the network-level ReadWtsJSON, which reads into a separate
// structure -- see SetWts method.
func (pj *Path) ReadWtsJSON(r io.Reader) error {
	pw, err := weights.PathReadJSON(r)
	if err != nil {
		return err // note: already logged
	}
	return pj.SetWts(pw)
}

// SetWts sets the weights for this pathway from weights.Path decoded values
func (pj *Path) SetWts(pw *weights.Path) error {
	if pw.MetaData != nil {
		if gs, ok := pw.MetaData["GScale"]; ok {
			pv, _ := strconv.ParseFloat(gs, 32)
			pj.GScale = float32(pv)
		}
	}
	var err error
	for i := range pw.Rs {
		pr := &pw.Rs[i]
		for si := range pr.Si {
			er := pj.SetSynValue("Wt", pr.Si[si], pr.Ri, pr.Wt[si]) // updates lin wt
			if er != nil {
				err = er
			}
		}
	}
	return err
}

// Build constructs the full connectivity among the layers as specified in this pathway.
// Calls PathBase.BuildStru and then allocates the synaptic values in Syns accordingly.
func (pj *Path) Build() error {
	if err := pj.BuildStru(); err != nil {
		return err
	}
	pj.Syns = make([]Synapse, len(pj.SConIndex))
	rsh := pj.Recv.Shape()
	//	ssh := pj.Send.Shape()
	rlen := rsh.Len()
	pj.GInc = make([]float32, rlen)
	pj.WbRecv = make([]WtBalRecvPath, rlen)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////
//  Init methods

// SetScalesRPool initializes synaptic Scale values using given tensor
// of values which has unique values for each recv neuron within a given pool.
func (pj *Path) SetScalesRPool(scales etensor.Tensor) {
	rNuY := scales.Dim(0)
	rNuX := scales.Dim(1)
	rNu := rNuY * rNuX
	rfsz := scales.Len() / rNu

	rsh := pj.Recv.Shape()
	rNpY := rsh.Dim(0)
	rNpX := rsh.Dim(1)
	r2d := false
	if rsh.NumDims() != 4 {
		r2d = true
		rNpY = 1
		rNpX = 1
	}

	for rpy := 0; rpy < rNpY; rpy++ {
		for rpx := 0; rpx < rNpX; rpx++ {
			for ruy := 0; ruy < rNuY; ruy++ {
				for rux := 0; rux < rNuX; rux++ {
					ri := 0
					if r2d {
						ri = rsh.Offset([]int{ruy, rux})
					} else {
						ri = rsh.Offset([]int{rpy, rpx, ruy, rux})
					}
					scst := (ruy*rNuX + rux) * rfsz
					nc := int(pj.RConN[ri])
					st := int(pj.RConIndexSt[ri])
					for ci := 0; ci < nc; ci++ {
						// si := int(pj.RConIndex[st+ci]) // could verify coords etc
						rsi := pj.RSynIndex[st+ci]
						sy := &pj.Syns[rsi]
						sc := scales.FloatValue1D(scst + ci)
						sy.Scale = float32(sc)
					}
				}
			}
		}
	}
}

// SetWtsFunc initializes synaptic Wt value using given function
// based on receiving and sending unit indexes.
func (pj *Path) SetWtsFunc(wtFun func(si, ri int, send, recv *etensor.Shape) float32) {
	rsh := pj.Recv.Shape()
	rn := rsh.Len()
	ssh := pj.Send.Shape()

	for ri := 0; ri < rn; ri++ {
		nc := int(pj.RConN[ri])
		st := int(pj.RConIndexSt[ri])
		for ci := 0; ci < nc; ci++ {
			si := int(pj.RConIndex[st+ci])
			wt := wtFun(si, ri, ssh, rsh)
			rsi := pj.RSynIndex[st+ci]
			sy := &pj.Syns[rsi]
			sy.Wt = wt * sy.Scale
			pj.Learn.LWtFmWt(sy)
		}
	}
}

// SetScalesFunc initializes synaptic Scale values using given function
// based on receiving and sending unit indexes.
func (pj *Path) SetScalesFunc(scaleFun func(si, ri int, send, recv *etensor.Shape) float32) {
	rsh := pj.Recv.Shape()
	rn := rsh.Len()
	ssh := pj.Send.Shape()

	for ri := 0; ri < rn; ri++ {
		nc := int(pj.RConN[ri])
		st := int(pj.RConIndexSt[ri])
		for ci := 0; ci < nc; ci++ {
			si := int(pj.RConIndex[st+ci])
			sc := scaleFun(si, ri, ssh, rsh)
			rsi := pj.RSynIndex[st+ci]
			sy := &pj.Syns[rsi]
			sy.Scale = sc
		}
	}
}

// InitWtsSyn initializes weight values based on WtInit randomness parameters
// for an individual synapse.
// It also updates the linear weight value based on the sigmoidal weight value.
func (pj *Path) InitWtsSyn(syn *Synapse) {
	if syn.Scale == 0 {
		syn.Scale = 1
	}
	syn.Wt = float32(pj.WtInit.Gen(-1))
	// enforce normalized weight range -- required for most uses and if not
	// then a new type of path should be used:
	if syn.Wt < 0 {
		syn.Wt = 0
	}
	if syn.Wt > 1 {
		syn.Wt = 1
	}
	syn.LWt = pj.Learn.WtSig.LinFmSigWt(syn.Wt)
	syn.Wt *= syn.Scale // note: scale comes after so LWt is always "pure" non-scaled value
	syn.DWt = 0
	syn.Norm = 0
	syn.Moment = 0
}

// InitWts initializes weight values according to Learn.WtInit params
func (pj *Path) InitWts() {
	for si := range pj.Syns {
		sy := &pj.Syns[si]
		pj.InitWtsSyn(sy)
	}
	for wi := range pj.WbRecv {
		wb := &pj.WbRecv[wi]
		wb.Init()
	}
	pj.LeabraPrj.InitGInc()
}

// InitWtSym initializes weight symmetry -- is given the reciprocal pathway where
// the Send and Recv layers are reversed.
func (pj *Path) InitWtSym(rpjp LeabraPath) {
	rpj := rpjp.AsLeabra()
	slay := pj.Send.(LeabraLayer).AsLeabra()
	ns := int32(len(slay.Neurons))
	for si := int32(0); si < ns; si++ {
		nc := pj.SConN[si]
		st := pj.SConIndexSt[si]
		for ci := int32(0); ci < nc; ci++ {
			sy := &pj.Syns[st+ci]
			ri := pj.SConIndex[st+ci]
			// now we need to find the reciprocal synapse on rpj!
			// look in ri for sending connections
			rsi := ri
			if len(rpj.SConN) == 0 {
				continue
			}
			rsnc := rpj.SConN[rsi]
			if rsnc == 0 {
				continue
			}
			rsst := rpj.SConIndexSt[rsi]
			rist := rpj.SConIndex[rsst]        // starting index in recv path
			ried := rpj.SConIndex[rsst+rsnc-1] // ending index
			if si < rist || si > ried {        // fast reject -- paths are always in order!
				continue
			}
			// start at index proportional to si relative to rist
			up := int32(0)
			if ried > rist {
				up = int32(float32(rsnc) * float32(si-rist) / float32(ried-rist))
			}
			dn := up - 1

			for {
				doing := false
				if up < rsnc {
					doing = true
					rrii := rsst + up
					rri := rpj.SConIndex[rrii]
					if rri == si {
						rsy := &rpj.Syns[rrii]
						rsy.Wt = sy.Wt
						rsy.LWt = sy.LWt
						rsy.Scale = sy.Scale
						// note: if we support SymFmTop then can have option to go other way
						break
					}
					up++
				}
				if dn >= 0 {
					doing = true
					rrii := rsst + dn
					rri := rpj.SConIndex[rrii]
					if rri == si {
						rsy := &rpj.Syns[rrii]
						rsy.Wt = sy.Wt
						rsy.LWt = sy.LWt
						rsy.Scale = sy.Scale
						// note: if we support SymFmTop then can have option to go other way
						break
					}
					dn--
				}
				if !doing {
					break
				}
			}
		}
	}
}

// InitGInc initializes the per-pathway GInc threadsafe increment -- not
// typically needed (called during InitWts only) but can be called when needed
func (pj *Path) InitGInc() {
	for ri := range pj.GInc {
		pj.GInc[ri] = 0
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//  Act methods

// SendGDelta sends the delta-activation from sending neuron index si,
// to integrate synaptic conductances on receivers
func (pj *Path) SendGDelta(si int, delta float32) {
	scdel := delta * pj.GScale
	nc := pj.SConN[si]
	st := pj.SConIndexSt[si]
	syns := pj.Syns[st : st+nc]
	scons := pj.SConIndex[st : st+nc]
	for ci := range syns {
		ri := scons[ci]
		pj.GInc[ri] += scdel * syns[ci].Wt
	}
}

// RecvGInc increments the receiver's GeRaw or GiRaw from that of all the pathways.
func (pj *Path) RecvGInc() {
	rlay := pj.Recv.(LeabraLayer).AsLeabra()
	if pj.Typ == emer.Inhib {
		for ri := range rlay.Neurons {
			rn := &rlay.Neurons[ri]
			rn.GiRaw += pj.GInc[ri]
			pj.GInc[ri] = 0
		}
	} else {
		for ri := range rlay.Neurons {
			rn := &rlay.Neurons[ri]
			rn.GeRaw += pj.GInc[ri]
			pj.GInc[ri] = 0
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//  Learn methods

// DWt computes the weight change (learning) -- on sending pathways
func (pj *Path) DWt() {
	if !pj.Learn.Learn {
		return
	}
	slay := pj.Send.(LeabraLayer).AsLeabra()
	rlay := pj.Recv.(LeabraLayer).AsLeabra()
	for si := range slay.Neurons {
		sn := &slay.Neurons[si]
		if sn.AvgS < pj.Learn.XCal.LrnThr && sn.AvgM < pj.Learn.XCal.LrnThr {
			continue
		}
		nc := int(pj.SConN[si])
		st := int(pj.SConIndexSt[si])
		syns := pj.Syns[st : st+nc]
		scons := pj.SConIndex[st : st+nc]
		for ci := range syns {
			sy := &syns[ci]
			ri := scons[ci]
			rn := &rlay.Neurons[ri]
			err, bcm := pj.Learn.CHLdWt(sn.AvgSLrn, sn.AvgM, rn.AvgSLrn, rn.AvgM, rn.AvgL)

			bcm *= pj.Learn.XCal.LongLrate(rn.AvgLLrn)
			err *= pj.Learn.XCal.MLrn
			dwt := bcm + err
			norm := float32(1)
			if pj.Learn.Norm.On {
				norm = pj.Learn.Norm.NormFmAbsDWt(&sy.Norm, math32.Abs(dwt))
			}
			if pj.Learn.Momentum.On {
				dwt = norm * pj.Learn.Momentum.MomentFmDWt(&sy.Moment, dwt)
			} else {
				dwt *= norm
			}
			sy.DWt += pj.Learn.Lrate * dwt
		}
		// aggregate max DWtNorm over sending synapses
		if pj.Learn.Norm.On {
			maxNorm := float32(0)
			for ci := range syns {
				sy := &syns[ci]
				if sy.Norm > maxNorm {
					maxNorm = sy.Norm
				}
			}
			for ci := range syns {
				sy := &syns[ci]
				sy.Norm = maxNorm
			}
		}
	}
}

// WtFmDWt updates the synaptic weight values from delta-weight changes -- on sending pathways
func (pj *Path) WtFmDWt() {
	if !pj.Learn.Learn {
		return
	}
	if pj.Learn.WtBal.On {
		for si := range pj.Syns {
			sy := &pj.Syns[si]
			ri := pj.SConIndex[si]
			wb := &pj.WbRecv[ri]
			pj.Learn.WtFmDWt(wb.Inc, wb.Dec, &sy.DWt, &sy.Wt, &sy.LWt, sy.Scale)
		}
	} else {
		for si := range pj.Syns {
			sy := &pj.Syns[si]
			pj.Learn.WtFmDWt(1, 1, &sy.DWt, &sy.Wt, &sy.LWt, sy.Scale)
		}
	}
}

// WtBalFmWt computes the Weight Balance factors based on average recv weights
func (pj *Path) WtBalFmWt() {
	if !pj.Learn.Learn || !pj.Learn.WtBal.On {
		return
	}

	rlay := pj.Recv.(LeabraLayer).AsLeabra()
	if !pj.Learn.WtBal.Targs && rlay.LeabraLay.IsTarget() {
		return
	}
	for ri := range rlay.Neurons {
		nc := int(pj.RConN[ri])
		if nc < 1 {
			continue
		}
		wb := &pj.WbRecv[ri]
		st := int(pj.RConIndexSt[ri])
		rsidxs := pj.RSynIndex[st : st+nc]
		sumWt := float32(0)
		sumN := 0
		for ci := range rsidxs {
			rsi := rsidxs[ci]
			sy := &pj.Syns[rsi]
			if sy.Wt >= pj.Learn.WtBal.AvgThr {
				sumWt += sy.Wt
				sumN++
			}
		}
		if sumN > 0 {
			sumWt /= float32(sumN)
		} else {
			sumWt = 0
		}
		wb.Avg = sumWt
		wb.Fact, wb.Inc, wb.Dec = pj.Learn.WtBal.WtBal(sumWt)
	}
}

// LrateMult sets the new Lrate parameter for Paths to LrateInit * mult.
// Useful for implementing learning rate schedules.
func (pj *Path) LrateMult(mult float32) {
	pj.Learn.Lrate = pj.Learn.LrateInit * mult
}

///////////////////////////////////////////////////////////////////////
//  WtBalRecvPath

// WtBalRecvPath are state variables used in computing the WtBal weight balance function
// There is one of these for each Recv Neuron participating in the pathway.
type WtBalRecvPath struct {

	// average of effective weight values that exceed WtBal.AvgThr across given Recv Neuron's connections for given Path
	Avg float32

	// overall weight balance factor that drives changes in WbInc vs. WbDec via a sigmoidal function -- this is the net strength of weight balance changes
	Fact float32

	// weight balance increment factor -- extra multiplier to add to weight increases to maintain overall weight balance
	Inc float32

	// weight balance decrement factor -- extra multiplier to add to weight decreases to maintain overall weight balance
	Dec float32
}

func (wb *WtBalRecvPath) Init() {
	wb.Avg = 0
	wb.Fact = 0
	wb.Inc = 1
	wb.Dec = 1
}
