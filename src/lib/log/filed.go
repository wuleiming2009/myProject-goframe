package log

import "go.uber.org/zap"

//http://km.innotechx.com/pages/viewpage.action?pageId=89041696 公用字段定义
var (
	requestId  = "request_id"
	machine    = "machine"
	app        = "app"
	method     = "method"
	node       = "node"
	logType    = "type"
	content    = "content"
	deviceID   = "device_id"
	uid        = "uid"
	requestURI = "request_uri"
	pageID     = "page_id"
	args       = "args"
)

func Field(key string, value interface{}) zap.Field {
	return zap.Reflect(key, value)
}

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int32(key string, value int32) zap.Field {
	return zap.Int32(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Uint64(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

func Float32(key string, value float32) zap.Field {
	return zap.Float32(key, value)
}

func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

func ErrorField(err error) zap.Field {
	if err != nil {
		return zap.String("error", err.Error())
	}
	return zap.Skip()
}

func ReqId(reqid string) zap.Field {
	if len(reqid) == 0 {
		return zap.Skip()
	}

	return zap.String(requestId, reqid)
}

func MachineField(name string) zap.Field {
	return zap.String(machine, name)
}

func AppIdField(appName string) zap.Field {
	return zap.String(app, appName)
}

func MethodField(methodName string) zap.Field {
	return zap.String(method, methodName)
}

func NodeField(nodeName string) zap.Field {
	return zap.String(node, nodeName)
}

func LogTypeField(_type string) zap.Field {
	return zap.String(logType, _type)
}

func ContentField(contents string) zap.Field {
	return zap.String(content, contents)
}

func DeviceIdField(deviceId string) zap.Field {
	return zap.String(deviceID, deviceId)
}

func UidField(userId string) zap.Field {
	return zap.String(uid, userId)
}

func RequestUrlField(reqUrl string) zap.Field {
	return zap.String(requestURI, reqUrl)
}

func PageIdField(pageId string) zap.Field {
	return zap.String(pageID, pageId)
}

func ArgsField(arg interface{}) zap.Field {
	return zap.Reflect(args, arg)
}
