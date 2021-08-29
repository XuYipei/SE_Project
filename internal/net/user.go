package net

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuplz/tuplz-be/internal/collect"
	"github.com/tuplz/tuplz-be/internal/state"
	"github.com/tuplz/tuplz-be/internal/user"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	UserInfoPost     = "userInfo"
	UserIdPost       = "userId"
	UserEmailPost    = "userEmail"
	UserPasswordPost = "userPassword"
	ProblemIdPost    = "problemId"
)

type UserReq struct {
	Id         string `json:"id"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	VerifyCode string `json:"verifyCode"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Key      string `json:"key"`
}
type UserRep struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}
type UsersRep struct {
	Status string `json:"status"`
	User   []User `json:"user"`
}

func RegisterUser(c *gin.Context) {
	var jsons struct {
		Name     string `json:"username" form:"username"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	c.ShouldBindJSON(&jsons)
	uInfo := user.User{
		Name:     jsons.Name,
		Password: jsons.Password,
		Email:    jsons.Email,
	}

	// log.print(jsons)
	var resp struct {
		RespHeader
		Id   string `json:"id"`
		Name string `json:"username"`
		Key  string `json:"key"`
	}

	result, err := user.Register(uInfo)
	resp.Id = result.Id
	resp.Name = result.Name
	resp.Handle(err)
	if err != nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	err = collect.UploadUser(result.Id)
	resp.Handle(err)
	if err != nil {
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Key, err = state.Login(result.Id)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func VerifyEmail(c *gin.Context) {
	var req struct {
		Id  string `json:"userId"`
		Key string `json:"key"`
	}
	c.ShouldBind(&req)
	req.Key = GetAuth(c)

	var resp struct {
		RespHeader
	}
	online, err := state.FindUpdate(req.Id, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	_, err = user.SendEmail(req.Id)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func VerifyUser(c *gin.Context) {
	var req struct {
		Id   string `json:"userId"`
		Code string `json:"verifyCode"`
		Key  string `json:"key"`
	}
	c.ShouldBind(&req)
	req.Key = GetAuth(c)

	// log.print(req)
	var resp struct {
		RespHeader
	}

	online, err := state.FindUpdate(req.Id, req.Key)
	if !online || err != nil {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	_, err = user.VerifyUser(req.Id, req.Code)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func LoginUserName(c *gin.Context) {
	var jsons struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}
	c.ShouldBind(&jsons)
	uInfo := user.User{
		Name:     jsons.Name,
		Password: jsons.Password,
	}

	result, err := user.LoginName(uInfo)
	resp := struct {
		RespHeader
		Id   string `json:"id"`
		Key  string `json:"key"`
		Name string `json:"username"`
	}{
		Id:   result.Id,
		Key:  "",
		Name: result.Name,
	}
	if err != nil {
		resp.Handle(err)
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Key, err = state.Login(result.Id)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func UpdateUser(c *gin.Context) {
	var req struct {
		Id     string `json:"id"`
		Name   string `json:"username"`
		Disc   string `json:"description"`
		Gender bool   `json:"gender"`
		Key    string
	}
	c.ShouldBindJSON(&req)
	req.Key = GetAuth(c)

	// log.print(req)

	userData := user.User{
		Name: req.Name,
		Id:   req.Id,
	}
	infoData := user.Information{
		Id:     req.Id,
		Disc:   req.Disc,
		Gender: req.Gender,
	}
	var resp struct {
		RespHeader
	}

	check, err := state.FindUpdate(req.Id, req.Key)
	online := (check && err == nil)
	if !online {
		resp.Handle(errors.New("state: state is not online"))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = user.Update(userData, infoData)
	resp.Handle(err)
	c.JSON(http.StatusOK, resp)
}

func FindUserId(c *gin.Context) {
	var req struct {
		Id  string `json:"userId" form:"userId"`
		Key string `json:"key" form:"key"`
	}
	req.Key = GetAuth(c)
	req.Id = c.Param("id")
	var resp struct {
		RespHeader
		User struct {
			Id       string `json:"userId"`
			Name     string `json:"username"`
			Email    string `json:"email"`
			Verified bool   `json:"isVerified"`
		} `json:"user"`
	}

	// log.print(req)

	state.FindUpdate(req.Id, req.Key)
	// online, err := state.FindUpdate(req.Id, req.Key)
	// if err != nil || !online {
	// 	resp.Handle(errors.New("state: state is not online"))
	// 	c.JSON(http.StatusOK, resp)
	// 	return
	// }

	uInfo, erru := user.FindOneUser(bson.M{"id": req.Id})
	_, erri := user.FindOneUserInfo(bson.M{"id": req.Id})
	resp.User.Id = uInfo.Id
	resp.User.Email = uInfo.Email
	resp.User.Verified = uInfo.Verify.State
	resp.User.Name = uInfo.Name

	if erru != nil {
		resp.Handle(erru)
	}
	if erri != nil {
		resp.Handle(erri)
	}
	c.JSON(http.StatusOK, resp)
}
