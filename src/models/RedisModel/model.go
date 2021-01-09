package RedisModel


type RedisResponse struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	ExpTime string `json:"exp"`
	Desc    string `json:"desc"`
}

func NewRedisResponse(key string, value string, expTime string) *RedisResponse {

	desc := value
	if len(value) >= 50 {
		valueRune := []rune(value)
		desc = string(valueRune[:50]) + "..."
	} else {
		desc = value
	}

	return &RedisResponse{Key: key, Value: value, ExpTime: expTime, Desc: desc}
}
