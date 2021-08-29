package main

var (
	IpAddress = "http://localhost"
	Port      = ":5000"
)

var (
	RegisterUserUrl  = "/api/user/register"
	VerifyEmailUrl   = "/api/user/verify/email"
	VerifyUserUrl    = "/api/user/verify/user"
	LoginUserNameUrl = "/api/user/login"
	UpdateUserUrl    = "/api/user/update"
	FindUserIdUrl    = "/api/users/:id"
	// RemoveUserUrl    = "/api/user/remove"
	// LogoutUserUrl = "/api/user/logout"
)

var (
	UploadCollectsUrl = "/api/collections"
	FindCollectsUrl   = "/api/collections"
	DeleteCollectsUrl = "/api/collections/:id"
	UpdateCollectsUrl = "/api/collections/:id"
	FindCollectUrl    = "/api/collections/:id"
	UploadCollectUrl  = "/api/problems/collection/:id"
)

var (
	FindRecmdProbIdUrl = "/api/problems/:id/recommends"
	FindRecommendUrl   = "/api/recommends/:id"
	FindCommentsUrl    = "/api/recommends/:id/comments"
	UploadCommentUrl   = "/api/recommends/:id/comments"
	UploadRecommendUrl = "/api/recommends"
)

var (
	FindProbIdUrl        = "/api/problems/:id"
	FindProbsUrl         = "/api/problems"
	FindProbRecommendUrl = "/api/recommendedproblems"
)

var (
	DeviceStatusUrl     = "/api/monitor/device"
	DatabaseStatsUrl    = "/api/monitor/database"
	DatabaseSynclockUrl = "/api/monitor/synclock"
	DatabaseUnlockUrl   = "/api/monitor/unlock"
	DatabaseDumpUrl     = "/api/monitor/dump"
	DatabaseStoreUrl    = "/api/monitor/store"
)
