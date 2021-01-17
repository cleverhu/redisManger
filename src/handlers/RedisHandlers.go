package handlers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"math"
	"redisManger/src/common"
	"redisManger/src/dbs"
	"redisManger/src/models/InfoModel"
	"redisManger/src/models/LogModel"
	"redisManger/src/models/LoginUserModel"
	"redisManger/src/models/RedisConfigModel"
	"redisManger/src/models/RedisModel"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rds *redis.Client
var orm *gorm.DB

func init() {
	rds = dbs.Rds
	orm = dbs.Orm
}

func Scan(ctx *gin.Context) {

	//rds.SetBit()
	method := strings.ToLower(ctx.DefaultQuery("method", ""))
	if method != "string" && method != "list" && method != "hash" && method != "set" && method != "zset" && method != "bitmap" {
		ctx.JSON(400, gin.H{"message": "search error"})
		return
	}

	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "0")
	size := ctx.DefaultQuery("size", "20")

	cursorInt, _ := strconv.ParseInt(cursor, 10, 64)
	sizeInt, _ := strconv.ParseInt(size, 10, 64)

	var stringRes []*RedisModel.StringModel
	var commonRes []*RedisModel.CommonModel
LABEL:
	scan := rds.Scan(uint64(cursorInt), "*"+search+"*", sizeInt)

	keys, nextCursor := scan.Val()

	m := &sync.Map{}
	for _, v := range keys {
		m.Store(v, "1")
	}

	m.Range(func(key, value interface{}) bool {
		exp := rds.TTL(key.(string)).Val()
		expTime := ""
		if exp.Seconds() < 0 {
			expTime = "-"
		} else {
			expTime = fmt.Sprintf("%.0f", exp.Seconds())
		}

		t := rds.Type(key.(string)).Val()
		if t == method {
			switch t {
			case "string":
				//v := rds.Get(key.(string)).Val()
				//fmt.Println(strings.TrimSpace(v))
				//fmt.Println(strings.TrimSpace(v) == "")//可以疑似区分bitmap
				r := RedisModel.NewStringModel(key.(string), rds.Get(key.(string)).Val(), expTime)
				stringRes = append(stringRes, r)
			case "list":
				c := &RedisModel.CommonModel{
					Key:     key.(string),
					Length:  rds.LLen(key.(string)).Val(),
					ExpTime: expTime,
				}
				commonRes = append(commonRes, c)
			case "set":
				c := &RedisModel.CommonModel{
					Key:     key.(string),
					Length:  rds.SCard(key.(string)).Val(),
					ExpTime: expTime,
				}
				commonRes = append(commonRes, c)
			case "hash":
				c := &RedisModel.CommonModel{
					Key:     key.(string),
					Length:  rds.HLen(key.(string)).Val(),
					ExpTime: expTime,
				}
				commonRes = append(commonRes, c)
			case "zset":
				c := &RedisModel.CommonModel{
					Key:     key.(string),
					Length:  rds.ZCard(key.(string)).Val(),
					ExpTime: expTime,
				}
				commonRes = append(commonRes, c)
			case "geo":
				c := &RedisModel.CommonModel{
					Key:     key.(string),
					ExpTime: expTime,
				}
				commonRes = append(commonRes, c)
			}
		}
		return true
	})
	if stringRes == nil && commonRes == nil && nextCursor != 0 {
		cursorInt = int64(nextCursor)
		goto LABEL
	}

	if method == "string" {
		ctx.JSON(200, gin.H{"data": stringRes, "next": nextCursor, "current": cursor})
	} else {
		ctx.JSON(200, gin.H{"data": commonRes, "next": nextCursor, "current": cursor})
	}

}

func DelByKeys(ctx *gin.Context) {
	keys := ctx.Param("keys")
	s := strings.Split(keys, ":")
	for i := 0; i < len(s); i++ {
		rds.Del(s[i])
	}
	ctx.JSON(200, gin.H{"message": "删除成功"})
}

