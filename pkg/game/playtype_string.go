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
	_ = x[WalkPickedOff-9]
	_ = x[StolenBase-10]
	_ = x[PickedOff-11]
	_ = x[CatcherInterference-12]
	_ = x[ReachedOnError-13]
	_ = x[FieldersChoice-14]
	_ = x[WildPitch-15]
	_ = x[PassedBall-16]
	_ = x[GroundOut-17]
	_ = x[FlyOut-18]
	_ = x[DoublePlay-19]
	_ = x[TriplePlay-20]
	_ = x[StrikeOut-21]
	_ = x[StrikeOutPassedBall-22]
	_ = x[StrikeOutWildPitch-23]
	_ = x[StrikeOutPickedOff-24]
	_ = x[StrikeOutStolenBase-25]
	_ = x[FoulFlyError-26]
	_ = x[NoPlay-27]
}

const _PlayType_name = "SingleDoubleTripleHomeRunCaughtStealingHitByPitchWalkWalkWildPitchWalkPassedBallWalkPickedOffStolenBasePickedOffCatcherInterferenceReachedOnErrorFieldersChoiceWildPitchPassedBallGroundOutFlyOutDoublePlayTriplePlayStrikeOutStrikeOutPassedBallStrikeOutWildPitchStrikeOutPickedOffStrikeOutStolenBaseFoulFlyErrorNoPlay"

var _PlayType_index = [...]uint16{0, 6, 12, 18, 25, 39, 49, 53, 66, 80, 93, 103, 112, 131, 145, 159, 168, 178, 187, 193, 203, 213, 222, 241, 259, 277, 296, 308, 314}

func (i PlayType) String() string {
	if i >= PlayType(len(_PlayType_index)-1) {
		return "PlayType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PlayType_name[_PlayType_index[i]:_PlayType_index[i+1]]
}
