package bgcontroller

import (
	"github.com/gin-gonic/gin"
	"www.echoecho.cn/main/app/model"
	"net/http"
	"strconv"
)


func GetLabelList(c * gin.Context){
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	countPerPage, _ := strconv.Atoi(c.DefaultQuery("count_per_page", strconv.Itoa(PER_PAGE)))

	var labels = make([]model.Label, 15, 15)
	model.Mdb.Offset(page * countPerPage).Limit(countPerPage).Order("id DESC",).Find(&labels)

	var data = make([]gin.H, 0, 15)
	for i := 0; i < len(labels); i ++ {

		data = append(data, gin.H{
			"id" : labels[i].ID,
			"name" : labels[i].Name,
			"created_at" : labels[i].CreatedAt,
		})
	}
	c.JSON(http.StatusOK, data)
}

// 通过 id 获取 label
func GetLabel(c * gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var label model.Label
	label.Find(id)
	if label.ID == 0 {
		c.JSON(http.StatusNotFound, "label 不存在")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": label.ID,
		"name": label.Name,
	})
}

type LabelMessage struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
}

// 新增 修改 label
func EditLabel(c * gin.Context) {
	var req LabelMessage

	if err := c.BindJSON(&req); err != nil {
		panic(err)
	}

	var label model.Label

	if req.ID == 0 {
		label.Name = req.Name
		label.Save()
	} else{
		label.Find(req.ID)
		if label.ID == 0 {
			c.String(http.StatusNotFound, "label 不存在")
			return
		}
		label.Name = req.Name
		label.Save()
	}

	c.JSON(http.StatusOK, label)
}

// 删除 label
func DeleteLabel(c * gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var label model.Label
	label.Find(id)
	if label.ID == 0 {
		c.String(http.StatusNotFound, "label 不存在")
		return
	}

	label.Delete()
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}