// onebot_types.go OneBot协议相关类型定义
// 实现OneBot v11标准的基本数据结构和常量定义
package main

import "encoding/json"

// OneBot协议版本
const (
	OneBotVersion = "11"
	OneBotImpl    = "yijing-onebot"
)

// OneBot事件类型常量
const (
	// 消息事件
	EventTypeMessage = "message"
	EventTypeNotice  = "notice"
	EventTypeRequest = "request"
	EventTypeMeta    = "meta_event"

	// 消息子类型
	MessageTypePrivate = "private"
	MessageTypeGroup   = "group"

	// 元事件子类型
	MetaEventTypeLifecycle = "lifecycle"
	MetaEventTypeHeartbeat = "heartbeat"
)

// OneBot动作类型常量
const (
	ActionSendPrivateMsg       = "send_private_msg"
	ActionSendGroupMsg         = "send_group_msg"
	ActionSendMsg              = "send_msg"
	ActionDeleteMsg            = "delete_msg"
	ActionGetMsg               = "get_msg"
	ActionGetForwardMsg        = "get_forward_msg"
	ActionSendLike             = "send_like"
	ActionSetGroupKick         = "set_group_kick"
	ActionSetGroupBan          = "set_group_ban"
	ActionSetGroupAnonymousBan = "set_group_anonymous_ban"
	ActionSetGroupWholeBan     = "set_group_whole_ban"
	ActionSetGroupAdmin        = "set_group_admin"
	ActionSetGroupAnonymous    = "set_group_anonymous"
	ActionSetGroupCard         = "set_group_card"
	ActionSetGroupName         = "set_group_name"
	ActionSetGroupLeave        = "set_group_leave"
	ActionSetGroupSpecialTitle = "set_group_special_title"
	ActionSetFriendAddRequest  = "set_friend_add_request"
	ActionSetGroupAddRequest   = "set_group_add_request"
	ActionGetLoginInfo         = "get_login_info"
	ActionGetStrangerInfo      = "get_stranger_info"
	ActionGetFriendList        = "get_friend_list"
	ActionGetGroupInfo         = "get_group_info"
	ActionGetGroupList         = "get_group_list"
	ActionGetGroupMemberInfo   = "get_group_member_info"
	ActionGetGroupMemberList   = "get_group_member_list"
	ActionGetGroupHonorInfo    = "get_group_honor_info"
	ActionGetCookies           = "get_cookies"
	ActionGetCsrfToken         = "get_csrf_token"
	ActionGetCredentials       = "get_credentials"
	ActionGetRecord            = "get_record"
	ActionGetImage             = "get_image"
	ActionCanSendImage         = "can_send_image"
	ActionCanSendRecord        = "can_send_record"
	ActionGetStatus            = "get_status"
	ActionGetVersionInfo       = "get_version_info"
	ActionSetRestart           = "set_restart"
	ActionCleanCache           = "clean_cache"
)

// OneBot消息段类型常量
const (
	SegmentTypeText      = "text"
	SegmentTypeFace      = "face"
	SegmentTypeImage     = "image"
	SegmentTypeRecord    = "record"
	SegmentTypeVideo     = "video"
	SegmentTypeAt        = "at"
	SegmentTypeRps       = "rps"
	SegmentTypeDice      = "dice"
	SegmentTypeShake     = "shake"
	SegmentTypePoke      = "poke"
	SegmentTypeAnonymous = "anonymous"
	SegmentTypeShare     = "share"
	SegmentTypeContact   = "contact"
	SegmentTypeLocation  = "location"
	SegmentTypeMusic     = "music"
	SegmentTypeReply     = "reply"
	SegmentTypeForward   = "forward"
	SegmentTypeNode      = "node"
	SegmentTypeXml       = "xml"
	SegmentTypeJson      = "json"
)

// OneBot基础事件结构
type OneBotEvent struct {
	Time        int64       `json:"time"`                   // 事件发生的时间戳
	SelfId      int64       `json:"self_id"`                // 收到事件的机器人QQ号
	PostType    string      `json:"post_type"`              // 上报类型
	MessageType string      `json:"message_type,omitempty"` // 消息类型
	SubType     string      `json:"sub_type,omitempty"`     // 子类型
	MessageId   int32       `json:"message_id,omitempty"`   // 消息ID
	UserId      int64       `json:"user_id,omitempty"`      // 发送者QQ号
	GroupId     int64       `json:"group_id,omitempty"`     // 群号
	Message     interface{} `json:"message,omitempty"`      // 消息内容
	RawMessage  string      `json:"raw_message,omitempty"`  // 原始消息内容
	Font        int32       `json:"font,omitempty"`         // 字体
	Sender      *Sender     `json:"sender,omitempty"`       // 发送人信息
	Anonymous   *Anonymous  `json:"anonymous,omitempty"`    // 匿名信息
}

// OneBot动作请求结构
type OneBotAction struct {
	Action string                 `json:"action"`           // 动作名称
	Params map[string]interface{} `json:"params,omitempty"` // 动作参数
	Echo   interface{}            `json:"echo,omitempty"`   // 回声，用于标识请求
}

