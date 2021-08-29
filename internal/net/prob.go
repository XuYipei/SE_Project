package net

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tuplz/tuplz-be/internal/collect"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"github.com/tuplz/tuplz-be/internal/prob"
	"github.com/tuplz/tuplz-be/internal/rcmdsys"
	"github.com/tuplz/tuplz-be/internal/state"
	"go.mongodb.org/mongo-driver/bson"
)

type ProbReq struct {
	Id        string `json:"id"`
	MaxLength int    `json:"maxLength"`
	UserId    string `json:"userId"`
	Key       string `json:"key"`
}

type ProbRep struct {
	Favorite bool         `json:"favourite" form:"favourite"`
	Id       string       `json:"id" form:"content"`
	Like     int          `json:"like" form:"content"`
	Dislike  int          `json:"dislike" form:"content"`
	Visit    int          `json:"visit" form:"content"`
	Content  prob.Content `json:"content" form:"content"`
	// UpdTime time.Time    `json:"updTime"`
}
type ProbSignal struct {
	// Status string  `json:"status"`
	RespHeader
	Prob ProbRep `json:"problem"`
}
type ProbMulti struct {
	Status string    `json:"status"`
	Prob   []ProbRep `json:"problems"`
}

func toProbRep(result prob.ProbStruct) ProbRep {
	return ProbRep{
		Id:      result.Id,
		Content: result.Content,
		Like:    result.Like,
		Dislike: result.Dislike,
		// UpdTime: result.UpdTime,
		Visit: result.Visit,
	}
}

// 根据题目 id 展示题目
func FindProbId(c *gin.Context) {
	var req struct {
		Id     string `json:"id" form:"id"`
		UserId string `json:"userId" form:"userId"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindQuery(&req)
	req.Id = c.Param("id")
	req.Key = GetAuth(c)

	var resp ProbSignal

	var result prob.ProbStruct
	var errP error
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		result, errP = prob.FindId(req.Id)
	}(&wg)

	// log.print(req)

	var wgF sync.WaitGroup
	var favor, online bool
	wgF.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wgF.Done()
		online, _ = state.FindUpdate(req.UserId, req.Key)
	}(&wgF)
	wgF.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wgF.Done()
		favor, _ = collect.FindFavorite(req.UserId, tools.StringToInt(req.Id))
	}(&wgF)

	wgF.Wait()
	wg.Wait()

	if online && errP == nil {
		rcmdsys.Upd(req.UserId, req.Id, 3.)
	}

	resp.Handle(errP)
	resp.Prob = toProbRep(result)
	resp.Prob.Favorite = favor && online
	prob.Visit(req.Id)
	c.JSON(http.StatusOK, resp)
}

// 展示所有题目
func FindProbs(c *gin.Context) {
	var req struct {
		UserId    string `json:"userId" form:"userId"`
		Key       string `json:"key" form:"key"`
		MaxLength int    `json:"maxLength" form:"maxLength"`
	}
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)

	// log.print(req)
	result, err := prob.FindProbs(req.MaxLength)
	state.FindUpdate(req.UserId, req.Key)

	rep := struct {
		Status string            `json:"status" form:"status"`
		Prob   []prob.ProbStruct `json:"problems" form:"problems"`
	}{Status: "success", Prob: result}

	if err != nil {
		rep.Status = "fail"
	}
	c.JSON(http.StatusOK, rep)
}

// 展示值得推荐的题目
func FindProbRcmd(c *gin.Context) {
	var req struct {
		UserId    string `json:"userId" form:"userId"`
		Key       string `json:"key" form:"key"`
		MaxLength int    `json:"maxLength" form:"maxLength"`
	}
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)

	var resp struct {
		RespHeader
		Prob []prob.ProbStruct `json:"problems" form:"problems"`
	}

	log.Println(req)

	online, _ := state.FindUpdate(req.UserId, req.Key)
	if !online {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	if req.MaxLength == 0 || req.MaxLength > 10 {
		req.MaxLength = 10
	}
	probs, err := rcmdsys.Query(req.UserId, req.MaxLength)

	resp.Handle(err)
	if err != nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Prob = make([]prob.ProbStruct, len(probs))
	for i, pId := range probs {
		probC, err := prob.FindOneProb(bson.M{"id": pId})
		if err != nil {
			resp.Handle(err)
			break
		}
		log.Println(probC)
		resp.Prob[i] = probC
	}
	c.JSON(http.StatusOK, resp)
}
