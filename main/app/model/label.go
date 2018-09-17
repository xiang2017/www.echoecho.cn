package model


type Label struct {
	BaseModel
	Name 		string
}


func (Label) TableName() string {
	return "labels"
}

func (label *Label) Save() error {
	if label.ID == 0 {
		res := Mdb.Create(label)
		return res.Error
	} else{
		res := Mdb.Save(label)
		return res.Error
	}
}

func (label *Label) Find(id int) error {
	res := Mdb.First(label, id)
	return res.Error
}

func (label *Label) Delete() error {
	res := Mdb.Delete(label)
	return res.Error
}

func (label *Label) FindOrCreate(name string){
	Mdb.Where("name=?", name).First(label)
	if (*label).ID == 0 {
		(*label).Name = name
		if res := Mdb.Create(label); res.Error != nil {
			panic(res.Error)
		}
	}
}

type QuestionLabelRel struct {
	BaseModel
	QuestionId	uint64
	LabelId		uint64
}

func (QuestionLabelRel) TableName() string{
	return "question_label_rel"
}

func GetLabelNamesByRel (relTable string, query interface{}, arg ...interface{}) []string{
	type LabelResult struct {
		Name 		string
	}
	var labels = make([]LabelResult, 10, 10)
	var labelNames = make([]string, 0, 10)
	Mdb.Table(relTable).Select("labels.name").Joins("left join labels on " + relTable + ".label_id = labels.id").Where(query, arg).Find(&labels)
	for i := 0; i < len(labels); i ++ {
		labelNames = append(labelNames, labels[i].Name)
	}
	return labelNames
}