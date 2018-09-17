package model


/* question */

type Question struct {
	BaseModel
	QuestionType		int
	Name				string
	BgName				string
	Content				string
	Answer				string
	AnswerExplain		string
	DoneCount			int
	RightCount			int
	QuestionDifficulty	byte
}

const (
	_ = iota
	QuestionTypeChoice
	QuestionTypeInput
)

func (Question) TableName() string {
	return "questions"
}

/* question option */

type QuestionOption struct {
	BaseModel
	QuestionId			uint64
	OptionType			byte
	Content				string
	OptionOrder			byte
}

const (
	_ = iota
	QuestionOptionText
	QuestionOptionImage
)

func (QuestionOption) TableName() string {
	return "question_options"
}

type QuestionKnowledgeRel struct {
	BaseModel
	QuestionId			int
	KnowledgeId			int
}

func (QuestionKnowledgeRel) TableName() string {
	return "question_knowledge_rel"
}
