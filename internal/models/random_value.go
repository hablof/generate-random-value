package models

type RandomValue struct {
	ID        uint64      `json:"id"`
	Value     string      `json:"value"`
	RequestID MaybeString `json:"-"`
}

type MaybeString struct {
	S       string
	IsValid bool
}

// func (rv *RandomValue) MarshalJSON() ([]byte, error) {
// 	if
// }
