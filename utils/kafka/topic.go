package kafka

const (
	TOPIC_BI_PROMOTION_RANK           = "bi-rank-promotion-event"
	TOPIC_RANK_EVENT                  = "rank-event"
	TOPIC_COHORT_INITIALIZATION_EVENT = "cohort-initialization-event" // 新用户注册推送
	TOPIC_TASK_UPDATE_EVENT           = "task-update-event"           // 任务完成
	TOPIC_PROP_EVENT                  = "prop-event"                  // 发送道具
	TOPIC_LOGIN_EVENT                 = "login-event"
	TOPIC_PURCHASE_EVENT              = "purchase-event"
	TOPIC_PLATFORM_EVENT              = "platform-event"
	TOPIC_SMS_TRANSACTIONAL_EVENT     = "sms-transactional-event"
	TOPIC_SMS_PROMOTIONAL_EVENT       = "sms-promotional-event"
)

const (
	TOPIC_TASK_PROGRESSES_EVENT = "task-progresses-event" // 任务进度
)

type TaskUpdate struct {
	TaskType             int64  `json:"task_type"`               // 任务类型：0每日任务 1新手任务 2周末任务 3city chart
	Message              string `json:"message"`                 // 消息文本
	TaskId               int64  `json:"task_id"`                 // 任务ID（任务分类的）
	TaskItemId           int64  `json:"task_item_id"`            // 具体任务ID
	TaskItemProgressesId int64  `json:"task_item_progresses_id"` // 任务进度
}

type SendProp struct {
	SendUserId    int64 `json:"send_user_id"`    // 任务类型：0每日任务 1新手任务 2周末任务 3city chart
	ReceiveUserId int64 `json:"receive_user_id"` // 消息文本
	PropId        int64 `json:"prop_id"`         // 道具ID
	TrackType     int64 `json:"track_type"`      // 动画轨迹类型 0无 1直线 2抛物线 3螺旋
}

type BuyErrorInfo struct {
	PurchaseKey string  `json:"purchase_key"`
	Price       float64 `json:"price"`
	Coins       int64   `json:"coins"`
}
