package client

const InterpDelayMs = 200 // matches server artificial latency
const MaxSnapshots = 64

type Interpolator struct {
	Snapshots []WorldSnapshot
}

func NewInterpolator() *Interpolator {
	return &Interpolator{
		Snapshots: make([]WorldSnapshot, 0, MaxSnapshots),
	}
}

func (i *Interpolator) AddSnapshot(s WorldSnapshot) {
	if len(i.Snapshots) >= MaxSnapshots {
		copy(i.Snapshots, i.Snapshots[1:])
		i.Snapshots = i.Snapshots[:MaxSnapshots-1]
	}

	i.Snapshots = append(i.Snapshots, s)
}

func (i *Interpolator) GetRenderState() WorldSnapshot {
	targetTime := NowMs() - InterpDelayMs

	if len(i.Snapshots) < 2 {
		// not enough information yet
		return WorldSnapshot{}
	}

	// find snapshots surrounding targetTime
	var older, newer WorldSnapshot
	found := false

	for idx := len(i.Snapshots) - 1; idx >= 0; idx-- {
		if i.Snapshots[idx].Timestamp <= targetTime {
			older = i.Snapshots[idx]

			if idx < len(i.Snapshots)-1 {
				newer = i.Snapshots[idx+1]
			} else {
				newer = older
			}

			found = true
			break
		}
	}

	if !found {
		// all snapshots newer — use the oldest
		return i.Snapshots[0]
	}

	// interpolate factor
	t0 := float32(older.Timestamp)
	t1 := float32(newer.Timestamp)
	tf := float32(targetTime)

	var factor float32 = 0
	if t1 > t0 {
		factor = (tf - t0) / (t1 - t0)
	}

	// interpolate players
	result := older
	result.Timestamp = targetTime

	for pIdx := range result.Players {
		id := result.Players[pIdx].ID

		// find same player in newer
		for _, np := range newer.Players {
			if np.ID == id {
				ox, oy := result.Players[pIdx].X, result.Players[pIdx].Y
				nx, ny := np.X, np.Y

				result.Players[pIdx].X = ox + (nx-ox)*factor
				result.Players[pIdx].Y = oy + (ny-oy)*factor
			}
		}
	}

	// coins are static until next tick → no interpolation needed

	return result
}
