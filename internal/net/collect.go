package net

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuplz/tuplz-be/internal/collect"
	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"github.com/tuplz/tuplz-be/internal/prob"
	"github.com/tuplz/tuplz-be/internal/rcmdsys"
	"github.com/tuplz/tuplz-be/internal/state"
)

// 更改题目收藏状态
func UploadCollect(c *gin.Context) {
	var req struct {
		UserId    string `json:"userId" form:"userId"`
		ProbId    int    `json:"probId" form:"probId"`
		Key       string `json:"key" form:"key"`
		CollectId int    `json:"collectionId" form:"collectionId"`
	}
	c.ShouldBindJSON(&req)
	req.Key = GetAuth(c)
	req.ProbId, _ = strconv.Atoi(c.Param("id"))

	var resp struct {
		RespHeader
		Prob ProbRep `json:"problem"`
	}

	var cId int
	var online, favorite bool
	var err, err0, err1, errP error
	var wg, wgProb sync.WaitGroup
	var probC prob.ProbStruct

	log.Println(req)

	wgProb.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		probC, errP = prob.FindId(strconv.Itoa(req.ProbId))
		resp.Prob = toProbRep(probC)
	}(&wgProb)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		online, err = state.FindUpdate(req.UserId, req.Key)
	}(&wg)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		favorite, cId = collect.FindFavorite(req.UserId, req.ProbId)
	}(&wg)
	wg.Wait()

	if !online || err != nil {
		wgProb.Wait()
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	wgProb.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rcmdsys.Upd(req.UserId, strconv.Itoa(req.ProbId), 5.)
	}(&wgProb)

	favor := false
	if !favorite {
		favor = true
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err0 = collect.UploadCollect(req.CollectId, req.ProbId)
		}(&wg)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err1 = prob.UpdateLike(strconv.Itoa(req.ProbId), 1)
		}(&wg)
	} else {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err0 = collect.DeleteCollect(req.UserId, cId, req.ProbId)
		}(&wg)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err1 = prob.UpdateLike(strconv.Itoa(req.ProbId), -1)
		}(&wg)
	}
	wg.Wait()
	wgProb.Wait()

	resp.Prob.Favorite = favor
	resp.Handle(err0)
	resp.Handle(err1)
	resp.Handle(errP)
	c.JSON(http.StatusOK, resp)
}

func UploadCollects(c *gin.Context) {
	var req struct {
		UserId string `json:"userId" form:"userId"`
		Title  string `json:"title" form:"title"`
		Disc   string `json:"disc" form:"disc"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindJSON(&req)
	req.Key = GetAuth(c)

	// log.print(req)

	var resp struct {
		RespHeader
	}

	online, err := state.FindUpdate(req.UserId, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = collect.UploadCollects(req.UserId, req.Title, req.Disc)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

type CollectHeader struct {
	Id      int       `json:"collectionId"`
	UserId  string    `json:"userId"`
	Title   string    `json:"title"`
	UpdTime time.Time `json:"updateTime"`
}

func FindCollect(c *gin.Context) {
	var req struct {
		Id     int    `json:"id" form:"id"`
		UserId string `json:"userId" form:"userId"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)
	req.Id, _ = strconv.Atoi(c.Param("id"))

	// log.print(req)

	state.FindUpdate(req.UserId, req.Key)
	result, err := collect.FindCollectId(req.Id)

	type problem struct {
		Content  prob.Content `json:"content"`
		Id       int          `json:"id"`
		Like     int          `json:"like"`
		Dislike  int          `json:"dislike"`
		Visit    int          `json:"visit"`
		Favorite bool         `json:"favorite"`
	}
	toProbs := func(data []prob.ProbStruct) []problem {
		result := make([]problem, len(data))
		for i, x := range data {
			result[i].Content = x.Content
			result[i].Id = x.Index
			result[i].Dislike = x.Dislike
			result[i].Like = x.Like
			result[i].Favorite = true
		}
		return result
	}

	type Collect struct {
		CollectHeader
		Problems []problem `json:"problems"`
	}
	var resp struct {
		RespHeader
		Collect Collect `json:"collection"`
	}
	resp.Collect = Collect{
		CollectHeader: CollectHeader{
			Id:      result.Id,
			UserId:  result.UserId,
			Title:   result.Name,
			UpdTime: result.UpdTime,
		},
		Problems: toProbs(result.ProbCs),
	}
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func FindCollects(c *gin.Context) {
	var req struct {
		UserId string `json:"userId" form:"userId"`
		Key    string `json:"Key" form:"Key"`
	}
	c.ShouldBindQuery(&req)
	req.Key = GetAuth(c)

	// log.print(req)

	state.FindUpdate(req.UserId, req.Key)
	result, err := collect.FindCollects(req.UserId)

	type Collect struct {
		CollectHeader
		Count int `json:"problemCount"`
	}
	var resp struct {
		RespHeader
		Data []Collect `json:"collections"`
	}
	resp.Data = make([]Collect, len(result))
	for i, x := range result {
		resp.Data[i] = Collect{
			Count: x.Count,
			CollectHeader: CollectHeader{
				Id:      x.Id,
				UserId:  req.UserId,
				Title:   x.Name,
				UpdTime: x.UpdTime,
			},
		}
	}
	resp.Handle(err)

	c.JSON(http.StatusOK, resp)
}

func UpdateCollects(c *gin.Context) {
	var req struct {
		Id     int    `json:"id" form:"id"`
		UserId string `json:"userId" form:"userId"`
		Title  string `json:"title" form:"title"`
		Disc   string `json:"disc" form:"disc"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindJSON(&req)
	req.Id = tools.StringToInt(c.Param("id"))
	req.Key = GetAuth(c)

	// log.print(req)

	var resp struct {
		RespHeader
	}

	online, err := state.FindUpdate(req.UserId, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = collect.UpdateCollects(req.UserId, req.Id, req.Title, req.Disc)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func DeleteCollects(c *gin.Context) {
	var req struct {
		Id     int    `json:"id" form:"id"`
		UserId string `json:"userId" form:"userId"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindJSON(&req)
	req.Id = tools.StringToInt(c.Param("id"))
	req.Key = GetAuth(c)

	// log.print(req)

	var resp struct {
		RespHeader
	}

	online, err := state.FindUpdate(req.UserId, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = collect.DeleteCollects(req.Id)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func CopyCollects(c *gin.Context) {
	var req struct {
		Id     int    `json:"id" form:"id"`
		UserId string `json:"userId" form:"userId"`
		Key    string `json:"key" form:"key"`
	}
	c.ShouldBindJSON(&req)
	req.Id = tools.StringToInt(c.Param("id"))
	req.Key = GetAuth(c)

	// log.print(req)

	var resp struct {
		RespHeader
	}

	online, err := state.FindUpdate(req.UserId, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = collect.DeleteCollects(req.Id)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}