// OneBot动作响应结构
type OneBotActionResponse struct {
	Status  string      `json:"status"`            // 状态，success 或 failed
	Retcode int         `json:"retcode"`           // 返回码，0表示成功
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Message string      `json:"message,omitempty"` // 错误消息
	Wording string      `json:"wording,omitempty"` // 错误描述
	Echo    interface{} `json:"echo,omitempty"`    // 回声
}

// 消息段结构
type MessageSegment struct {
	Type string                 `json:"type"` // 段类型
	Data map[string]interface{} `json:"data"` // 段数据
}

// 发送者信息
type Sender struct {
	UserId   int64  `json:"user_id,omitempty"`  // QQ号
	Nickname string `json:"nickname,omitempty"` // 昵称
	Sex      string `json:"sex,omitempty"`      // 性别
	Age      int32  `json:"age,omitempty"`      // 年龄
	Card     string `json:"card,omitempty"`     // 群名片/备注
	Area     string `json:"area,omitempty"`     // 地区
	Level    string `json:"level,omitempty"`    // 成员等级
	Role     string `json:"role,omitempty"`     // 角色
	Title    string `json:"title,omitempty"`    // 专属头衔
}

// 匿名信息
type Anonymous struct {
	Id   int64  `json:"id"`   // 匿名用户ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户flag
}

// 生命周期事件
type LifecycleEvent struct {
	OneBotEvent
	SubType string `json:"sub_type"` // enable, disable, connect
}

// 心跳事件
type HeartbeatEvent struct {
	OneBotEvent
	Status   interface{} `json:"status"`   // 状态信息
	Interval int64       `json:"interval"` // 心跳间隔（毫秒）
}

// 私聊消息事件
type PrivateMessageEvent struct {
	OneBotEvent
	TempSource int32 `json:"temp_source,omitempty"` // 临时会话来源
}

// 群消息事件
type GroupMessageEvent struct {
	OneBotEvent
}

// OneBot元事件
type MetaEvent struct {
	OneBotEvent
	MetaEventType string `json:"meta_event_type"` // 元事件类型
}

// 创建文本消息段
func NewTextSegment(text string) MessageSegment {
	return MessageSegment{
		Type: SegmentTypeText,
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

// 创建图片消息段
func NewImageSegment(file string) MessageSegment {
	return MessageSegment{
		Type: SegmentTypeImage,
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

// 创建AT消息段
func NewAtSegment(qq int64) MessageSegment {
	return MessageSegment{
		Type: SegmentTypeAt,
		Data: map[string]interface{}{
			"qq": qq,
		},
	}
}

// 创建回复消息段
func NewReplySegment(messageId int32) MessageSegment {
	return MessageSegment{
		Type: SegmentTypeReply,
		Data: map[string]interface{}{
			"id": messageId,
		},
	}
}

// 消息段数组类型
type Message []MessageSegment

// 将消息段数组转换为字符串
func (m Message) String() string {
	var result string
	for _, seg := range m {
		switch seg.Type {
		case SegmentTypeText:
			if text, ok := seg.Data["text"].(string); ok {
				result += text
			}
		case SegmentTypeAt:
			if qq, ok := seg.Data["qq"]; ok {
				result += "[CQ:at,qq=" + jsonToString(qq) + "]"
			}
		case SegmentTypeImage:
			if file, ok := seg.Data["file"].(string); ok {
				result += "[CQ:image,file=" + file + "]"
			}
		default:
			result += "[CQ:" + seg.Type + "]"
		}
	}
	return result
}

// 辅助函数：将接口类型转换为字符串
func jsonToString(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// OneBot状态信息
type OneBotStatus struct {
	AppInitialized bool   `json:"app_initialized"` // 原 APP 是否初始化
	AppEnabled     bool   `json:"app_enabled"`     // 原 APP 是否启用
	PluginsGood    bool   `json:"plugins_good"`    // 插件正常
	AppGood        bool   `json:"app_good"`        // 原 APP 正常
	Online         bool   `json:"online"`          // 是否在线
	Good           bool   `json:"good"`            // OneBot 实现是否正常运行
	Stat           *Stats `json:"stat,omitempty"`  // 统计信息
}

// 统计信息
type Stats struct {
	PacketReceived  int64 `json:"packet_received"`  // 收到的数据包数量
	PacketSent      int64 `json:"packet_sent"`      // 发送的数据包数量
	PacketLost      int64 `json:"packet_lost"`      // 丢失的数据包数量
	MessageReceived int64 `json:"message_received"` // 接收消息数
	MessageSent     int64 `json:"message_sent"`     // 发送消息数
	DisconnectTimes int64 `json:"disconnect_times"` // 连接断开次数
	LostTimes       int64 `json:"lost_times"`       // 连接丢失次数
}

// 版本信息
type VersionInfo struct {
	AppName         string `json:"app_name"`         // 应用标识
	AppVersion      string `json:"app_version"`      // 应用版本
	ProtocolVersion string `json:"protocol_version"` // OneBot 标准版本
	OneBotVersion   string `json:"onebot_version"`   // OneBot 实现版本
}

// 登录信息
type LoginInfo struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
}
