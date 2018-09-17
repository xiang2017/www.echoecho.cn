package bgcontroller

import (
	"github.com/gin-gonic/gin"
	"www.echoecho.cn/main/app/model"
	"html"
	"net/http"
	"encoding/json"
	"strconv"
	"www.echoecho.cn/main/app/helper"
)

type OptionMessage struct {
	OptionType 		string	`json:"option_type"`
	Content			string	`json:"content"`
}

type KnowledgeMessage struct {
	ID 				int 	`json:"id"`
}

type QuestionMessage struct {
	ID 				int		`json:"id"`
	QuestionType	string	`json:"question_type"`
	Name			string	`json:"name"`
	BgName			string	`json:"bg_name"`
	Content			string	`json:"content"`
	Options			[]OptionMessage	`json:"options"`
	Answer			json.Number	`json:"answer,Number"`
	AnswerExplain	string	`json:"answer_explain"`
	QuestionDifficulty	string	`json:"question_difficulty"`
	Labels			[]string	`json:"labels"`
	Knowledge 		[]KnowledgeMessage `json:"knowledge"`
}

func QuestionList(c *gin.Context){
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	var questions = make([]model.Question, 15, 15)
	if res := model.Mdb.Offset(page * PER_PAGE).Limit(PER_PAGE).Order("id desc").Find(&questions); res.Error != nil {
		panic(res.Error)
	}

	var count int
	if res := model.Mdb.Table(model.Question{}.TableName()).Count(&count); res.Error != nil {
		panic(res.Error)
	}

	var data = make([]gin.H, 0, 15)
	for i := 0; i < len(questions); i ++ {
		// labels


		// knowledge

		data = append(data, gin.H {
			"id": questions[i].ID,
			"name": questions[i].Name,
			"bg_name": questions[i].BgName,
			"done_count": questions[i].DoneCount,
			"right_count": questions[i].RightCount,
			"question_difficulty": questions[i].QuestionDifficulty,
			"created_at": questions[i].CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"data":  data,
	})
}

// add or edit question
func EditQuestion(c *gin.Context){
	var message QuestionMessage
	if err := c.BindJSON(&message); err != nil {
		panic(err)
	}

	// parse int
	var questionType int
	var err error
	if questionType, err = strconv.Atoi(message.QuestionType); err != nil {
		panic(err)
	}

	// validate
	if message.Name == "" {
		c.String(http.StatusBadRequest, "题目名不能为空")
		return
	}
	if message.BgName == "" {
		c.String(http.StatusBadRequest, "题目后台名不能为空")
		return
	}
	if questionType == model.QuestionTypeChoice && len(message.Options) < 3 {
		c.String(http.StatusBadRequest, "选择题必须至少三个选项")
		return
	}

	var question model.Question

	if message.ID != 0 {
		model.Mdb.First(&question, message.ID)
		if question.ID == 0 {
			c.String(http.StatusNotFound, "问题不存在")
			return
		}
	}

	// set field
	question.Name = message.Name

	question.QuestionType = questionType
	question.BgName = message.BgName
	question.Content = html.EscapeString(message.Content)

	if message.ID != 0 {
		if question.QuestionType != questionType {
			c.String(http.StatusConflict, "问题类型不能修改")
			return
		}
		if question.QuestionType == model.QuestionTypeChoice && question.Answer != string(message.Answer) {
			c.String(http.StatusConflict, "选择题选项不能修改")
			return
		}
	}

	question.Answer = string(message.Answer)
	question.AnswerExplain = message.AnswerExplain
	difficulty, _ := strconv.Atoi(message.QuestionDifficulty)
	question.QuestionDifficulty = byte(difficulty)

	if message.ID != 0 {
		if res := model.Mdb.Save(&question); res.Error != nil {
			panic(res.Error)
		}
	} else {
		if res := model.Mdb.Create(&question); res.Error != nil {
			panic(res.Error)
		}
	}

	// update options
	if question.QuestionType == model.QuestionTypeChoice {
		var options = make([]model.QuestionOption, 4, 4)
		if res := model.Mdb.Where("question_id = ? ", question.ID).Order("option_order asc").Find(&options); res.Error != nil {
			panic(res.Error)
		}
		i := 0
		for ; i < len(options); i ++ {
			if i >= len(message.Options) {
				if res := model.Mdb.Delete(&options[i]); res.Error != nil{
					panic(res.Error)
				}
				continue
			}
			options[i].Content = message.Options[i].Content
			options[i].OptionType = model.QuestionOptionText
			options[i].OptionOrder = byte(i)
			if res := model.Mdb.Save(&options[i]); res.Error != nil {
				panic(res.Error)
			}
		}
		for ; i < len(message.Options); i ++ {
			var option model.QuestionOption
			option.Content = message.Options[i].Content
			option.OptionType = model.QuestionOptionText
			option.OptionOrder = byte(i)
			option.QuestionId = question.ID
			if res := model.Mdb.Create(&option); res.Error != nil {
				panic(res.Error)
			}
		}
	}

	// update labels
	type RelResult struct {
		Name		string
		ID 			int
	}
	var labelRels = make([]RelResult, 5, 5)
	model.Mdb.Table("question_label_rel").Select("labels.name, labels.id").Joins("left join labels on question_label_rel.label_id = labels.id").Where("question_label_rel.question_id=?", question.ID).Find(&labelRels)
	for i := 0; i < len(labelRels); i ++ {
		if !helper.InStringArray(labelRels[i].Name, &message.Labels, true) {
			// delete
			model.Mdb.Where("question_id=? and label_id=?", question.ID, labelRels[i].ID).Delete(model.QuestionLabelRel{})
		} else{
			message.Labels = helper.RemoveStringFromArray(labelRels[i].Name, &message.Labels, true)
		}
	}

	// add Label rel
	for i := 0; i < len(message.Labels); i ++ {
		// 查询是否存在该 label，不存在，则 add
		var label model.Label
		model.Mdb.Where("name=?", message.Labels[i]).First(&label)

		if label.ID == 0 {
			label.Name = message.Labels[i]
			if res := model.Mdb.Create(&label); res.Error != nil {
				panic(res.Error)
			}
		}

		var rel model.QuestionLabelRel
		rel.QuestionId = question.ID
		rel.LabelId = label.ID
		if res := model.Mdb.Create(&rel); res.Error != nil {
			panic(res.Error)
		}
	}

	// update knowledge rel
	var ids []int
	for i := 0; i < len(message.Knowledge); i ++ {
		ids = append(ids, message.Knowledge[i].ID)
	}
	var rels = make([]model.QuestionKnowledgeRel, 1, 1)
	model.Mdb.Where("question_id=?", question.ID).Find(&rels)
	for i := 0; i < len(rels); i ++{
		if helper.InIntArray(rels[i].KnowledgeId, &ids) {
			ids = helper.RemoveIndexFromArray(i, &ids)
		} else{
			model.Mdb.Delete(&rels[i])
		}
	}
	// 最后把 ids 中的数据跟 question 关联
	for i := 0; i < len(ids); i ++ {
		var rel model.QuestionKnowledgeRel
		rel.KnowledgeId = ids[i]
		rel.QuestionId = int(question.ID)
		if res := model.Mdb.Create(&rel); res.Error != nil {
			c.String(http.StatusInternalServerError, res.Error.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": question.ID,
	})
}

// 获取问题信息
func GetQuestionInfo(c *gin.Context) {
	id := c.Query("id")
	if "" == id {
		c.String(http.StatusBadRequest, "id can not be nil")
		return
	}

	var question model.Question
	if res := model.Mdb.Where("id = ?", id).First(&question); res.Error != nil {
		c.String(http.StatusNotFound, res.Error.Error())
		return
	}

	// 查询问题的标签
	type LabelResult struct {
		Name 		string
	}
	var labels = make([]LabelResult, 10, 10)
	var labelNames = make([]string, 0, 10)
	model.Mdb.Table("question_label_rel").Select("labels.name").Joins("left join labels on question_label_rel.label_id = labels.id").Where("question_label_rel.question_id=?", question.ID).Find(&labels)
	for i := 0; i < len(labels); i ++ {
		labelNames = append(labelNames, labels[i].Name)
	}

	type KnowledgeResult struct {
		ID 			int
		Name 		string
	}
	var knowledgeList = make([]KnowledgeResult, 0, 3)
	model.Mdb.Table("question_knowledge_rel").Select("knowledge.*").Joins("left join knowledge on knowledge.id=question_knowledge_rel.knowledge_id").Where("question_knowledge_rel.question_id=?", question.ID).Find(&knowledgeList)
	var knowledge = make([]gin.H, 0, 2)
	for i := 0; i < len(knowledgeList); i ++ {
		knowledge = append(knowledge, gin.H{
			"id": knowledgeList[i].ID,
			"name": knowledgeList[i].Name,
		})
	}

	var data = gin.H{
		"id": question.ID,
		"name": question.Name,
		"bg_name": question.BgName,
		"question_type": strconv.Itoa(question.QuestionType),
		"content": html.UnescapeString(question.Content),
		"answer": question.Answer,
		"answer_explain": question.AnswerExplain,
		"question_difficulty": strconv.Itoa(int(question.QuestionDifficulty)),
		"labels": labelNames,
		"knowledge": knowledge,
	}

	// 选择题的选项
	if question.QuestionType == model.QuestionTypeChoice {
		var questionOptions = make([]model.QuestionOption, 4, 4)
		model.Mdb.Where("question_id=?", question.ID).Order("option_order asc").Find(&questionOptions)
		var options = make([]gin.H, 0, 0)
		for i := 0; i < len(questionOptions); i ++ {
			options = append(options, gin.H{
				"option_type": strconv.Itoa(int(questionOptions[i].OptionType)),
				"content": questionOptions[i].Content,
			})
		}
		data["options"] = options
		data["answer"], _ = strconv.Atoi(question.Answer)
	} else{
		data["options"] = make([]interface{}, 0, 0)
	}

	c.JSON(http.StatusOK, data)
}

// 删除问题
func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	var question model.Question
	model.Mdb.Where("id = ?", id).First(&question)



	if question.ID == 0 {
		c.String(http.StatusNotFound, "问题不存在")
		return
	}


	// 删除关联 label
	model.Mdb.Where("question_id=?", id).Delete(model.QuestionLabelRel{})

	// 如果是选择题 删除 options
	if question.QuestionType == model.QuestionTypeChoice {
		model.Mdb.Where("question_id=?", id).Delete(model.QuestionOption{})
	}

	model.Mdb.Delete(&question)

	c.String(http.StatusOK, "success")
}
