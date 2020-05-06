package service

import (
	"context"
	"elastic-be/constants"
	model "elastic-be/model"
	"elastic-be/model/dto"
	"elastic-be/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic"
)

// ElasticServiceInterface ...
type ElasticServiceInterface struct {
}

func getPost(page int, count int) []model.Post {

	var postings []model.Post
	fromPage := (page - 1) * count
	// totalGetPage := fromPage + count

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	q := elastic.NewBoolQuery()
	// q.Must(elastic.NewRangeQuery("Like").From(0).To(100))
	// q.Must(elastic.NewTermQuery("restoId", 79))

	rows, errES := esclient.Search().Index("emenu_pos").Type("_doc").Query(q).Sort("urut", false).From(fromPage).Size(count).Do(context.Background())
	if errES != nil {
		fmt.Println("err es :", errES.Error())
		return postings
	}

	for _, hit := range rows.Hits.Hits {
		var posting model.Post
		_ = json.Unmarshal(hit.Source, &posting)
		postings = append(postings, posting)
	}

	for _, p := range postings {
		fmt.Printf("Posting found Name: %s,  like:%v, dislike: %v , time : %v \n", p.Post, p.Like, p.Dislike, p.InsertDate)
	}

	return postings

}

func addPosting(post *model.Post) {
	// var esClient elastic.Client
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	urut := utils.MakeTimesLong() + post.RestoID
	posting := model.Post{
		ID:          utils.GenerateUUID(),
		Post:        post.Post,
		Like:        0,
		Dislike:     0,
		Pict:        post.Pict,
		InsertDate:  time.Now(),
		RestoID:     post.RestoID,
		RestoName:   post.RestoName,
		RestoImgURL: post.RestoImgURL,
		Urut:        urut,
	}

	dataJSON, err := json.Marshal(posting)
	js := string(dataJSON)
	esclient.Index().
		Index("emenu_pos").
		BodyJson(js).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println("[Elastic][Insert Posting]Insertion Successful")

}

func addComment(comment *model.Comment) (newid string, msg string) {
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	urut := utils.MakeTimesLong()
	datacomment := model.Comment{
		ID:              utils.GenerateUUID(),
		CommendParentID: comment.CommendParentID,
		AuthorName:      comment.AuthorName,
		AuthorImageURL:  comment.AuthorImageURL,
		Message:         comment.Message,
		InsertDate:      time.Now(),
		PostID:          comment.PostID,
		Urut:            urut,
	}

	dataJSON, err := json.Marshal(datacomment)
	js := string(dataJSON)

	res, errAdd := esclient.Index().
		Index("emenu_comment").
		BodyJson(js).
		Do(context.Background())

	if errAdd != nil {
		// panic(errAdd)
		println("error add to elastic " + errAdd.Error())
		return "0", errAdd.Error()
	}

	fmt.Println("[Elastic][Insert Comment]Insertion Successful")
	return res.Id, "ok"
}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://52.221.255.231:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err
}

func InitializeServiceInterface() *ElasticServiceInterface {

	return &ElasticServiceInterface{}
}

func (service *ElasticServiceInterface) AddPosting(posting *model.Post) model.Response {
	var res model.Response

	addPosting(posting)

	res.Rc = constants.ERR_CODE_00
	res.Msg = constants.ERR_CODE_00_MSG
	res.Data = "tes add "
	res.TotalData = 1

	return res

}

func (service *ElasticServiceInterface) AddComment(comment *model.Comment) model.Response {

	var res model.Response

	elID, er := getElasticIDByPostId(comment.PostID)
	if elID == "0" {
		res.Rc = "99"
		res.Msg = "Post ID not found "
		res.Data = "0"
		res.TotalData = 1
		return res
	}

	print("el id =", elID, " err =", er)

	newID, errMsg := addComment(comment)

	res.Rc = constants.ERR_CODE_00
	res.Msg = constants.ERR_CODE_00_MSG
	if errMsg != "ok" {
		res.Rc = constants.ERR_CODE_10
		res.Msg = constants.ERR_CODE_10_MSG + " " + errMsg
	}
	res.Data = newID
	res.TotalData = 1

	return res

}

