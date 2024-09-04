package schemas

type URI struct {
	UserId        int64 `uri:"userId"`
	PostId        int64 `uri:"postId"`
	PostCommentId int64 `uri:"postCommentId"`
}
