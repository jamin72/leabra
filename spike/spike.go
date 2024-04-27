// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spike

import (
	"cogentcore.org/core/math32"
	"github.com/emer/leabra/v2/leabra"
)

// ActParams is full set of activation params including those from base
// leabra and the additional Spiking-specific ones.
type ActParams struct {
	leabra.ActParams

	// spiking parameters
	Spike SpikeParams `view:"inline"`
}

func (sk *ActParams) Defaults() {
	sk.ActParams.Defaults()
	sk.Spike.Defaults()
}

func (sk *ActParams) Update() {
	sk.ActParams.Update()
	sk.Spike.Update()
}

// CopyFromAct copies ActParams from source (e.g., rate-code params)
func (sk *ActParams) CopyFromAct(act *leabra.ActParams) {
	sk.ActParams = *act
	sk.Update()
}

func (sk *ActParams) SpikeVmFmG(nrn *leabra.Neuron) {
	updtVm := true
	if sk.Spike.Tr > 0 && nrn.ISI >= 0 && nrn.ISI < float32(sk.Spike.Tr) {
		updtVm = false // don't update the spiking vm during refract
	}

	nwVm := nrn.Vm
	if updtVm {
		ge := nrn.Ge * sk.Gbar.E
		gi := nrn.Gi * sk.Gbar.I
		gk := sk.Gbar.K * (nrn.GknaFast + nrn.GknaMed + nrn.GknaSlow)
		nrn.Gk = gk
		vmEff := nrn.Vm
		// midpoint method: take a half-step in vmEff
		inet1 := sk.InetFmG(vmEff, ge, gi, gk)
		vmEff += .5 * sk.Dt.VmDt * inet1 // go half way
		inet2 := sk.InetFmG(vmEff, ge, gi, gk)
		// add spike current if relevant
		if sk.Spike.Exp {
			inet2 += sk.Gbar.L * sk.Spike.ExpSlope *
				math32.Exp((vmEff-sk.XX1.Thr)/sk.Spike.ExpSlope)
		}
		nwVm += sk.Dt.VmDt * inet2
		nrn.Inet = inet2
	}

	if sk.Noise.Type == leabra.VmNoise {
		nwVm += nrn.Noise
	}
	nrn.Vm = sk.VmRange.ClipValue(nwVm)
}

// SpikeActFmVm computes the discrete spiking activation
// from membrane potential Vm
func (sk *ActParams) SpikeActFmVm(nrn *leabra.Neuron) {
	var thr float32
	if sk.Spike.Exp {
		thr = sk.Spike.ExpThr
	} else {
		thr = sk.XX1.Thr
	}
	if nrn.Vm > thr {
		nrn.Spike = 1
		nrn.Vm = sk.Spike.VmR
		nrn.Inet = 0
		if nrn.ISIAvg == -1 {
			nrn.ISIAvg = -2
		} else if nrn.ISI > 0 { // must have spiked to update
			sk.Spike.AvgFmISI(&nrn.ISIAvg, nrn.ISI+1)
		}
		nrn.ISI = 0
	} else {
		nrn.Spike = 0
		if nrn.ISI >= 0 {
			nrn.ISI += 1
		}
		if nrn.ISIAvg >= 0 && nrn.ISI > 0 && nrn.ISI > 1.2*nrn.ISIAvg {
			sk.Spike.AvgFmISI(&nrn.ISIAvg, nrn.ISI)
		}
	}

	nwAct := sk.Spike.ActFmISI(nrn.ISIAvg, .001, 1) // todo: use real #'s
	if nwAct > 1 {
		nwAct = 1
	}
	nwAct = nrn.Act + sk.Dt.VmDt*(nwAct-nrn.Act)
	nrn.ActDel = nwAct - nrn.Act
	nrn.Act = nwAct
	if sk.KNa.On {
		sk.KNa.GcFmSpike(&nrn.GknaFast, &nrn.GknaMed, &nrn.GknaSlow, nrn.Spike > .5)
	}
}

// SpikeParams contains spiking activation function params.
// Implements the AdEx adaptive exponential function
type SpikeParams struct {

	// if true, turn on exponential excitatory current that drives Vm rapidly upward for spiking as it gets past its nominal firing threshold (Thr) -- nicely captures the Hodgkin Huxley dynamics of Na and K channels -- uses Brette & Gurstner 2005 AdEx formulation -- this mechanism has an unfortunate interaction with the continuous inhibitory currents generated by the standard FFFB inhibitory function, which cause this mechanism to desensitize and fail to spike
	Exp bool `def:"false"`

	// slope in Vm (2 mV = .02 in normalized units) for extra exponential excitatory current that drives Vm rapidly upward for spiking as it gets past its nominal firing threshold (Thr) -- nicely captures the Hodgkin Huxley dynamics of Na and K channels -- uses Brette & Gurstner 2005 AdEx formulation -- a value of 0 disables this mechanism
	ExpSlope float32 `viewif:"Exp" def:"0.02"`

	// membrane potential threshold for actually triggering a spike when using the exponential mechanism
	ExpThr float32 `viewif:"Exp" def:"1.2"`

	// post-spiking membrane potential to reset to, produces refractory effect if lower than VmInit -- 0.30 is appropriate biologically based value for AdEx (Brette & Gurstner, 2005) parameters
	VmR float32 `def:"0.3,0,0.15"`

	// post-spiking explicit refractory period, in cycles -- prevents Vm updating for this number of cycles post firing
	Tr int `def:"3"`

	// for translating spiking interval (rate) into rate-code activation equivalent (and vice-versa, for clamped layers), what is the maximum firing rate associated with a maximum activation value (max act is typically 1.0 -- depends on act_range)
	MaxHz float32 `def:"180" min:"1"`

	// constant for integrating the spiking interval in estimating spiking rate
	RateTau float32 `def:"5" min:"1"`

	// rate = 1 / tau
	RateDt float32 `view:"-"`
}

func (sk *SpikeParams) Defaults() {
	sk.Exp = false
	sk.ExpSlope = 0.02
	sk.ExpThr = 1.2
	sk.VmR = 0.3
	sk.Tr = 3
	sk.MaxHz = 180
	sk.RateTau = 5
	sk.Update()
}

func (sk *SpikeParams) Update() {
	sk.RateDt = 1 / sk.RateTau
}

// ActToISI compute spiking interval from a given rate-coded activation,
// based on time increment (.001 = 1msec default), Act.Dt.Integ
func (sk *SpikeParams) ActToISI(act, timeInc, integ float32) float32 {
	if act == 0 {
		return 0
	}
	return (1 / (timeInc * integ * act * sk.MaxHz))
}

// ActFmISI computes rate-code activation from estimated spiking interval
func (sk *SpikeParams) ActFmISI(isi, timeInc, integ float32) float32 {
	if isi <= 0 {
		return 0
	}
	maxInt := 1.0 / (timeInc * integ * sk.MaxHz) // interval at max hz..
	return maxInt / isi                          // normalized
}

// AvgFmISI updates spiking ISI from current isi interval value
func (sk *SpikeParams) AvgFmISI(avg *float32, isi float32) {
	if *avg <= 0 {
		*avg = isi
	} else if isi < 0.8**avg {
		*avg = isi // if significantly less than we take that
	} else { // integrate on slower
		*avg += sk.RateDt * (isi - *avg) // running avg updt
	}
}
