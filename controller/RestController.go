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

	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

type RestController struct {
}

func (controller *RestController) AddPosting(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()
	res := model.Response{}

	req := model.Post{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Request body error:", err)
		res.Rc = constants.ERR_CODE_03
		res.Msg = constants.ERR_CODE_03_MSG
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	reqByte, _ := json.Marshal(req)
	log.Println("req -> ", string(reqByte))

	res = service.InitializeServiceInterface().AddPosting(&req)

	ctx.JSON(http.StatusOK, res)
}

func (controller *RestController) GetByFilterPaging(ctx *gin.Context) {
	fmt.Println(">>> Posting Controoler - Get All <<<")
	parent := context.Background()
	defer parent.Done()

	req := dto.RequesDto{}
	res := model.Response{}

	page, errPage := strconv.Atoi(ctx.Param("page"))
	if errPage != nil {
		log.Println("error", errPage)
		res.Rc = constants.ERR_CODE_02
		res.Msg = constants.ERR_CODE_02_MSG
		ctx.JSON(http.StatusOK, res)
		return
	}

	count, errCount := strconv.Atoi(ctx.Param("count"))
	if errCount != nil {
		logs.Info("error", errPage)
		res.Rc = constants.ERR_CODE_02
		res.Msg = constants.ERR_CODE_02_MSG
		ctx.JSON(http.StatusOK, res)
		return
	}

	res = service.InitializeServiceInterface().GetDataByFilterPaging(req, page, count)

	ctx.JSON(http.StatusOK, res)

}

func (controller *RestController) PostLike(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	// req := dto.RequesDto{}
	// res := model.Response{}

	postID := ctx.Param("postId")

	res := service.InitializeServiceInterface().LikePost(postID)

	ctx.JSON(http.StatusOK, res)

}

func (controller *RestController) PostDislike(ctx *gin.Context) {
	parent := context.Background()
	defer parent.Done()

	// req := dto.RequesDto{}
	// res := model.Response{}

	postID := ctx.Param("postId")

	res := service.InitializeServiceInterface().DislikePost(postID)

	ctx.JSON(http.StatusOK, res)

}
