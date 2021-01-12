package handlers

import (
	"crypto/md5"
	"fmt"
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
	method := strings.ToLower(ctx.DefaultQuery("method", ""))
	if method != "string" && method != "list" && method != "set" {
		ctx.JSON(400, gin.H{"message": "search error"})
		return
	}

	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "0")
	size := ctx.DefaultQuery("size", "20")

	cursorInt, _ := strconv.ParseInt(cursor, 10, 64)
	sizeInt, _ := strconv.ParseInt(size, 10, 64)

	var stringRes []*RedisModel.StringModel
	var listRes []*RedisModel.ListModel
	var setRes []*RedisModel.SetModel
LABEL:
	scan := rds.Scan(uint64(cursorInt), "*"+search+"*", sizeInt)

	keys, nextCursor := scan.Val()

	m := &sync.Map{}
	for _, v := range keys {
		m.Store(v, "1")
	}

	m.Range(func(key, value interface{}) bool {
		t := rds.Type(key.(string)).Val()
		if t == method {
			switch t {
			case "string":
				exp := rds.TTL(key.(string)).Val()
				expTime := ""
				if exp.Seconds() < 0 {
					expTime = "-"
				} else {
					expTime = fmt.Sprintf("%.0f", exp.Seconds())
				}

				r := RedisModel.NewStringModel(key.(string), rds.Get(key.(string)).Val(), expTime)
				stringRes = append(stringRes, r)
			case "list":
				exp := rds.TTL(key.(string)).Val()
				expTime := ""
				if exp.Seconds() < 0 {
					expTime = "-"
				} else {
					expTime = fmt.Sprintf("%.0f", exp.Seconds())
				}
				r := RedisModel.NewListModel(key.(string), rds.LLen(key.(string)).Val(), expTime)
				listRes = append(listRes, r)
			case "set":
				exp := rds.TTL(key.(string)).Val()
				expTime := ""
				if exp.Seconds() < 0 {
					expTime = "-"
				} else {
					expTime = fmt.Sprintf("%.0f", exp.Seconds())
				}
				r := RedisModel.NewSetModel(key.(string), rds.SCard(key.(string)).Val(), expTime)
				setRes = append(setRes, r)
			}
		}
		return true
	})
	if stringRes == nil && listRes == nil && setRes == nil && nextCursor != 0 {
		cursorInt = int64(nextCursor)
		goto LABEL
	}

	if method == "string" {
		ctx.JSON(200, gin.H{"data": stringRes, "next": nextCursor, "current": cursor})
	} else if method == "list" {
		ctx.JSON(200, gin.H{"data": listRes, "next": nextCursor, "current": cursor})
	} else if method == "set" {
		ctx.JSON(200, gin.H{"data": setRes, "next": nextCursor, "current": cursor})
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
		Key   string      `json:"key"`
		Value string      `json:"value"`
		Exp   interface{} `json:"exp"`
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

	var expInt int64
	switch f.Exp.(type) {
	case string:
		expInt, _ = strconv.ParseInt(f.Exp.(string), 10, 64)
	case int:
		fmt.Printf("%T", f.Exp)
		expInt = int64(f.Exp.(int))
	}

	err = rds.Set(f.Key, f.Value, time.Duration(expInt*int64(math.Pow10(9)))).Err()
	if err != nil {
		ctx.JSON(400, gin.H{"message": "赋值失败"})
		return
	}
	ctx.JSON(200, gin.H{"message": "添加成功"})

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
		rm := RedisModel.NewListValueModel(key, v, (pageInt-1)*sizeInt+int64(i))
		vms = append(vms, rm)
	}
	ctx.JSON(200, gin.H{"data": vms, "len": lLen, "key": key})
}

func ListPost(ctx *gin.Context) {
	type form struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Exp    string `json:"exp"`
		Exists bool   `json:"exist"`
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
	if f.Exp != "" {
		expInt, _ := strconv.ParseInt(f.Exp, 10, 64)
		rds.Expire(f.Key, time.Duration(expInt*int64(math.Pow10(9))))
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

		fmt.Print(len(res), "-", cursor, "|")
		if cursor == 0 {
			if length == lLen {
				i--
			}
			break
		}
	}
	fmt.Println()
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
		Key   string `json:"key"`
		Value string `json:"value"`
		Exp   string `json:"exp"`
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

	if f.Exp != "" {
		expInt, _ := strconv.ParseInt(f.Exp, 10, 64)
		rds.Expire(f.Key, time.Duration(expInt*int64(math.Pow10(9))))
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
