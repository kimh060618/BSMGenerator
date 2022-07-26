package bsm

type BSMMessage struct {
	Type      int   `json:"type"`
	TimeStamp int64 `json:"time"`
	Level     int   `json:"level"`
}
