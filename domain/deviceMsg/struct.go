package deviceMsg

type TimeValue struct {
	Timestamp int64 `json:"timestamp,omitempty"` //毫秒时间戳
	Value     any   `json:"value"`               //值
}
