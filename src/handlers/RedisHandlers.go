package handlers

import (
	"github.com/gin-gonic/gin"
	"redisManger/src/dbs"
	"redisManger/src/models/RedisModel"
	"strconv"
	"sync"
)

func GET(ctx *gin.Context) {
	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "0")
	size := ctx.DefaultQuery("size", "100")

	cursorInt, _ := strconv.ParseInt(cursor, 10, 64)
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	scan := dbs.Rds.Scan(uint64(cursorInt), search+"*", sizeInt)

	keys, nextCursor := scan.Val()
	m := &sync.Map{}
	for _, v := range keys {
		m.Store(v, dbs.Rds.Get(v).Val())
	}
	var res []*RedisModel.RedisResponse

	m.Range(func(key, value interface{}) bool {
		r := RedisModel.NewRedisResponse(key.(string), value.(string))
		res = append(res, r)
		return true
	})

	ctx.JSON(200, gin.H{"data": res, "nextCursor": nextCursor})
}
