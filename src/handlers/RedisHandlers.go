package handlers

import (
	"github.com/gin-gonic/gin"
	"redisManger/src/dbs"
)

func GetTop20(ctx *gin.Context) {
	scan := dbs.Rds.Scan(0, "", 20)
	keys, _ := scan.Val()
	//fmt.Println(cursor)
	//fmt.Println(len(keys))
	//for _, v := range keys {
	//	fmt.Println(v)
	//}
	ctx.JSON(200, keys)
}
