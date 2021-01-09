package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"redisManger/src/dbs"
	"redisManger/src/models/RedisConfigModel"
	"redisManger/src/models/RedisModel"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GET(ctx *gin.Context) {
	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "0")
	size := ctx.DefaultQuery("size", "20")

	cursorInt, _ := strconv.ParseInt(cursor, 10, 64)
	sizeInt, _ := strconv.ParseInt(size, 10, 64)

LABEL:
	scan := dbs.Rds.Scan(uint64(cursorInt), "*"+search+"*", sizeInt)

	keys, nextCursor := scan.Val()
	m := &sync.Map{}
	for _, v := range keys {
		m.Store(v, dbs.Rds.Get(v).Val())
	}
	var res []*RedisModel.RedisResponse

	m.Range(func(key, value interface{}) bool {
		exp := dbs.Rds.TTL(key.(string)).Val()
		//fmt.Println(exp)
		expTime := ""
		if exp.Seconds() < 0 {
			expTime = "-"
		} else {
			expTime = fmt.Sprintf("%.0f", exp.Seconds())
		}
		//value = decode.UnicodeToUTF8(value.(string))

		r := RedisModel.NewRedisResponse(key.(string), value.(string), expTime)
		res = append(res, r)
		return true
	})
	//fmt.Println(res, nextCursor)
	if res == nil && nextCursor != 0 {
		cursorInt = int64(nextCursor)
		goto LABEL
	}

	ctx.JSON(200, gin.H{"data": res, "next": nextCursor, "current": cursor})
}

func DEL(ctx *gin.Context) {
	id := ctx.Param("id")
	s := strings.Split(id, ":")
	for i := 0; i < len(s); i++ {
		dbs.Rds.Del(s[i])
	}
	ctx.JSON(200, gin.H{"message": "success"})
}

func ADD(ctx *gin.Context) {
	type form struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Exp   int    `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "fail"})
		return
	}
	err = dbs.Rds.Set(f.Key, f.Value, time.Duration(int64(f.Exp)*int64(math.Pow10(9)))).Err()
	if err != nil {
		ctx.JSON(400, gin.H{"message": "fail"})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})

}

func Config(ctx *gin.Context) {

	values := dbs.Rds.ConfigGet("*").Val()
	//data key value
	var keys []string
	var cfs []*RedisConfigModel.Config
	m := make(map[string]string, 0)
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		value := values[i+1].(string)
		keys = append(keys, key)
		m[key] = value
	}

	sort.Strings(keys)

	for _, k := range keys {
		c := RedisConfigModel.NewConfig(k, m[k])
		cfs = append(cfs, c)
		//}
	}

	ctx.JSON(200, gin.H{"config": cfs})
}

func UpdateConfig(ctx *gin.Context){
	ctx.ShouldBindJSON()
}
