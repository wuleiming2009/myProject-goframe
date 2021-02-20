package log

type key int

var (
	zapKey   key = 2
	dlKey    key = 1
	reqIdKey key = 3
)

func GetZapKey() key {
	return zapKey
}

func GetDlKey() key {
	return dlKey
}

func GetReqIdKey() key {
	return reqIdKey
}
