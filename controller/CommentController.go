package controller

import (
	"context"
	"elastic-be/constants"
	"elastic-be/model"
	"elastic-be/model/dto"
	"elastic-be/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
}

func (controller *CommentController) AddComment(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()
	res := model.Response{}

	req := model.Comment{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Request body error:", err)
		res.Rc = constants.ERR_CODE_03
		res.Msg = constants.ERR_CODE_03_MSG
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	reqByte, _ := json.Marshal(req)
	log.Println("req -> ", string(reqByte))

	res = service.InitializeServiceInterface().AddComment(&req)

	ctx.JSON(http.StatusOK, res)
}

func (controller *CommentController) GetByFilterPaging(ctx *gin.Context) {
	fmt.Println(">>> Comment Controoler - Get All <<<")
	parent := context.Background()
	defer parent.Done()

	req := dto.RequesDto{}
	res := model.Response{}

	postId := ctx.Param("postId")
	if postId == "" {
		log.Println("error posting id tidak ada")
		res.Rc = constants.ERR_CODE_02
		res.Msg = constants.ERR_CODE_02_MSG
		ctx.JSON(http.StatusOK, res)
		return
	}

	res = service.InitializeServiceInterface().GetDataCommentByFilterPaging(req, postId)

	ctx.JSON(http.StatusOK, res)

}
