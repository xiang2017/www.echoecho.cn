package bgcontroller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"www.echoecho.cn/main/app/model"
	"html"
	"www.echoecho.cn/main/app/helper"
	"strconv"
)


func EditKnowledge(c *gin.Context) {

	type KnowledgeMessage struct {
		ID			int		`json:"id"`
		Name 		string	`json:"name"`
		Content		string 	`json:"content"`
		Labels 		[]string `json:"labels"`
	}

	var message KnowledgeMessage
	if err := c.BindJSON(&message); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if message.Name == "" {
		c.String(http.StatusBadRequest, "名字不能为空")
		return
	}


	if message.Content == "" {
		c.String(http.StatusBadRequest, "内容不能为空")
		return
	}

	var knowledge model.Knowledge
	if message.ID != 0 {
		model.Mdb.Where("id=?", message.ID).First(&knowledge)
		if knowledge.ID == 0 {
			c.String(http.StatusNotFound, "knowledge 不存在")
			return
		}
	}

	knowledge.Name = message.Name
	knowledge.Content = html.EscapeString(message.Content)

	if message.ID != 0 {
		if res := model.Mdb.Save(&knowledge); res.Error != nil {
			panic(res.Error)
		}
	} else{
		if res := model.Mdb.Create(&knowledge); res.Error != nil {
			panic(res.Error)
		}
	}

	// labels
	type RelResult struct {
		Name		string
		ID 			int
	}
	var labelRel = make([]RelResult, 15, 15)
	model.Mdb.Table("knowledge_label_rel").Select("labels.*").Joins("left join labels on labels.id=knowledge_label_rel.label_id").Where("knowledge_label_rel.knowledge_id=?", knowledge.ID).Find(&labelRel)
	for i := 0; i < len(labelRel); i ++ {
		if !helper.InStringArray(labelRel[i].Name, &message.Labels, true) {
			// delete
			model.Mdb.Where("knowledge_id=? and label_id=?", knowledge.ID, labelRel[i].ID).Delete(model.KnowledgeLabelRel{})
		} else{
			message.Labels = helper.RemoveStringFromArray(labelRel[i].Name, &message.Labels, true)
		}
	}
	// add label rel
	for i := 0; i < len(message.Labels); i ++ {
		// 查询是否存在该 label，不存在，则 add
		var label model.Label
		label.FindOrCreate(message.Labels[i])

		var rel model.KnowledgeLabelRel
		rel.KnowledgeId = knowledge.ID
		rel.LabelId = label.ID
		if res := model.Mdb.Create(&rel); res.Error != nil {
			panic(res.Error)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": knowledge.ID,
	})
}


func KnowledgeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	var knowledgeList = make([]model.Knowledge, 15, 15)
	if res := model.Mdb.Offset(page * PER_PAGE).Limit(PER_PAGE).Find(&knowledgeList); res.Error != nil {
		panic(res.Error)
	}

	var count int
	if res := model.Mdb.Table(model.Knowledge{}.TableName()).Count(&count); res.Error != nil {
		panic(res.Error)
	}

	var data = make([]gin.H, 0, 15)
	for i := 0; i < len(knowledgeList); i ++ {
		data = append(data, gin.H{
			"id": knowledgeList[i].ID,
			"name": knowledgeList[i].Name,
			"created_at": knowledgeList[i].CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"count": count,
	})
}


func GetKnowledge(c *gin.Context) {
	id := c.Param("id")
	var knowledge model.Knowledge
	if res := model.Mdb.Where("id=?", id).First(&knowledge); res.Error != nil {
		c.String(http.StatusNotFound, res.Error.Error())
		return
	}

	// labels
	labels := model.GetLabelNamesByRel("knowledge_label_rel", "knowledge_label_rel.knowledge_id=?", knowledge.ID)

	var data = gin.H{
		"id": knowledge.ID,
		"name": knowledge.Name,
		"content": html.UnescapeString(knowledge.Content),
		"labels": labels,
	}

	c.JSON(http.StatusOK, data)
}