func StringUpdate(ctx *gin.Context) {
	type form struct {
		Key   string  `json:"key"`
		Value string  `json:"value"`
		Exp   float64 `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	t := rds.Type(f.Key).Val()

	if t != "none" && t != "string" {
		ctx.JSON(400, gin.H{"message": "修改失败key已经存在且不是string类型无法添加"})
		return
	}

	err = rds.Set(f.Key, f.Value, time.Duration(int64(f.Exp)*int64(math.Pow10(9)))).Err()
	if err != nil {
		ctx.JSON(400, gin.H{"message": "赋值失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "添加成功"})

}

func ToFile(ctx *gin.Context) {

	m := sync.Map{}
	var cursor uint64
	var result = make([]string, 0)
	var err = errors.New("")
	for true {

		result, cursor, err = rds.Scan(cursor, "*", 100).Result()
		fmt.Println(err)
		if err != nil {
			ctx.JSON(400, gin.H{"message": "导出失败"})
			return
		}

		for _, v := range result {
			if rds.Type(v).Val() == "string" {
				m.Store(v, rds.Get(v).Val())
			}
		}

		if cursor == 0 {
			break
		}
	}
	token := ctx.Request.Header.Get("Token")
	_, claims, _ := common.ParseToken(token)
	f := excelize.NewFile()
	_ = f.SetCellValue("Sheet1", "A1", "key")
	_ = f.SetCellValue("Sheet1", "B1", "value")
	var i = 2
	m.Range(func(key, value interface{}) bool {
		_ = f.SetCellValue("Sheet1", "A"+fmt.Sprintf("%d", i), key)
		_ = f.SetCellValue("Sheet1", "B"+fmt.Sprintf("%d", i), value)
		i++
		return true
	})
	fileName := claims.UserName + "_" + time.Now().Format("2006-01-02-15:04:05导出string记录") + ".xlsx"
	err = f.SaveAs("./files/" + fileName)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "导出成功", "data": gin.H{"url": "http://m.deeplythink.com/excels/" + fileName}})
}

func CommonToFile(ctx *gin.Context) {
	type form struct {
		KeyAxis   string `json:"keyAxis"`
		KeyName   string `json:"keyName"`
		ValueAxis string `json:"valueAxis"`
		ValueName string `json:"valueName"`
		SheetName string `json:"sheetName"`
	}
	myForm := &form{}
	err := ctx.ShouldBindJSON(myForm)
	fmt.Println(myForm)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		fmt.Println(err)
		return
	}

	type cell struct {
		name  string
		index int64
	}
	k := &cell{
		name:  "",
		index: 0,
	}
	v := &cell{
		name:  "",
		index: 0,
	}

	c := regexp.MustCompile(`([A-Z]+)([0-9]+)`)

	if c.MatchString(myForm.KeyAxis) {
		matches := c.FindStringSubmatch(myForm.KeyAxis)
		k.name = matches[1]
		k.index, _ = strconv.ParseInt(matches[2], 10, 64)
		if k.index < 1 {
			k.index = 1
		}
	} else {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		return
	}

	if c.MatchString(myForm.ValueAxis) {
		matches := c.FindStringSubmatch(myForm.ValueAxis)
		v.name = matches[1]
		v.index, _ = strconv.ParseInt(matches[2], 10, 64)
		if v.index < 1 {
			v.index = 1
		}
	} else {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		return
	}

	m := sync.Map{}
	var cursor uint64
	var result = make([]string, 0)

	for true {

		result, cursor, err = rds.Scan(cursor, "*", 100).Result()
		fmt.Println(err)
		if err != nil {
			ctx.JSON(400, gin.H{"message": "导出失败"})
			return
		}

		for _, v := range result {
			if rds.Type(v).Val() == "string" {
				m.Store(v, rds.Get(v).Val())
			}
		}

		if cursor == 0 {
			break
		}
	}

	token := ctx.Request.Header.Get("Token")
	_, claims, _ := common.ParseToken(token)
	f := excelize.NewFile()
	_ = f.SetCellValue(myForm.SheetName, k.name+fmt.Sprintf("%d", k.index), myForm.KeyName)
	_ = f.SetCellValue(myForm.SheetName, v.name+fmt.Sprintf("%d", v.index), myForm.ValueName)

	m.Range(func(key, value interface{}) bool {
		k.index++
		v.index++
		_ = f.SetCellValue(myForm.SheetName, k.name+fmt.Sprintf("%d", k.index), key)
		_ = f.SetCellValue(myForm.SheetName, v.name+fmt.Sprintf("%d", v.index), value)

		return true
	})
	fileName := claims.UserName + "_" + time.Now().Format("2006-01-02-15:04:05导出string记录") + ".xlsx"
	err = f.SaveAs("./files/" + fileName)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "导出成功", "data": gin.H{"url": "http://redis.deeplythink.com/files/" + fileName}})
}

func CommonInsertToMySql(ctx *gin.Context) {
	type form struct {
		Dns           string `json:"dns"`
		TableName     string `json:"tableName"`
		ColumnOfKey   string `json:"columnOfKey"`
		ColumnOfValue string `json:"columnOfValue"`
	}
	myForm := &form{}
	err := ctx.ShouldBindJSON(myForm)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败,请检查参数"})
		return
	}

	db, err := gorm.Open("mysql", myForm.Dns)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败,连接数据库失败"})
		return
	}
	db.LogMode(true)
	m := sync.Map{}
	var cursor uint64
	var result = make([]string, 0)

	for true {
		result, cursor, err = rds.Scan(cursor, "*", 100).Result()
		fmt.Println(err)
		if err != nil {
			ctx.JSON(400, gin.H{"message": "导出失败"})
			return
		}

		for _, v := range result {
			if rds.Type(v).Val() == "string" {
				m.Store(v, rds.Get(v).Val())
			}
		}

		if cursor == 0 {
			break
		}
	}

	m.Range(func(key, value interface{}) bool {
		err = db.Table(myForm.TableName).Exec("insert into "+myForm.TableName+" ("+myForm.ColumnOfKey+","+myForm.ColumnOfValue+") values (?,?)", key, value).Error
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		ctx.JSON(400, gin.H{"message": "导出失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "导出成功"})
}

func ListGetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	lLen := rds.LLen(key).Val()
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	if pageInt < 1 {
		pageInt = 1
	}
	size := ctx.DefaultQuery("size", "20")
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	if sizeInt < 1 {
		sizeInt = 1
	}
	data := rds.LRange(key, (pageInt-1)*sizeInt, pageInt*sizeInt+-1).Val()
	var vms []*RedisModel.ListValueModel
	for i, v := range data {
		rm := &RedisModel.ListValueModel{
			Key:   key,
			Value: v,
			Index: (pageInt-1)*sizeInt + int64(i),
		}
		vms = append(vms, rm)
	}
	ctx.JSON(200, gin.H{"data": vms, "len": lLen, "key": key})
}

func ListPost(ctx *gin.Context) {
	type form struct {
		Key    string  `json:"key"`
		Value  string  `json:"value"`
		Exp    float64 `json:"exp"`
		Exists bool    `json:"exist"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	if f.Exists {
		rds.Del(f.Key)
	}

	err = rds.LPush(f.Key, f.Value).Err()

	if err != nil {
		ctx.JSON(400, gin.H{"message": "添加失败"})
		return
	}
	if f.Exp > 0 {
		rds.Expire(f.Key, time.Duration(int64(f.Exp)*int64(math.Pow10(9))))
	}

	ctx.JSON(200, gin.H{"message": "添加成功"})
}

func ListExists(ctx *gin.Context) {
	key := ctx.Param("key")
	keys := rds.Keys(key).Val()
	if len(keys) < 1 {
		ctx.JSON(200, gin.H{"result": 1, "message": "可以创建"})
		return
	} else {
		lLen := rds.LLen(key).Val()
		if lLen == 0 {
			ctx.JSON(200, gin.H{"result": -1, "message": "此键已经存在，并且不是list类型，无法添加。"})
			return
		}
		ctx.JSON(200, gin.H{"result": 0, "message": "此键已经存在，是否选择覆盖？"})
		return
	}
}

func ListRemoveValue(ctx *gin.Context) {
	m := RedisModel.ListValueModel{}
	err := ctx.ShouldBindJSON(&m)
	fmt.Println(m)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "input error"})
		return
	}
	sum := md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))

	err = rds.LSet(m.Key, m.Index, fmt.Sprintf("%x", sum)).Err()
	fmt.Println(err)
	err = rds.LRem(m.Key, 1, fmt.Sprintf("%x", sum)).Err()
	fmt.Println(err)
	ctx.JSON(200, gin.H{"message": "delete success"})
}

