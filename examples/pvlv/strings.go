// Code generated by "stringer -linecomment -output=strings.go -type=StepGrain,StopStepCond"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Cycle-0]
	_ = x[Quarter-1]
	_ = x[SettleMinus-2]
	_ = x[SettlePlus-3]
	_ = x[AlphaCycle-4]
	_ = x[SGTrial-5]
	_ = x[Epoch-6]
	_ = x[MultiRunSequence-7]
	_ = x[StepGrainN-8]
}

const _StepGrain_name = "CycleQuarterSettleMinusSettlePlusAlphaCycleTrialEpochMultiRunSequenceStepGrainN"

var _StepGrain_index = [...]uint8{0, 5, 12, 23, 33, 43, 48, 53, 69, 79}

func (i StepGrain) String() string {
	if i < 0 || i >= StepGrain(len(_StepGrain_index)-1) {
		return "StepGrain(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _StepGrain_name[_StepGrain_index[i]:_StepGrain_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SSNone-0]
	_ = x[SSError-1]
	_ = x[SSCorrect-2]
	_ = x[SSTrialNameMatch-3]
	_ = x[SSTrialNameNonmatch-4]
	_ = x[StopStepCondN-5]
}

const _StopStepCond_name = "SSNoneErrorCorrectTrial NameNot Trial NameStopStepCondN"

var _StopStepCond_index = [...]uint8{0, 6, 11, 18, 28, 42, 55}

func (i StopStepCond) String() string {
	if i < 0 || i >= StopStepCond(len(_StopStepCond_index)-1) {
		return "StopStepCond(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _StopStepCond_name[_StopStepCond_index[i]:_StopStepCond_index[i+1]]
}