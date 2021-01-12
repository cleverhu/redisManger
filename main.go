package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"redisManger/src/common"
	"redisManger/src/dbs"
	"redisManger/src/handlers"
	"redisManger/src/utils/myHttp"
	"time"
)

func main() {
	//
	//insertData()
	//return
	r := gin.New()
	r.Use(common.Cors())
	r.Use(common.Auth())

	{
		str := r.Group("/string")
		str.GET("/scan", handlers.Scan)
		str.DELETE("/:keys", handlers.DelByKeys)
		str.POST("/", handlers.StringUpdate)

		str.GET("/toFile", handlers.ToFile)
	}

	{
		list := r.Group("/list")
		list.GET("/scan", handlers.Scan)
		list.GET("/get/:key", handlers.ListGetByKey)
		list.DELETE("/:keys", handlers.DelByKeys)

		list.POST("/", handlers.ListPost)

		list.GET("/exists/:key", handlers.ListExists)

		list.POST("/remove", handlers.ListRemoveValue)
		list.POST("/insert", handlers.ListInsert)
	}

	{
		set := r.Group("/set")
		set.GET("/scan", handlers.Scan)
		set.POST("/", handlers.SetPost)
		set.GET("/get/:key", handlers.SetGetByKey)
		set.DELETE("/:keys", handlers.DelByKeys)
		set.POST("/remove", handlers.SetRemoveValue)
		set.POST("/getCommon", handlers.SetGetCommon)
	}

	{
		hash := r.Group("/hash")

		hash.GET("/scan", handlers.Scan)
		hash.DELETE("/:keys", handlers.DelByKeys)

		hash.GET("/get/:key", handlers.HashGetByKey)
		hash.POST("/remove", handlers.HashRemoveValue)
		hash.POST("/", handlers.HashPost)

	}

	{
		cfg := r.Group("/config")
		cfg.GET("/", handlers.Config)
		cfg.POST("/", handlers.UpdateConfig)
	}

	r.GET("/info/", handlers.Info)

	{
		user := r.Group("/user")
		user.POST("/login", handlers.Login)
		user.POST("/check", handlers.Validate)
	}

	{
		connect := r.Group("/connect")

		connect.POST("/test", handlers.ConnectTest)

		connect.POST("/save", handlers.ConnectSave)
		connect.GET("/get", handlers.ConnectGet)
	}

	{
		logs := r.Group("/logs", handlers.LogsGet)
		logs.POST("/get")
	}

	log.Fatal(r.Run(":80"))

}

