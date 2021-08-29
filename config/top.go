package config

import "time"

var (
	/*
		ClientUser: 用户名
		ClinetPassword: 用户密码
		DatabaseName: 数据库名
		UserCollectionName: 用户表名
		RecommendCollectionName: 推荐表名
	*/
	ClientUser     string = "admin"
	ClientPassword string = "123456"
	Database       string = "seproject"
)
var (
	UserTable           string = "user"
	UserContentTable    string = "userContent"
	CommentTable        string = "comment"
	CommentContentTable string = "commentContent"
	StateTable          string = "state"
	ProblemTable        string = "problem"
	ProbContentTable    string = "problemContent"
	HyperTable          string = "hyper"
	CollectTable        string = "collect"
	CollectContentTable string = "collectContent"
)

var (
	ProblemDir  string = "./assets/problemset"
	DatabaseDir string = "./assets/database"
)

var (
	LoadProblemWorkers int = 64
)

type inits struct {
	InitHyper   bool
	InitUser    bool
	InitProb    bool
	InitRecmd   bool
	InitCollect bool
}

var (
	EmailUser string = "seporjectteam4@163.com"
	EmailAuth string = "BDVPJHVKVTDANTBM"
	EmailHost string = "smtp.163.com"
	EmailPort string = "25"
)

var Limits = struct {
	UOnline        time.Duration
	VerifyTime     time.Duration
	ProbContent    int
	ProbSections   int
	CommentContent int
}{
	UOnline:        3600 * time.Minute,
	VerifyTime:     5 * time.Minute,
	ProbContent:    50,
	CommentContent: 50,
	ProbSections:   1,
}

var Init = inits{
	InitHyper:   false,
	InitUser:    false,
	InitProb:    false,
	InitRecmd:   false,
	InitCollect: false,
}
