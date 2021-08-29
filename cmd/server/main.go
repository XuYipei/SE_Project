package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuplz/tuplz-be/config"
	"github.com/tuplz/tuplz-be/internal/collect"
	"github.com/tuplz/tuplz-be/internal/comment"
	"github.com/tuplz/tuplz-be/internal/net"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"github.com/tuplz/tuplz-be/internal/prob"
	"github.com/tuplz/tuplz-be/internal/state"
	"github.com/tuplz/tuplz-be/internal/user"
	"go.mongodb.org/mongo-driver/bson"
)

func runUserServ(router *gin.Engine) {
	router.POST(RegisterUserUrl, net.RegisterUser)
	router.POST(VerifyEmailUrl, net.VerifyEmail)
	router.POST(VerifyUserUrl, net.VerifyUser)
	router.POST(UpdateUserUrl, net.UpdateUser)
	router.POST(LoginUserNameUrl, net.LoginUserName)
	router.GET(FindUserIdUrl, net.FindUserId)
}

func runCollectServ(router *gin.Engine) {
	router.GET(FindCollectUrl, net.FindCollect)
	router.GET(FindCollectsUrl, net.FindCollects)
	router.POST(UploadCollectsUrl, net.UploadCollects)
	router.POST(UploadCollectUrl, net.UploadCollect)
	router.PUT(UpdateCollectsUrl, net.UpdateCollects)
	router.DELETE(DeleteCollectsUrl, net.DeleteCollects)
}

func runCommentServ(router *gin.Engine) {
	router.GET(FindRecmdProbIdUrl, net.FindRecommendsProbId)
	router.GET(FindCommentsUrl, net.FindComments)
	router.GET(FindRecommendUrl, net.FindRecommendId)
	router.POST(UploadCommentUrl, net.UploadReply)
	router.POST(UploadRecommendUrl, net.UploadRecommend)
}

func runProbServ(router *gin.Engine) {
	router.GET(FindProbsUrl, net.FindProbs)
	router.GET(FindProbIdUrl, net.FindProbId)
	router.GET(FindProbRecommendUrl, net.FindProbRcmd)
}

func runMonitorServ(router *gin.Engine) {
	router.GET(DatabaseStatsUrl, net.FindDbStatus)
	router.GET(DeviceStatusUrl, net.FindDeviceStatus)
	router.POST(DatabaseSynclockUrl, net.LockSyncDatabase)
	router.POST(DatabaseUnlockUrl, net.UnlockDatabase)
	router.POST(DatabaseDumpUrl, net.DumpDatabase)
	router.POST(DatabaseStoreUrl, net.StoreDatabase)
}

func RunServer() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(Cors())

	runUserServ(router)
	runCommentServ(router)
	runProbServ(router)
	runMonitorServ(router)
	runCollectServ(router)

	router.Run(Port)
}

func InitHyper(reset bool, commentInit, userInit, problemInit, collectInit bool) (int, int, int, int, error) {
	type kvPair struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}
	var err0, err1, err2, err3 error
	if reset {
		err := tools.InitCollection(config.HyperTable)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		err0 = tools.InsertOne(config.Database, config.HyperTable, kvPair{Key: "comment", Value: 0})
		err1 = tools.InsertOne(config.Database, config.HyperTable, kvPair{Key: "user", Value: 0})
		err2 = tools.InsertOne(config.Database, config.HyperTable, kvPair{Key: "problem", Value: 0})
		err3 = tools.InsertOne(config.Database, config.HyperTable, kvPair{Key: "collect", Value: 0})
	}

	if commentInit && err0 == nil {
		err0 = tools.UpdateOne(config.Database, config.HyperTable, bson.M{"key": "comment"}, bson.M{"$set": bson.M{"value": 0}})
	}
	if userInit && err1 == nil {
		err1 = tools.UpdateOne(config.Database, config.HyperTable, bson.M{"key": "user"}, bson.M{"$set": bson.M{"value": 0}})
	}
	if problemInit && err2 == nil {
		err2 = tools.UpdateOne(config.Database, config.HyperTable, bson.M{"key": "problem"}, bson.M{"$set": bson.M{"value": 0}})
	}
	if collectInit && err3 == nil {
		err3 = tools.UpdateOne(config.Database, config.HyperTable, bson.M{"key": "collect"}, bson.M{"$set": bson.M{"value": 0}})
	}

	if err0 != nil {
		return 0, 0, 0, 0, err0
	}
	if err1 != nil {
		return 0, 0, 0, 0, err1
	}
	if err2 != nil {
		return 0, 0, 0, 0, err2
	}
	if err3 != nil {
		return 0, 0, 0, 0, err2
	}

	var recCount, userCount, probCount, collectCount kvPair
	mg := tools.FindOne(config.Database, config.HyperTable, bson.M{"key": "comment"})
	err := mg.Decode(&recCount)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	mg = tools.FindOne(config.Database, config.HyperTable, bson.M{"key": "user"})
	err = mg.Decode(&userCount)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	mg = tools.FindOne(config.Database, config.HyperTable, bson.M{"key": "problem"})
	err = mg.Decode(&probCount)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	mg = tools.FindOne(config.Database, config.HyperTable, bson.M{"key": "collect"})
	err = mg.Decode(&collectCount)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return userCount.Value, recCount.Value, probCount.Value, collectCount.Value, nil
}

func InitDatabase() error {
	err := state.Init()
	if err != nil {
		return err
	}
	userCount, recommendCount, problemCount, collectCount, err :=
		InitHyper(config.Init.InitHyper, config.Init.InitRecmd, config.Init.InitUser, config.Init.InitProb, config.Init.InitCollect)
	if err != nil {
		return err
	}

	err = comment.Init(config.Init.InitRecmd, recommendCount)
	if err != nil {
		return err
	}
	err = prob.Init(config.Init.InitProb, problemCount)
	if err != nil {
		return err
	}
	err = user.Init(config.Init.InitUser, userCount)
	if err != nil {
		return err
	}
	err = collect.Init(config.Init.InitCollect, collectCount)
	if err != nil {
		return err
	}

	return nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	err := InitDatabase()
	if err != nil {
		fmt.Println(err)
	}
	RunServer()
}
