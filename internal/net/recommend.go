package net

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tuplz/tuplz-be/internal/comment"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"github.com/tuplz/tuplz-be/internal/state"
	"github.com/tuplz/tuplz-be/internal/user"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

type Recommend struct {
	UserId  string `json:"userId"`
	ProbId  int    `json:"probId" form:"problemId"`
	Message string `json:"message" form:"recommendReason"`
}

type RecommendRep struct {
	Status    string    `json:"status"`
	Recommend Recommend `json:"recommend"`
}

type RecommendsRep struct {
	Status     string      `json:"status"`
	Recommdnes []Recommend `json:"recommends"`
}

func toRecommend(cInfo comment.Comment) Recommend {
	return Recommend{
		UserId:  cInfo.Content.UserId,
		ProbId:  cInfo.ProbId,
		Message: cInfo.Content.Content,
	}
}

// 上传推荐
func UploadRecommend(c *gin.Context) {
	var req struct {
		UserName string `json:"userName" form:"userName"`
		UserId   string `json:"userId" form:"userId"`
		ProbId   int    `json:"problemId" form:"problemId"`
		Content  string `json:"recommendReason" form:"recommendReason"`
		Key      string `json:"key" form:"key"`
	}
	c.ShouldBindJSON(&req)
	req.Key = GetAuth(c)
	uInfo, _ := user.FindOneUser(bson.M{"id": req.UserId})
	req.UserName = uInfo.Name

	// log.print(req)

	var resp struct {
		RespHeader
		Recommend Recommend
	}

	rInfo := comment.Comment{
		ProbId: req.ProbId,
		Content: comment.Content{
			UserName: req.UserName,
			UserId:   req.UserId,
			Content:  req.Content,
			Time:     time.Now(),
		},
	}

	check, err := state.FindUpdate(req.UserId, req.Key)
	online := (check && err == nil)
	if !online {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	rInfo, err = comment.Upload(rInfo)
	resp.Recommend = toRecommend(rInfo)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

// 上传回复
func UploadReply(c *gin.Context) {
	var req struct {
		RecommendId int    `json:"recommendId" form:"recommendId"`
		UserId      string `json:"userId" form:"userId"`
		UserName    string `json:"username" form:"username"`
		Reply       int    `json:"replyTo" form:"replyTo"`
		Content     string `json:"commentContent" form:"commentContent"`
		Key         string
	}
	c.ShouldBindJSON(&req)
	req.Key = GetAuth(c)
	uInfo, _ := user.FindOneUser(bson.M{"id": req.UserId})
	req.UserName = uInfo.Name
	req.RecommendId, _ = strconv.Atoi(c.Param("id"))
	if req.Reply <= 0 {
		req.Reply = req.RecommendId
	}

	// log.print(req)

	var resp struct {
		RespHeader
	}
	rInfo := comment.Comment{
		Reply:     req.Reply,
		Recommend: req.RecommendId,
		Content: comment.Content{
			UserId:   req.UserId,
			UserName: req.UserName,
			Content:  req.Content,
			Time:     time.Now(),
		},
	}
	check, err := state.FindUpdate(req.UserId, req.Key)
	online := (check && err == nil)
	if !online {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	rInfo, err = comment.Upload(rInfo)
	if err == nil {
		err = comment.AddResponses(
			req.RecommendId,
			rInfo.Id,
		)
	}
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

// 展示推荐的回复
func FindComments(c *gin.Context) {
	var req struct {
		RecommendId int    `form:"recommendId"`
		UserId      string `form:"userId"`
		Key         string `form:"key"`
		MaxLength   int    `form:"maxLength"`
	}
	var err error
	req.RecommendId, _ = strconv.Atoi(c.Param("id"))
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)

	result, err := comment.FindCommentId(req.RecommendId)
	type resp struct {
		Id        int       `json:"commentId" form:"commentId"`
		UserId    string    `json:"userId" form:"userId"`
		UserName  string    `json:"username" form:"username"`
		IsReply   bool      `json:"isReply" form:"isReply"`
		Content   string    `json:"commentContent" form:"commentContent"`
		Time      time.Time `json:"updateTime" form:"updateTime"`
		Responses []int     `json:"responses" form:"responses"`
		ReplyTo   int       `json:"replyTo" form:"replyTo"`
	}
	transform := func(input comment.Comment) resp {
		return resp{
			Id:        input.Id,
			UserId:    input.Content.UserId,
			UserName:  input.Content.UserName,
			Content:   input.Content.Content,
			Time:      input.Content.Time,
			ReplyTo:   input.Reply,
			Responses: input.Responses,
			IsReply:   false,
		}
	}
	recommends := make([]resp, len(result.Responses))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":   "fail",
			"comments": recommends,
		})
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		state.FindUpdate(req.UserId, req.Key)
	}(&wg)
	for i, id := range result.Responses {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, id int) {
			defer wg.Done()
			reply, errC := comment.FindCommentId(id)
			if errC != nil {
				err = errC
			}
			recommends[i] = transform(reply)
		}(&wg, i, id)
	}
	wg.Wait()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":   "fail",
			"comments": recommends,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"comments": recommends,
	})
}

func FindRecommendId(c *gin.Context) {
	var req struct {
		Id     int    `form:"id"`
		UserId string `form:"userId"`
		Key    string `form:"key"`
	}
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)
	req.Id, _ = strconv.Atoi(c.Param("id"))

	state.FindUpdate(req.UserId, req.Key)
	result, err := comment.FindCommentId(req.Id)

	resp := struct {
		Id       int       `json:"recommendId" form:"recommendId"`
		UserId   string    `json:"userId" form:"userId"`
		UserName string    `json:"username" foem:"username"`
		ProbId   int       `json:"problemId" form:"problemId"`
		Content  string    `json:"message" form:"message"`
		Time     time.Time `json:"updateTime" form:"updateTime"`
	}{
		Id:       result.Id,
		UserId:   result.Content.UserId,
		UserName: result.Content.UserName,
		ProbId:   result.ProbId,
		Content:  result.Content.Content,
		Time:     result.Content.Time,
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"recommend": resp,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    "fail",
			"recommend": resp,
		})
	}
}

func FindRecommendsProbId(c *gin.Context) {
	var req struct {
		ProbId    string `form:"problemId"`
		UserId    string `form:"userId"`
		Key       string `form:"key"`
		MaxLength int    `form:"maxLength"`
	}
	req.ProbId = c.Param("id")
	req.Key = GetAuth(c)
	c.ShouldBindQuery(&req)

	// log.print(req)

	state.FindUpdate(req.UserId, req.Key)
	result, err := comment.FindRecommendsProbId(
		tools.StringToInt(req.ProbId),
		req.MaxLength,
	)

	// log.print(result)

	recommends := make(
		[]struct {
			Id       int       `json:"recommendId" form:"recommendId"`
			UserId   string    `json:"userId" form:"userId"`
			UserName string    `json:"username" foem:"username"`
			ProbId   int       `json:"problemId" form:"problemId"`
			Content  string    `json:"message" form:"message"`
			Time     time.Time `json:"updateTime" form:"updateTime"`
		}, len(result),
	)
	for i, x := range result {
		recommends[i].Id = x.Id
		recommends[i].ProbId = x.ProbId
		recommends[i].UserId = x.Content.UserId
		recommends[i].UserName = x.Content.UserName
		recommends[i].Content = x.Content.Content
		recommends[i].UserName = x.Content.UserName
		recommends[i].Time = x.Content.Time
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"recommends": recommends,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":     "fail",
			"recommends": recommends,
		})
	}
}
