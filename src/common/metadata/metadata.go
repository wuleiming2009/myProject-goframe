package metadata

type Key string

const (
	KeyUserAgent Key = "user_agent"
	KeyReqId         = "req_id"
	UserId         = "user_id"
)

const (
	headerKeyReqId = "X-Reqid"
)
