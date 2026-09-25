package constants

const (
	MsgOK            = "ok"
	MsgInvalidParams = "请求参数错误"
	MsgUnauthorized  = "未登录或登录已过期"
	MsgForbidden     = "没有操作权限"
	MsgNotFound      = "资源不存在"
	MsgInternalError = "服务器内部错误"
	MsgLoginSuccess  = "登录成功"
	MsgRateLimited   = "请求过于频繁，请稍后再试"
)

// 展览就绪（readiness）文案。
const (
	MsgReadinessReady      = "全部关联作品均可公开展出，展览已就绪。"
	MsgReadinessBlocked    = "存在不可公开展出的关联作品，请处理阻断原因后再公开发布。"
	MsgReadinessEmpty      = "本展尚未收录任何作品，请先添加作品。"
	MsgReadinessEmptyBrief = "空展览暂不能公开展出"
	MsgReadinessMissingFmt = "作品[id=%s]在作品库中不存在或已被删除，关联引用失效。"
	MsgReadinessDraftFmt   = "作品[id=%s]仍为草稿（status=Draft），未发布不能公开展出。"
	MsgReadinessPendingFmt = "作品[id=%s]审核状态为待审核（reviewStatus=Pending），审核通过前不能公开。"
	MsgReadinessRejectFmt  = "作品[id=%s]审核未通过（reviewStatus=Rejected），已被驳回不能公开。"
	MsgReadinessFlagFmt    = "作品[id=%s]被审核标记（reviewStatus=Flagged），需复核后才能公开。"
	MsgReadinessRevUnkFmt  = "作品[id=%s]审核状态异常（reviewStatus=%s），无法确认能否公开。"
	MsgReadinessSoldFmt    = "作品[id=%s]状态为已售（status=Sold），不能继续公开展出。"
	MsgReadinessArchFmt    = "作品[id=%s]已下架（status=Archived），不能继续公开展出。"
	MsgReadinessStsUnkFmt  = "作品[id=%s]作品状态异常（status=%s），无法确认能否公开。"
	MsgReadinessNoImgFmt   = "作品[id=%s]缺少可公开展示的图片（imageUrls 为空）。"
)
