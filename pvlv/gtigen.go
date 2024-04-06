// Code generated by "core generate"; DO NOT EDIT.

package pvlv

import (
	"cogentcore.org/core/gti"
)

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IAmygPrjn", IDName: "i-amyg-prjn", Doc: "IAmygPrjn has one method, AsAmygModPrjn, which recasts the projection as a moddulatory projection"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ISetScalePrjn", IDName: "i-set-scale-prjn", Doc: "ISetScalePrjn initializes weights, including special scale calculations"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.AmygModPrjn", IDName: "amyg-mod-prjn", Doc: "AmygModPrjn holds parameters and state variables for modulatory projections to amygdala layers", Embeds: []gti.Field{{Name: "Prjn"}}, Fields: []gti.Field{{Name: "SetScale", Doc: "only for Leabra algorithm: if initializing the weights, set the connection scaling parameter in addition to intializing the weights -- for specifically-supported specs, this will for example set a gaussian scaling parameter on top of random initial weights, instead of just setting the initial weights to a gaussian weighted value -- for other specs that do not support a custom init_wts function, this will set the scale values to what the random weights would otherwise be set to, and set the initial weight value to a constant (init_wt_val)"}, {Name: "SetScaleMin", Doc: "minimum scale value for SetScale projections"}, {Name: "SetScaleMax", Doc: "maximum scale value for SetScale projections"}, {Name: "InitWtVal", Doc: "constant initial weight value for specs that do not support a custom init_wts function and have set_scale set: the scale values are set to what the random weights would otherwise be set to, and the initial weight value is set to this constant: the net actual weight value is scale * init_wt_val.."}, {Name: "DALRGain", Doc: "gain multiplier on abs(DA) learning rate multiplier"}, {Name: "DALRBase", Doc: "constant baseline amount of learning prior to abs(DA) factor -- should be near zero otherwise offsets in activation will drive learning in the absence of DA significance"}, {Name: "DALrnThr", Doc: "minimum threshold for phasic abs(da) signals to count as non-zero;  useful to screen out spurious da signals due to tiny VSPatch-to-LHb signals on t2 & t4 timesteps that can accumulate over many trials - 0.02 seems to work okay"}, {Name: "ActDeltaThr", Doc: "minimum threshold for delta activation to count as non-zero;  useful to screen out spurious learning due to unintended delta activity - 0.02 seems to work okay for both acquisition and extinction guys"}, {Name: "ActLrnMod", Doc: "if true, recv unit deep_lrn value modulates learning"}, {Name: "ActLrnThr", Doc: "only ru->deep_lrn values > this get to learn - 0.05f seems to work okay"}, {Name: "DaMod", Doc: "parameters for dopaminergic modulation"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IBlAmygLayer", IDName: "i-bl-amyg-layer", Doc: "IBlAmygLayer has one method, AsBlAmygLayer, that returns a pointer to the layer specifically as a BLA layer."})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.BlAmygLayer", IDName: "bl-amyg-layer", Doc: "BlAmygLayer contains values specific to BLA layers, including Interlayer Inhibition (ILI)", Embeds: []gti.Field{{Name: "ModLayer", Doc: "modulation state"}}, Fields: []gti.Field{{Name: "Valence", Doc: "positive or negative valence"}, {Name: "ILI", Doc: "inter-layer inhibition parameters and state"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ICElAmygLayer", IDName: "ic-el-amyg-layer"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.AcqExt", IDName: "acq-ext"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.CElAmygLayer", IDName: "c-el-amyg-layer", Embeds: []gti.Field{{Name: "ModLayer"}}, Fields: []gti.Field{{Name: "CElTyp", Doc: "basic parameters determining what type CEl layer this is"}, {Name: "AcqDeepMod", Doc: "use deep_mod_net for value from acquisition / go units, instead of inhibition current (otherwise use gi_syn) -- allows simpler parameter setting without titrating inhibition and this learning modulation signal"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.CElAmygLayerType", IDName: "c-el-amyg-layer-type", Fields: []gti.Field{{Name: "AcqExt", Doc: "acquisition or extinction"}, {Name: "Valence", Doc: "positive or negative DA valence"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Stim", IDName: "stim"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Context", IDName: "context"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Valence", IDName: "valence", Doc: "Valence"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IUS", IDName: "ius", Doc: "US, either positive or negative Valence"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.US", IDName: "us"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.PosUS", IDName: "pos-us", Doc: "positive and negative subtypes of US"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.NegUS", IDName: "neg-us"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Tick", IDName: "tick", Doc: "Tick"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.USTimeState", IDName: "us-time-state", Fields: []gti.Field{{Name: "Stm", Doc: "CS value"}, {Name: "US", Doc: "a US value or absent (USNone)"}, {Name: "Val", Doc: "PV d, POS, NEG, or absent (ValNone)"}, {Name: "Tck", Doc: "Within-trial timestep"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.PackedUSTimeState", IDName: "packed-us-time-state"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Inputs", IDName: "inputs"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.LHbRMTgGains", IDName: "l-hb-rm-tg-gains", Doc: "Gain constants for LHbRMTg inputs", Fields: []gti.Field{{Name: "All", Doc: "final overall gain on everything"}, {Name: "VSPatchPosD1", Doc: "patch D1 APPETITIVE pathway - versus pos PV outcomes"}, {Name: "VSPatchPosD2", Doc: "patch D2 APPETITIVE pathway versus vspatch_pos_D1"}, {Name: "VSPatchPosDisinhib", Doc: "proportion of positive reward prediction error (RPE) to use if RPE results from a predicted omission of positive"}, {Name: "VSMatrixPosD1", Doc: "gain on VS matrix D1 APPETITIVE guys"}, {Name: "VSMatrixPosD2", Doc: "VS matrix D2 APPETITIVE"}, {Name: "VSPatchNegD1", Doc: "VS patch D1 pathway versus neg PV outcomes"}, {Name: "VSPatchNegD2", Doc: "VS patch D2 pathway versus vspatch_neg_D1"}, {Name: "VSMatrixNegD1", Doc: "VS matrix D1 AVERSIVE"}, {Name: "VSMatrixNegD2", Doc: "VS matrix D2 AVERSIVE"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.LHbRMTgLayer", IDName: "l-hb-rm-tg-layer", Embeds: []gti.Field{{Name: "Layer"}}, Fields: []gti.Field{{Name: "RcvFrom"}, {Name: "Gains"}, {Name: "PVNegDiscount", Doc: "reduction in effective PVNeg net value (when positive) so that negative outcomes can never be completely predicted away -- still allows for positive da for less-bad outcomes"}, {Name: "InternalState"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.LHBRMTgInternalState", IDName: "lhbrm-tg-internal-state", Fields: []gti.Field{{Name: "VSPatchPosD1"}, {Name: "VSPatchPosD2"}, {Name: "VSPatchNegD1"}, {Name: "VSPatchNegD2"}, {Name: "VSMatrixPosD1"}, {Name: "VSMatrixPosD2"}, {Name: "VSMatrixNegD1"}, {Name: "VSMatrixNegD2"}, {Name: "PosPV"}, {Name: "NegPV"}, {Name: "VSPatchPosNet"}, {Name: "VSPatchNegNet"}, {Name: "VSMatrixPosNet"}, {Name: "VSMatrixNegNet"}, {Name: "NetPos"}, {Name: "NetNeg"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IModLayer", IDName: "i-mod-layer"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.AvgMaxModLayer", IDName: "avg-max-mod-layer"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModSender", IDName: "mod-sender", Doc: "ModSender has methods for sending modulation, and setting the value to be sent."})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModReceiver", IDName: "mod-receiver", Doc: "ModReceiver has one method to integrate incoming modulation, and another"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModLayer", IDName: "mod-layer", Doc: "ModLayer is a layer that RECEIVES modulatory input", Embeds: []gti.Field{{Name: "Layer"}, {Name: "ModParams", Doc: "parameters shared by all modulator receiver layers"}, {Name: "Modulators", Doc: "layer-level neuromodulator levels"}}, Fields: []gti.Field{{Name: "ModNeurs", Doc: "neuron-level modulation state"}, {Name: "ModPools", Doc: "pools for maintaining aggregate values"}, {Name: "ModReceivers", Doc: "layer names and scale values for mods sent from this layer"}, {Name: "DaMod", Doc: "parameters for dopaminergic modulation"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModPool", IDName: "mod-pool", Doc: "ModPool is similar to a standard Pool structure, and uses the same code to compute running statistics.", Fields: []gti.Field{{Name: "ModNetStats"}, {Name: "ModSent", Doc: "modulation level transmitted to receiver layers"}, {Name: "ModSendThreshold", Doc: "threshold for sending modulation. values below this are not added to the pool-level total"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.DaModParams", IDName: "da-mod-params", Doc: "DaModParams specifies parameters shared by all layers that receive dopaminergic modulatory input.", Fields: []gti.Field{{Name: "On", Doc: "whether to use dopamine modulation"}, {Name: "RecepType", Doc: "dopamine receptor type, D1 or D2"}, {Name: "BurstGain", Doc: "multiplicative gain factor applied to positive dopamine signals -- this operates on the raw dopamine signal prior to any effect of D2 receptors in reversing its sign!"}, {Name: "DipGain", Doc: "multiplicative gain factor applied to negative dopamine signals -- this operates on the raw dopamine signal prior to any effect of D2 receptors in reversing its sign! should be small for acq, but roughly equal to burst_da_gain for ext"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModParams", IDName: "mod-params", Doc: "ModParams contains values that control a receiving layer's response to modulatory inputs", Fields: []gti.Field{{Name: "Minus", Doc: "how much to multiply Da in the minus phase to add to Ge input -- use negative values for NoGo/indirect pathway/D2 type neurons"}, {Name: "Plus", Doc: "how much to multiply Da in the plus phase to add to Ge input -- use negative values for NoGo/indirect pathway/D2 type neurons"}, {Name: "NegGain", Doc: "for negative dopamine, how much to change the default gain value as a function of dopamine: gain = gain * (1 + da * NegNain) -- da is multiplied by minus or plus depending on phase"}, {Name: "PosGain", Doc: "for positive dopamine, how much to change the default gain value as a function of dopamine: gain = gain * (1 + da * PosGain) -- da is multiplied by minus or plus depending on phase"}, {Name: "ActModZero", Doc: "for modulation coming from the BLA via deep_mod_net -- when this modulation signal is below zero, does it have the ability to zero out the patch activations?  i.e., is the modulation required to enable patch firing?"}, {Name: "ModNetThreshold", Doc: "threshold on deep_mod_net before deep mod is applied -- if not receiving even this amount of overall input from deep_mod sender, then do not use the deep_mod_net to drive deep_mod and deep_lrn values -- only for SUPER units -- based on LAYER level maximum for base LeabraLayerSpec, PVLV classes are based on actual deep_mod_net for each unit"}, {Name: "ModSendThreshold", Doc: "threshold for including neuron activation in total to send (for ModNet)"}, {Name: "IsModSender", Doc: "does this layer send modulation to other layers?"}, {Name: "IsModReceiver", Doc: "does this layer receive modulation from other layers?"}, {Name: "IsPVReceiver", Doc: "does this layer receive a direct PV input?"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModRcvrParams", IDName: "mod-rcvr-params", Doc: "ModRcvrParams specifies the name of a layer that receives modulatory input, and a scale factor--critical for inputs from\nlarge layers such as BLA.", Fields: []gti.Field{{Name: "RcvName", Doc: "name of receiving layer"}, {Name: "Scale", Doc: "scale factor for modulation to this receiver"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Modulators", IDName: "modulators", Doc: "Modulators are modulatory neurotransmitters. Currently ACh and SE are only placeholders.", Fields: []gti.Field{{Name: "DA", Doc: "current dopamine level for this layer"}, {Name: "ACh", Doc: "current acetylcholine level for this layer"}, {Name: "SE", Doc: "current serotonin level for this layer"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModNeuron", IDName: "mod-neuron", Doc: "ModNeuron encapsulates the variables used by all layers that receive modulatory input", Embeds: []gti.Field{{Name: "Modulators", Doc: "neuron-level modulator activation"}}, Fields: []gti.Field{{Name: "ModAct", Doc: "activity level for modulation"}, {Name: "ModLevel", Doc: "degree of full modulation to apply"}, {Name: "ModNet", Doc: "modulation input from sender"}, {Name: "ModLrn", Doc: "multiplier for DA modulation of learning rate"}, {Name: "PVAct", Doc: "direct activation from US"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.DaRType", IDName: "da-r-type", Doc: "Dopamine receptor type, for D1R and D2R dopamine receptors"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.MSNLayer", IDName: "msn-layer", Embeds: []gti.Field{{Name: "ModLayer"}}, Fields: []gti.Field{{Name: "Compartment", Doc: "patch or matrix"}, {Name: "DIState", Doc: "slice of delayed inhibition state for this layer."}, {Name: "DIParams"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IMSNLayer", IDName: "imsn-layer"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.MSNParams", IDName: "msn-params", Doc: "Parameters for Dorsal Striatum Medium Spiny Neuron computation", Fields: []gti.Field{{Name: "Compartment", Doc: "patch or matrix"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.StriatalCompartment", IDName: "striatal-compartment"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.DelayedInhibParams", IDName: "delayed-inhib-params", Doc: "Delayed inhibition for matrix compartment layers", Fields: []gti.Field{{Name: "Active", Doc: "add in a portion of inhibition from previous time period"}, {Name: "PrvQ", Doc: "proportion of per-unit net input on previous gamma-frequency quarter to add in as inhibition"}, {Name: "PrvTrl", Doc: "proportion of per-unit net input on previous trial to add in as inhibition"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.MSNTraceParams", IDName: "msn-trace-params", Doc: "Params for for trace-based learning", Fields: []gti.Field{{Name: "Deriv", Doc: "use the sigmoid derivative factor 2 * act * (1-act) in modulating learning -- otherwise just multiply by msn activation directly -- this is generally beneficial for learning to prevent weights from continuing to increase when activations are already strong (and vice-versa for decreases)"}, {Name: "Decay", Doc: "multiplier on trace activation for decaying prior traces -- new trace magnitude drives decay of prior trace -- if gating activation is low, then new trace can be low and decay is slow, so increasing this factor causes learning to be more targeted on recent gating changes"}, {Name: "GateLRScale", Doc: "learning rate scale factor, if"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.DelInhState", IDName: "del-inh-state", Doc: "DelInhState contains extra variables for MSNLayer neurons -- stored separately", Fields: []gti.Field{{Name: "GePrvQ", Doc: "netin from previous quarter, used for delayed inhibition"}, {Name: "GePrvTrl", Doc: "netin from previous \"trial\" (alpha cycle), used for delayed inhibition"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.TraceSyn", IDName: "trace-syn", Doc: "TraceSyn holds extra synaptic state for trace projections", Fields: []gti.Field{{Name: "NTr", Doc: "new trace -- drives updates to trace value -- su * (1-ru_msn) for gated, or su * ru_msn for not-gated (or for non-thalamic cases)"}, {Name: "Tr", Doc: " current ongoing trace of activations, which drive learning -- adds ntr and clears after learning on current values -- includes both thal gated (+ and other nongated, - inputs)"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.DALrnRule", IDName: "da-lrn-rule"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.MSNPrjn", IDName: "msn-prjn", Doc: "MSNPrjn does dopamine-modulated, for striatum-like layers", Embeds: []gti.Field{{Name: "Prjn"}}, Fields: []gti.Field{{Name: "LearningRule"}, {Name: "Trace", Doc: "special parameters for striatum trace learning"}, {Name: "TrSyns", Doc: "trace synaptic state values, ordered by the sending layer units which owns them -- one-to-one with SConIndex array"}, {Name: "SLActVar", Doc: "sending layer activation variable name"}, {Name: "RLActVar", Doc: "receiving layer activation variable name"}, {Name: "MaxVSActMod", Doc: "for VS matrix TRACE_NO_THAL_VS and DA_HEBB_VS learning rules, this is the maximum value that the deep_mod_net modulatory inputs from the basal amygdala (up state enabling signal) can contribute to learning"}, {Name: "DaMod", Doc: "parameters for dopaminergic modulation"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.IMSNPrjn", IDName: "imsn-prjn"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.Network", IDName: "network", Embeds: []gti.Field{{Name: "Network"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.INetwork", IDName: "i-network"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.ModNeuronVar", IDName: "mod-neuron-var", Doc: "NeuronVars are indexes into extra neuron-level variables"})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.PPTgLayer", IDName: "pp-tg-layer", Doc: "The PPTg passes on a positively-rectified version of its input signal.", Embeds: []gti.Field{{Name: "Layer"}}, Fields: []gti.Field{{Name: "Ge"}, {Name: "GePrev"}, {Name: "SendAct"}, {Name: "DA"}, {Name: "DNetGain", Doc: "gain on input activation"}, {Name: "ActThreshold", Doc: "activation threshold for passing through"}, {Name: "ClampActivation", Doc: "clamp activation directly, after applying gain"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.PVLayer", IDName: "pv-layer", Doc: "Primary Value input layer. Sends activation directly to its receivers, bypassing the standard mechanisms.", Embeds: []gti.Field{{Name: "Layer"}}, Fields: []gti.Field{{Name: "Net"}, {Name: "SendPVQuarter"}, {Name: "PVReceivers"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.VTADAGains", IDName: "vtada-gains", Doc: "Gain constants for inputs to the VTA", Fields: []gti.Field{{Name: "DA", Doc: "overall multiplier for dopamine values"}, {Name: "PPTg", Doc: "gain on bursts from PPTg"}, {Name: "LHb", Doc: "gain on dips/bursts from LHbRMTg"}, {Name: "PV", Doc: "gain on positive PV component of total phasic DA signal (net after subtracting VSPatchIndir (PVi) shunt signal)"}, {Name: "PVIBurstShunt", Doc: "gain on VSPatch projection that shunts bursting in VTA (for VTAp = VSPatchPosD1, for VTAn = VSPatchNegD2)"}, {Name: "PVIAntiBurstShunt", Doc: "gain on VSPatch projection that opposes shunting of bursting in VTA (for VTAp = VSPatchPosD2, for VTAn = VSPatchNegD1)"}, {Name: "PVIDipShunt", Doc: "gain on VSPatch projection that shunts dipping of VTA (currently only VTAp supported = VSPatchNegD2) -- optional and somewhat controversial"}, {Name: "PVIAntiDipShunt", Doc: "gain on VSPatch projection that opposes the shunting of dipping in VTA (currently only VTAp supported = VSPatchNegD1)"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.VTALayer", IDName: "vta-layer", Doc: "VTA internal state", Embeds: []gti.Field{{Name: "ClampDaLayer"}}, Fields: []gti.Field{{Name: "SendVal"}, {Name: "Valence", Doc: "VTA layer DA valence, positive or negative"}, {Name: "TonicDA", Doc: "set a tonic 'dopamine' (DA) level (offset to add to da values)"}, {Name: "DAGains", Doc: "gains for various VTA inputs"}, {Name: "RecvFrom"}, {Name: "InternalState", Doc: "input values--for debugging only"}}})

var _ = gti.AddType(&gti.Type{Name: "github.com/emer/leabra/v2/pvlv.VTAState", IDName: "vta-state", Doc: "monitoring and debugging only. Received values from all inputs", Fields: []gti.Field{{Name: "PPTgDAp"}, {Name: "LHbDA"}, {Name: "PosPVAct"}, {Name: "VSPosPVI"}, {Name: "VSNegPVI"}, {Name: "BurstLHbDA"}, {Name: "DipLHbDA"}, {Name: "TotBurstDA"}, {Name: "TotDipDA"}, {Name: "NetDipDA"}, {Name: "NetDA"}, {Name: "SendVal"}}})
