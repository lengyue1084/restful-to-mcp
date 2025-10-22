package _const

import "errors"

type (
	AuthStatus          int8
	AuthTypeStatus      int8
	HaveToolsStatus     int8
	McpServerTypeStatus int8
	Status              int8
	ServerStatus        int8
	CommonStatus        int8
	SourceType          int8
)

const (
	IsAuthNo  AuthStatus = 1
	IsAuthYes AuthStatus = 2

	// IsAuthNoAuth 0未知，1 不开启，2开启service授权，3开启平台授权，4开启所有的授权'
	IsAuthNoAuth       AuthTypeStatus = 1
	IsAuthServiceAuth  AuthTypeStatus = 2
	IsAuthPlatformAuth AuthTypeStatus = 3
	IsAuthAllAuth      AuthTypeStatus = 4

	HaveToolsNo          HaveToolsStatus     = 1
	HaveToolsYes         HaveToolsStatus     = 2
	McpServerTypeOpenapi McpServerTypeStatus = 1
	McpServerTypeGrpc    McpServerTypeStatus = 2

	// StatusHidden 通用的是否显示
	StatusHidden    Status       = 1
	StatusDisplay   Status       = 2
	CommonStatusNo  CommonStatus = 1
	CommonStatusYes CommonStatus = 2

	ServerNotSetToken    ServerStatus = 1
	ServerHadSetToken    ServerStatus = 2
	ServerTokenIsWorking ServerStatus = 3

	SourceTypeFile SourceType = 1
	SourceTypeForm SourceType = 2
)

const (
	McpServerRefreshTicketTime = 60 * 10
	CommonContextTimeOut       = 5
)

const (
	PlatformToken = "x-mcp-platform-token"
	ServiceToken  = "x-mcp-service-token"
)
const (
	TraceId     = "traceId"
	SpanId      = "spanId"
	ServerToken = "serverToken"
)

func (s ServerStatus) Sting() string {
	switch s {
	case ServerNotSetToken:
		return "未设置"
	case ServerHadSetToken:
		return "已设置"
	case ServerTokenIsWorking:
		return "正常"
	default:
		return "未知"
	}

}

func (s Status) String() string {

	switch s {
	case StatusDisplay:
		return "显示"
	case StatusHidden:
		return "隐藏"
	default:
		return "未知"

	}
}

func (a AuthStatus) String() string {
	switch a {
	case IsAuthNo:
		return "否"
	case IsAuthYes:
		return "是"
	default:
		return "未知"
	}
}
func (h HaveToolsStatus) String() string {
	switch h {
	case HaveToolsNo:
		return "否"
	case HaveToolsYes:
		return "是"
	default:
		return "未知"
	}

}

func (m McpServerTypeStatus) String() string {
	switch m {
	case McpServerTypeOpenapi:
		return "OpenAPI"
	case McpServerTypeGrpc:
		return "gRPC"
	default:
		return "未知"
	}
}

type contextKey string

const (
	TxKey            contextKey = "tx"
	CommonRetryTimes            = 5
)

var (
	McpFileIsNotExist = errors.New("mcp file is not exist")
)
