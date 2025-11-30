package client

import "sort"

const InterpDelayMs = 200
const MaxSnapshots = 64

type Interpolator struct {
	Snapshots []WorldSnapshot
}

func NewInterpolator() *Interpolator {
	return &Interpolator{Snapshots: make([]WorldSnapshot, 0, MaxSnapshots)}
}

func (it *Interpolator) AddSnapshot(s WorldSnapshot) {
	if len(it.Snapshots) >= MaxSnapshots {
		copy(it.Snapshots, it.Snapshots[1:])
		it.Snapshots = it.Snapshots[:MaxSnapshots-1]
	}
	it.Snapshots = append(it.Snapshots, s)
	sort.Slice(it.Snapshots, func(i, j int) bool {
		return it.Snapshots[i].Timestamp < it.Snapshots[j].Timestamp
	})
}

func (it *Interpolator) GetRenderState() WorldSnapshot {
	if len(it.Snapshots) == 0 {
		return WorldSnapshot{}
	}
	target := NowMs() - InterpDelayMs

	if len(it.Snapshots) == 1 {
		return it.Snapshots[0]
	}

	var older, newer *WorldSnapshot
	for i := len(it.Snapshots) - 1; i >= 0; i-- {
		if it.Snapshots[i].Timestamp <= target {
			older = &it.Snapshots[i]
			if i+1 < len(it.Snapshots) {
				newer = &it.Snapshots[i+1]
			} else {
				newer = &it.Snapshots[i]
			}
			break
		}
	}

	if older == nil {
		return it.Snapshots[0]
	}

	if newer == nil || older.Timestamp == newer.Timestamp {
		return *older
	}

	t0 := float32(older.Timestamp)
	t1 := float32(newer.Timestamp)
	tf := float32(target)
	var f float32
	if t1 > t0 {
		f = (tf - t0) / (t1 - t0)
	} else {
		f = 0
	}

	out := *older
	out.Timestamp = target

	for i := range out.Players {
		id := out.Players[i].ID
		for _, np := range newer.Players {
			if np.ID == id {
				out.Players[i].X = out.Players[i].X + (np.X-out.Players[i].X)*f
				out.Players[i].Y = out.Players[i].Y + (np.Y-out.Players[i].Y)*f
				out.Players[i].Score = np.Score
			}
		}
	}

	return out
}