func ListInsert(ctx *gin.Context) {
	type form struct {
		Method int    `json:"method"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	if err != nil {
		ctx.JSON(400, "输入错误")
		return
	}

	if f.Method == 0 {
		err = rds.LPush(f.Key, f.Value).Err()
	} else {
		err = rds.RPush(f.Key, f.Value).Err()
	}
	if err != nil {
		ctx.JSON(400, gin.H{"message": "插入错误"})
		return
	}
	ctx.JSON(200, gin.H{"message": "插入成功"})
}

func HashGetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	lLen := rds.HLen(key).Val()
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	if pageInt < 1 {
		pageInt = 1
	}
	size := ctx.DefaultQuery("size", "5")
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	if sizeInt < 1 {
		sizeInt = 1
	}

	var result []*RedisModel.HashValueModel
	res := make(map[string]string)
	var i int64
	if lLen >= (pageInt-1)*sizeInt {
		res = rds.HGetAll(key).Val()
		for field, value := range res {
			if i >= (pageInt-1)*sizeInt {
				s := &RedisModel.HashValueModel{
					Field: field,
					Key:   key,
					Value: value,
				}
				result = append(result, s)
			}
			i++
			if i >= pageInt*sizeInt {
				break
			}
		}
	}

	ctx.JSON(200, gin.H{"data": result, "total": lLen, "key": key})
}

func HashRemoveValue(ctx *gin.Context) {
	s := RedisModel.HashValueModel{}
	err := ctx.ShouldBindJSON(&s)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	err = rds.HDel(s.Key, s.Field).Err()
	fmt.Println(s)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "删除错误"})
		return
	}
	ctx.JSON(200, gin.H{"message": "删除成功"})
}

func HashPost(ctx *gin.Context) {
	type form struct {
		Key   string  `json:"key"`
		Field string  `json:"field"`
		Value string  `json:"value"`
		Force bool    `json:"force"`
		Exp   float64 `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	t := rds.Type(f.Key).Val()
	if t != "none" && t != "hash" {
		ctx.JSON(400, gin.H{"message": "此键存在且不是hash类型，无法添加或修改"})
		return
	}
	fmt.Println(f)
	if f.Force {
		err = rds.HSet(f.Key, f.Field, f.Value).Err()
	} else {
		b := rds.HSetNX(f.Key, f.Field, f.Value).Val()
		if !b {
			ctx.JSON(400, gin.H{"message": "已经存在无法修改"})
			return
		}
	}

	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "添加失败"})
		return
	}

	if f.Exp > 0 {
		rds.Expire(f.Key, time.Duration(int64(f.Exp)*int64(math.Pow10(9))))
	}

	ctx.JSON(200, gin.H{"message": "修改成功"})
}

func SetGetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	lLen := rds.SCard(key).Val()
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	var cursor uint64
	if pageInt < 1 {
		pageInt = 1
	}
	size := ctx.DefaultQuery("size", "5")
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	if sizeInt < 1 {
		sizeInt = 1
	}
	var i int64

	res := make([]string, 0)
	var length int64
	for true {
		i++
		if i <= pageInt {
			res, cursor, _ = rds.SScan(key, cursor, "*", sizeInt).Result()
			length += int64(len(res))
		} else {
			_, cursor, _ = rds.SScan(key, cursor, "*", sizeInt).Result()
		}

		//fmt.Print(len(res), "-", cursor, "|")
		if cursor == 0 {
			if length == lLen {
				i--
			}
			break
		}
	}
	//fmt.Println()
	var result []*RedisModel.SetValueModel
	for _, v := range res {
		s := &RedisModel.SetValueModel{
			Key:   key,
			Value: v,
		}
		result = append(result, s)
	}

	ctx.JSON(200, gin.H{"data": result, "total": lLen, "key": key, "total_page": i})
}

func SetRemoveValue(ctx *gin.Context) {
	s := RedisModel.SetValueModel{}
	err := ctx.ShouldBindJSON(&s)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	err = rds.SRem(s.Key, s.Value).Err()
	if err != nil {
		ctx.JSON(400, gin.H{"message": "删除错误"})
		return
	}
	ctx.JSON(200, gin.H{"message": "删除成功"})
}

