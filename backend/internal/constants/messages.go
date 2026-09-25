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

	// 展览公开就绪评估的阻断原因文案（同时用于后端返回与前端展示）。
	MsgReadinessEmptyExhibition  = "展览尚未收录任何作品，暂无可公开内容"
	MsgReadinessMissingArtworks  = "%d 件关联作品不存在或已被删除"
	MsgReadinessDraftArtworks    = "%d 件作品尚未发布"
	MsgReadinessPendingReview    = "%d 件作品仍在等待审核"
	MsgReadinessRejectedReview   = "%d 件作品审核未通过"
	MsgReadinessSoldArtworks     = "%d 件作品已售出"
	MsgReadinessArchivedArtworks = "%d 件作品已下架"
	MsgReadinessNoImage          = "%d 件作品缺少图片"
)
