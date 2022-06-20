// Code generated by "stringer -type=PlayType"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Single-0]
	_ = x[Double-1]
	_ = x[Triple-2]
	_ = x[HomeRun-3]
	_ = x[CaughtStealing-4]
	_ = x[HitByPitch-5]
	_ = x[Walk-6]
	_ = x[WalkWildPitch-7]
	_ = x[WalkPassedBall-8]
	_ = x[StolenBase-9]
	_ = x[PickedOff-10]
	_ = x[CatcherInterference-11]
	_ = x[ReachedOnError-12]
	_ = x[FieldersChoice-13]
	_ = x[WildPitch-14]
	_ = x[PassedBall-15]
	_ = x[GroundOut-16]
	_ = x[FlyOut-17]
	_ = x[DoublePlay-18]
	_ = x[TriplePlay-19]
	_ = x[StrikeOut-20]
	_ = x[StrikeOutPassedBall-21]
	_ = x[StrikeOutWildPitch-22]
	_ = x[StrikeOutPickedOff-23]
	_ = x[FoulFlyError-24]
	_ = x[NoPlay-25]
}

const _PlayType_name = "SingleDoubleTripleHomeRunCaughtStealingHitByPitchWalkWalkWildPitchWalkPassedBallStolenBasePickedOffCatcherInterferenceReachedOnErrorFieldersChoiceWildPitchPassedBallGroundOutFlyOutDoublePlayTriplePlayStrikeOutStrikeOutPassedBallStrikeOutWildPitchStrikeOutPickedOffFoulFlyErrorNoPlay"

var _PlayType_index = [...]uint16{0, 6, 12, 18, 25, 39, 49, 53, 66, 80, 90, 99, 118, 132, 146, 155, 165, 174, 180, 190, 200, 209, 228, 246, 264, 276, 282}

func (i PlayType) String() string {
	if i >= PlayType(len(_PlayType_index)-1) {
		return "PlayType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PlayType_name[_PlayType_index[i]:_PlayType_index[i+1]]
}
