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

type ListModel struct {
	Key     string `json:"key"`
	Length  int64  `json:"len"`
	ExpTime string `json:"exp"`
}

func NewListModel(key string, length int64, expTime string) *ListModel {
	return &ListModel{Key: key, Length: length, ExpTime: expTime}
}

type ListValueModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Index int64  `json:"index"`
}

func NewListValueModel(key string, value string, index int64) *ListValueModel {
	return &ListValueModel{Key: key, Value: value, Index: index}
}
