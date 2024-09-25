package timestamp

import "time"

type Timestamp int64 // UNIX time ms
type Duration int64  // ms

func (r Timestamp) Add(duration Duration) Timestamp {
	return r + Timestamp(duration)
}

func Now() Timestamp {
	return Timestamp(time.Now().UnixMilli())
}
