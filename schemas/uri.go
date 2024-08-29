package schemas

type URI struct {
	UserId            int64 `uri:"userId"`
	TestId            int64 `uri:"testId"`
	QuestionId        int64 `uri:"questionId"`
	TestQuestionaryId int64 `uri:"testQuestionaryId"`
}
