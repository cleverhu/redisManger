package RedisModel

type StringModel struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	ExpTime string `json:"exp"`
	Desc    string `json:"desc"`
}

func NewStringModel(key string, value string, expTime string) *StringModel {

	desc := value
	if len(value) >= 50 {
		valueRune := []rune(value)
		desc = string(valueRune[:50]) + "..."
	} else {
		desc = value
	}

	return &StringModel{Key: key, Value: value, ExpTime: expTime, Desc: desc}
}

type CommonModel struct {
	Key     string `json:"key"`
	Length  int64  `json:"len"`
	ExpTime string `json:"exp"`
}

type GEOModel struct {
	Key       string `json:"key"`
	Longitude float64  `json:"Longitude"`
	Latitude  float64  `json:"Latitude"`
	ExpTime   string `json:"exp"`
}

type ListValueModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index int64  `json:"index"`
}

type SetValueModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HashValueModel struct {
	Key   string `json:"key"`
	Field string `json:"field"`
	Value string `json:"value"`
}

type ZSetValueModel struct {
	Key    string  `json:"key"`
	Member string  `json:"member"`
	Score  float64 `json:"score"`
}