func SetPost(ctx *gin.Context) {
	type form struct {
		Key   string  `json:"key"`
		Value string  `json:"value"`
		Exp   float64 `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	t := rds.Type(f.Key).Val()
	if t != "none" && t != "set" {
		ctx.JSON(400, gin.H{"message": "此键存在且不是set类型，无法添加"})
		return
	}
	err = rds.SAdd(f.Key, f.Value).Err()
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "添加失败"})
		return
	}
	if f.Exp > 0 {
		rds.Expire(f.Key, time.Duration(int64(f.Exp)*int64(math.Pow10(9))))
	}

	ctx.JSON(200, gin.H{"message": "添加成功"})
}

func SetGetCommon(ctx *gin.Context) {
	type form struct {
		Keys   []string `json:"keys"`
		Method int      `json:"method"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	result := make([]string, 0)

	if f.Method == 0 {
		result = rds.SInter(f.Keys...).Val()
	} else {
		result = rds.SUnion(f.Keys...).Val()
	}

	ctx.JSON(200, gin.H{"message": "查询成功", "data": result})
}

func ZSetGetByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	lLen := rds.ZCard(key).Val()
	page := ctx.DefaultQuery("page", "1")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	if pageInt < 1 {
		pageInt = 1
	}
	size := ctx.DefaultQuery("size", "5")
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	if sizeInt < 1 {
		sizeInt = 1
	}
	var result []*RedisModel.ZSetValueModel
	vals := rds.ZRangeWithScores(key, 0, -1).Val()
	//fmt.Println(vals)
	var index int64
	for _, v := range vals {
		if index >= (pageInt-1)*sizeInt {
			r := &RedisModel.ZSetValueModel{
				Key:    key,
				Member: v.Member.(string),
				Score:  v.Score,
			}
			result = append(result, r)
		}

		if len(result) >= int(sizeInt) {
			break
		}

		index++
	}
	ctx.JSON(200, gin.H{"data": result, "total": lLen, "key": key})
}

func ZSetPost(ctx *gin.Context) {
	type form struct {
		Key    string  `json:"key"`
		Member string  `json:"member"`
		Score  float64 `json:"score"`
		Exp    float64 `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	t := rds.Type(f.Key).Val()
	if t != "none" && t != "zset" {
		ctx.JSON(400, gin.H{"message": "此键存在且不是set类型，无法添加"})
		return
	}
	err = rds.ZAdd(f.Key, redis.Z{
		Score:  f.Score,
		Member: f.Member,
	}).Err()
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "添加失败"})
		return
	}

	if f.Exp > 0 {
		rds.Expire(f.Key, time.Duration(int64(f.Exp)*int64(math.Pow10(9))))
	}

	ctx.JSON(200, gin.H{"message": "添加成功"})
}

func ZSetRemoveValue(ctx *gin.Context) {
	zs := RedisModel.ZSetValueModel{}
	err := ctx.ShouldBindJSON(&zs)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	err = rds.ZRem(zs.Key, zs.Member).Err()
	if err != nil {
		ctx.JSON(400, gin.H{"message": "删除错误"})
		return
	}
	ctx.JSON(200, gin.H{"message": "删除成功"})
}

func GEOPost(ctx *gin.Context) {
	type form struct {
		Key       string  `json:"key"`
		Member    string  `json:"member"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Exp       float64 `json:"exp"`
	}
	f := &form{}
	err := ctx.ShouldBindJSON(&f)
	fmt.Println(err, f)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	t := rds.Type(f.Key).Val()
	if t != "none" && t != "geo" {
		ctx.JSON(400, gin.H{"message": "此键存在且不是geo类型，无法添加"})
		return
	}
	err = rds.GeoAdd(f.Key, &redis.GeoLocation{
		Name:      f.Member,
		Longitude: f.Longitude,
		Latitude:  f.Latitude,
	}).Err()

	if err != nil {
		ctx.JSON(400, gin.H{"message": "添加失败"})
		return
	}

	dbs.Orm.Table("geo").Exec("insert into geo (geo_key,geo_member) values(?,?)", f.Key, f.Member)

	if f.Exp > 0 {
		rds.Expire(f.Key, time.Duration(int64(f.Exp)*int64(math.Pow10(9))))
	}

	ctx.JSON(200, gin.H{"message": "添加成功"})
}