func (service *ElasticServiceInterface) GetDataByFilterPaging(req dto.RequesDto, page int, count int) model.Response {
	var res model.Response

	log.Println("get data : ", res)

	res.Rc = constants.ERR_CODE_00
	res.Msg = constants.ERR_CODE_00_MSG
	res.Data = getPost(page, count)
	res.TotalData = 1

	return res

}

func (service *ElasticServiceInterface) GetDataCommentByFilterPaging(req dto.RequesDto, postId string) model.Response {
	var res model.Response

	log.Println("get data : ", res)

	dataCommens := getComment(postId)
	res.Rc = constants.ERR_CODE_00
	res.Msg = constants.ERR_CODE_00_MSG
	res.Data = dataCommens
	res.TotalData = len(dataCommens)

	return res

}

func (service *ElasticServiceInterface) LikePost(postId string) model.Response {
	var res model.Response

	errCode, errMsg := likePost(postId)
	if errCode == "00" {
		res.Rc = constants.ERR_CODE_00
		res.Msg = constants.ERR_CODE_00_MSG
	} else {
		res.Rc = errCode
		res.Msg = "Data not update"
	}

	res.Data = errMsg
	res.TotalData = 1

	return res

}

func likePost(postID string) (code string, msg string) {
	// result := false

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	elID, er := getElasticIDByPostId(postID)
	if elID == "0" {
		return "99", er
	}

	print("el id =", elID, " err =", er)

	ctx := context.Background()

	update, err := esclient.Update().Index("emenu_pos").Type("_doc").Id(elID).
		Script(elastic.NewScript("ctx._source.like += params.num ").Lang("painless").Param("num", 1)).
		Upsert(map[string]interface{}{"like": 0}).
		Do(ctx)

	if err != nil {
		// Handle error
		// panic(err)
		print(err.Error())
		return "99", err.Error()
	}

	// update.ForcedRefresh
	return "00", update.Result
}

func (service *ElasticServiceInterface) DislikePost(postId string) model.Response {
	var res model.Response

	errCode, errMsg := dislikePost(postId)
	if errCode == "00" {
		res.Rc = constants.ERR_CODE_00
		res.Msg = constants.ERR_CODE_00_MSG
	} else {
		res.Rc = errCode
		res.Msg = "Data not update"
	}
	res.Data = errMsg
	res.TotalData = 1

	return res

}

func dislikePost(postID string) (code string, msg string) {
	// result := false

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	elID, er := getElasticIDByPostId(postID)
	if elID == "0" {
		return "99", er
	}

	print("el id =", elID, " err =", er)

	ctx := context.Background()

	update, err := esclient.Update().Index("emenu_pos").Type("_doc").Id(elID).
		Script(elastic.NewScript("ctx._source.dislike += params.num ").Lang("painless").Param("num", 1)).
		Upsert(map[string]interface{}{"dislike": 0}).
		Do(ctx)

	if err != nil {
		// Handle error
		// panic(err)
		print(err.Error())
		return "", ""
	}

	// update.ForcedRefresh
	return "00", update.Result
}

func getElasticIDByPostId(postID string) (string, string) {

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	q := elastic.NewBoolQuery()
	q.Must(elastic.NewTermQuery("id", postID))

	rows, errES := esclient.Search().Index("emenu_pos").Type("_doc").Query(q).Do(context.Background())
	if errES != nil {
		fmt.Println("err es :", errES.Error())
		return "0", errES.Error()
	}

	for _, hit := range rows.Hits.Hits {

		return hit.Id, ""
	}

	return "0", "post not found"

}

func getComment(postId string) []model.Comment {

	var comments []model.Comment

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	q := elastic.NewBoolQuery()
	// q.Must(elastic.NewRangeQuery("Like").From(0).To(100))
	q.Must(elastic.NewTermQuery("postID", postId))

	rows, errES := esclient.Search().Index("emenu_comment").Type("_doc").Query(q).Sort("urut", false).Do(context.Background())
	if errES != nil {
		fmt.Println("err es :", errES.Error())
		return comments
	}

	for _, hit := range rows.Hits.Hits {
		var comment model.Comment
		_ = json.Unmarshal(hit.Source, &comment)
		comments = append(comments, comment)
	}

	for _, c := range comments {
		fmt.Printf("Comment found Msg: %s,  author:%v \n", c.Message, c.AuthorName)
	}

	return comments

}
