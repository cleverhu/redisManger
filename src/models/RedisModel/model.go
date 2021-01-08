package RedisModel

type RedisResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewRedisResponse(key string, value string) *RedisResponse {
	return &RedisResponse{Key: key, Value: value}
}
