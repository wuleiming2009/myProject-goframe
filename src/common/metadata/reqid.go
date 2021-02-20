package metadata

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"os"
	"time"

	"github.com/spaolacci/murmur3"
)

var pid = uint16(time.Now().UnixNano() & 65535)
var machineFlag uint16

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	machineFlag = uint16(murmur3.Sum32([]byte(hostname)) & 65535)
}

func GenReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint16(b[:], pid)
	binary.LittleEndian.PutUint16(b[2:], machineFlag)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func NewContextWithReqId(ctx context.Context, reqId ...string) context.Context {
	if reqId == nil {
		reqId = []string{GenReqId()}
	}
	return context.WithValue(ctx, KeyReqId, reqId[0])
}

func ReqIdFromContext(ctx context.Context) string {
	if reqId, ok := ctx.Value(KeyReqId).(string); ok {
		return reqId
	}
	return ""
}

func NewContextWithUserId(ctx context.Context, userId uint64) context.Context {
	return context.WithValue(ctx, UserId, userId)
}

func UserIdFromContext(ctx context.Context) uint64 {
	if userId, ok := ctx.Value(UserId).(uint64); ok {
		return userId
	}
	return 0
}