func Config(ctx *gin.Context) {

	values := rds.ConfigGet("*").Val()
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

func UpdateConfig(ctx *gin.Context) {
	r := RedisConfigModel.NewUpdateRedisRequest()
	err := ctx.ShouldBindJSON(r)
	if err != nil {
		ctx.JSON(200, gin.H{"message": "输入错误"})
		return
	}
	err = rds.ConfigSet(r.Key, r.Value).Err()
	if err != nil {
		ctx.JSON(200, gin.H{"message": fmt.Sprintf("update config error:%v", err)})
		return
	}

	if r.EditFile {
		err := rds.ConfigRewrite()
		if err != nil {
			ctx.JSON(200, gin.H{"message": fmt.Sprintf("update config success, but write to config file error:%v", err)})
			return
		}

	}

	ctx.JSON(200, gin.H{"message": "修改成功"})

}

func Info(ctx *gin.Context) {
	str := rds.Info().Val()
	ret := InfoModel.GetInfo(str)
	ctx.JSON(200, ret)
}

func Login(ctx *gin.Context) {
	user := LoginUserModel.LoginUser{}
	err := ctx.ShouldBindJSON(&user)
	fmt.Println(user, err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入信息有误"})
		return
	}
	orm.Table("users").Raw("select * from users where password = ? and username = ?", user.Username, user.Password).First(&user)
	if user.ID == 0 {
		ctx.JSON(400, gin.H{"message": "账号或者密码有错"})
		return
	}
	orm.Table("users").Exec("update users set update_time = ? where id = ?", time.Now(), user.ID)
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "token发放失败"})
		return
	}

	orm.Table("user_log").Exec("insert into user_log (uid,aid,path,method,logtime) values (?,?,?,?,?)", user.ID, common.USER, ctx.Request.URL.Path, ctx.Request.Method, time.Now())
	ctx.JSON(200, gin.H{"message": "登录成功", "token": token})
}

func Validate(ctx *gin.Context) {
	type form struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	f := &form{}
	_ = ctx.ShouldBindJSON(&f)
	_, claims, err := common.ParseToken(f.Token)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "不合法的token"})
		return
	}
	f.Username = claims.UserName

	ctx.JSON(200, f)
}

func ConnectTest(ctx *gin.Context) {
	type form struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}
	f := &form{}

	err := ctx.ShouldBindJSON(&f)
	fmt.Println(f)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "连接失败"})
		return
	}
	//
	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", f.Host, f.Port), Password: f.Password, DB: f.DB})
	defer client.Close()
	err = client.Ping().Err()
	fmt.Println(err)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "连接失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "连接成功"})
}

func ConnectSave(ctx *gin.Context) {
	type form struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}
	f := &form{}

	err := ctx.ShouldBindJSON(&f)
	//fmt.Println(f)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	token := ctx.Request.Header.Get("Token")

	_, claims, _ := common.ParseToken(token)
	fmt.Println(claims)
	id := claims.UserId
	fmt.Println(id)
	orm.Table("user_redis_info").Exec("delete from user_redis_info  where uid = ?", id)
	err = orm.Table("user_redis_info").Exec("insert into user_redis_info (uid,host,port,password,db) values (?,?,?,?,?)", id, f.Host, f.Port, f.Password, f.DB).Error
	if err != nil {
		ctx.JSON(400, gin.H{"message": "保存失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "保存成功"})
}

func ConnectGet(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Token")
	_, claims, _ := common.ParseToken(token)
	id := claims.UserId
	type form struct {
		Host     string `json:"host" gorm:"column:host"`
		Port     int    `json:"port" gorm:"port"`
		Password string `json:"password" gorm:"password"`
		DB       int    `json:"db" gorm:"db"`
	}
	f := &form{}
	orm.Table("user_redis_info").Select("host,port,password,db").Where("uid=?", id).First(&f)
	ctx.JSON(200, gin.H{"message": "get success", "data": f})
}

func LogsGet(ctx *gin.Context) {
	type form struct {
		Page   int    `json:"page"`
		Size   int    `json:"size"`
		Search string `json:"search"`
	}

	f := &form{}
	err := ctx.ShouldBindJSON(&f)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "输入错误"})
		return
	}
	fmt.Println(f)
	type result struct {
		Logs  []*LogModel.LogImpl `json:"logs"`
		Total int64               `json:"total"`
	}
	total := &struct {
		Total int64 `json:"total"`
	}{}
	var logs []*LogModel.LogImpl

	token := ctx.Request.Header.Get("Token")
	_, claims, _ := common.ParseToken(token)

	res := &result{}
	orm.Raw("select count(*) as total FROM user_log LEFT JOIN actions on user_log.aid = actions.id where uid = ? and (action like '%"+f.Search+"%' or path like '%"+f.Search+"%')", claims.UserId).Find(&total)
	res.Total = total.Total

	orm.Raw("select uid,action,path,method,logtime FROM user_log LEFT JOIN actions on user_log.aid = actions.id where uid = ? and (action like '%"+f.Search+"%' or path like '%"+f.Search+"%') order by logtime desc limit ?,? ", claims.UserId, (f.Page-1)*f.Size, f.Size).Find(&logs)
	res.Logs = logs
	ctx.JSON(200, gin.H{"message": "查询成功", "result": res})
}
