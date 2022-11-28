package vo

type Record[T any] struct {
	Label string `json:"label"`
	Value T      `json:"value"`
}