func insertData() {
	//task := make(chan int, 10)
	//for i := 0; i < 20000; i++ {
	//	task <- 1
	//	go func(i int) {
	//		dbs.Rds.Set(fmt.Sprintf("%d", i), fmt.Sprintf("hi %d", i), time.Duration(rand.Int()))
	//		<-task
	//	}(i)
	//}

	type result struct {
		Id          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
	}

	for i := 1; i <= 1; i++ {
		_, str, _, _ := myHttp.Request("get", "https://so.csdn.net/api/v2/search?q=golang&t=all&p="+fmt.Sprintf("%d", i)+"&s=0&tm=0&lv=-1&ft=0&l=&u=&platform=pc", nil, nil, "cookie:uuid_tt_dd=10_7922026970-1592972787364-295687; _ga=GA1.2.2129730115.1592972789; dc_sid=c11f7409ea3f2db5aef3536ddfc37ea0; UN=XiaoHuLearn; Hm_ct_6bcd52f51e9b3dce32bec4a3997715ac=6525*1*10_7922026970-1592972787364-295687!5744*1*XiaoHuLearn; c_segment=7; p_uid=U110000; UserName=XiaoHuLearn; UserInfo=3835fc0946fc49ce9d469f4044fb12ba; UserToken=3835fc0946fc49ce9d469f4044fb12ba; UserNick=XiaoHuLearn; AU=F22; BT=1599741433952; Hm_up_6bcd52f51e9b3dce32bec4a3997715ac=%7B%22uid_%22%3A%7B%22value%22%3A%22XiaoHuLearn%22%2C%22scope%22%3A1%7D%2C%22islogin%22%3A%7B%22value%22%3A%221%22%2C%22scope%22%3A1%7D%2C%22isonline%22%3A%7B%22value%22%3A%221%22%2C%22scope%22%3A1%7D%2C%22isvip%22%3A%7B%22value%22%3A%221%22%2C%22scope%22%3A1%7D%7D; SESSION=9c606e00-ee10-4840-b09b-19db2cc2d076; __gads=ID=614dc74697a7dab9-226bad5b1ac300ff:T=1598172190:RT=1598172190:R:S=ALNI_MbyjRZuyBOBFIuWuVikvv7rvU8F6g; aliyun_webUmidToken=T2gACDm43YE8Xi_zVyVECwy062DIOsRiDS29xvk-kwixqWIqiLiAp_FFEdEMVvEGEk8FwTvav7gAplJ6Hy7JSM6L; aliyun_UAToken=137#Zgc9hE9oUNqMSfZii9EFJ9zNIQk/t5SejH2dJ6AHGLJ95F8KgbfGRbgd3aENtnZQf+icxDDxTGWTRKT6Tm1QGkGpp+YheW0Ve/iSjHtgbaCfyv+NApyuiE/PYZewsdvC63Kdlsqna/3iJBrwni6NNV+6RgKJhqwWxD6+nDsJkMq++Wlm5hYRTKM5crxVBbQKHR8cI79VPGHtcnH+/pTNy6aJyTpVxl30rWS3shSzldNGzZBIRBoMKDMua6dvDIdUXBxEEnUMoQOJB3KvtGGSVOePurTA1e7KPjmcy+xdH3eWDifhMzQ7FhCWuN2cYUyJjQ6ZPsQBbDh3xzx5G581g44jJhTw6A6JJABYLqf0Lj+A0ka1pE6Z8c+kKcihksQt0qfIt76/xVWWL7vsRcdgpNpBmhmDDZFsmoiAJpjYVJJ64i82h0PTQUHf1XTMjh5fEOvBK4jdDpkzX57pcOOM8nuzAUKd+Jdjmr6fCAS+i+p3O4dbTS2xo728oa/R2GXaYjnnC9kN6ww08KEKZEDwGte4KcU1518+g+axdknldudw+47Ved+2at2Lm2kldaAu2DTp/5OjlH1MO6Ipan+6w3A0Lq5UuUUZ/qWVm+WjZh6E5A0By88n3ced4myZN3hyr9iIMp7s4kKmFNOO+ASdH/ff+1RvJmkR1UZtS2+jKv2K2hE47CIoflKGrh0f7I3v/76qxBOyN0UXwmuOvrIgw6tuxUa2Tv2x1Zcnzv/TcQpfZ5bSEhkp3HiRyp0dVKUh9AsM0Bg8DaVRrytDhhPTbIEuS7Fep0otAT46yM2vz8fmde4jQywnfrpnjBbHXSQ/Eylwb3gbQbpVIuAHxZW2M64pgf2pBImOZciIMCdFD1CEYTUx1lgIp8xCQeLUdqippRJm1IESbUFslyJx1lQypXicQon++ZDpMpkRwqwy+B6pYSUS1lQiSdC9LInJ+ZDVpT17YIEQNAXp1TxZ1qgippimQefo+ZXVpkUmuIEy+piVYSJS1qQipXpcQeLk9Glh4dP9iIS26HSO0juOmbZJMc9oaOOwJEAw5l62P2UYdmoXeSR7oYPYQAELtcZTSlYCvj9UkJuRlo5QxEAAqFmz5ruZY2JwgiOt9wQo+Y4cjQSlv/Li4mQc4pZ5Oswxh6iev3MUg0CjIEATVf2anJeNDW+X6Ui4u2tmyHriZbNnxBWSkVFSq+wEb0TJaEfTBeYf7w6XRPxLtBvrcJVs87j32hUGunpIYGsYvvVm6l6TNueAANijQK7qd6vSFJMfOOddcL436S4p4iSLxgZg40zoz2UvLFHf4YMwtJftmcJKpYS333YSFgo5L8peE4dDpUUAfadNeiy2sp3MqnNireIrNnXqVbDGUO+EQ4FzQYWCBfiP95MBT3ZqKW1dAGEh5QsjZnNbxoijwvAT7IHw/R48AGZyMDlXdHiRViqiB/wPxhx2cM6LU6znMdVZz9em6AJn0YSYNB6HxHQ4JZ3sLl2D06WZVoRWqZ1G9526KoC7G/O9iPF+8z9To7PpA3atER6wH3QLfUUbC3JfnJHR5jR+wj4oy8Jt5zTF/yU/SUetMdqplRHpYGnlToA39Mcp6jQTgNvsMcxC4hXjEKSY; c_first_ref=www.baidu.com; announcement-new=%7B%22isLogin%22%3Atrue%2C%22announcementUrl%22%3A%22https%3A%2F%2Fblog.csdn.net%2Fblogdevteam%2Farticle%2Fdetails%2F112280974%3Futm_source%3Dgonggao_0107%22%2C%22announcementCount%22%3A1%2C%22announcementExpire%22%3A249394776%7D; c_hasSub=true; dc_session_id=10_1610106432478.737991; searchHistoryArray-new=%5B%22golang%22%2C%22go%22%5D; c_utm_term=golang; c_utm_medium=distribute.pc_search_result.none-task-blog-2%7Eall%7Esobaiduweb%7Edefault-0-89917456.first_rank_v2_pc_rank_v29; log_Id_click=111; BAIDU_SSP_lcr=https://www.baidu.com/link?url=6g5WAIymblcq6VwQGcKigwARjvWwgG_-KArP45UiG3ZTADvSIrHMaCX3hY6GbAOq&wd=&eqid=8bd91e400004125a000000045ff84809; c_ref=https%3A//www.baidu.com/link; c_first_page=https%3A//blog.csdn.net/qq_39893313/article/details/106575947; c_page_id=default; dc_tos=qmm6g3; log_Id_pv=754; c_pref=https%3A//www.baidu.com/link; Hm_lvt_6bcd52f51e9b3dce32bec4a3997715ac=1610106434,1610106590,1610106910,1610106916; Hm_lpvt_6bcd52f51e9b3dce32bec4a3997715ac=1610106916; log_Id_view=2177", "", 30)
		//fmt.Println(str)
		var j interface{}
		err := json.Unmarshal([]byte(str), &j)
		if err == nil {
			d := j.(map[string]interface{})["result_vos"].([]interface{})

			for _, v := range d {
				//	fmt.Println(v)
				m := v.(map[string]interface{})

				k := "hash_" + m["id"].(string)
				//title := decode.UnicodeToUTF8(m["title"].(string))
				////fmt.Println(k, title)
				//desc := decode.UnicodeToUTF8(m["description"].(string))
				//fmt.Println(desc)
				r := &result{
					Id:          k,
					Title:       m["title"].(string),
					Description: m["description"].(string),
					Url:         m["url"].(string),
				}
				//data, _ := json.Marshal(r)

				rand.Seed(time.Now().UnixNano())
				//dbs.Rds.Set(k, string(data), time.Duration(rand.Int63()))
				//for key, value := range m {
				//	dbs.Rds.RPush(r.Id, fmt.Sprintf("%v:%v", key, value))
				//}

				for key, value := range m {
					//dbs.Rds.SAdd(r.Id, fmt.Sprintf("%v:%v", key, value))
					dbs.Rds.HSet(r.Id, key, value)
				}

				dbs.Rds.Expire(k, time.Duration(rand.Int63()))
				//dbs.Rds.Expire(r.Id, time.Duration(1*math.Pow10(9)))
			}
		}

	}

	return
}
