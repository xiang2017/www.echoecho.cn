package model


type Knowledge struct {
	BaseModel
	Name 		string
	Content		string
}

func (Knowledge) TableName() string {
	return "knowledge"
}

type KnowledgeLabelRel struct {
	BaseModel
	KnowledgeId		uint64
	LabelId			uint64
}

func (KnowledgeLabelRel) TableName() string {
	return "knowledge_label_rel"
}