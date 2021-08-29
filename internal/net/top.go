package net

var ErrCodeTable = map[string]int{
	"state: state is not online":           1,
	"username: name already exists":        2,
	"emailVerify: already verified":        3,
	"emailVerify: wrond verification code": 4,
}

type RespHeader struct {
	Status  string `json:"status" form:"status"`
	ErrCode int    `json:"errCode" form:"errCode"`
}

func (r *RespHeader) Handle(err error) {
	if err == nil {
		if r.Status != "fail" {
			r.Status = "success"
		}
		r.ErrCode = 0
	} else {
		r.Status = "fail"
		ec, ok := ErrCodeTable[err.Error()]
		if ok {
			r.ErrCode = ec
		} else {
			r.ErrCode = 0
		}
	}
}
