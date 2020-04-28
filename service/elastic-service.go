package service

import (
	"context"
	"elastic-be/constants"
	"elastic-be/model"
	"elastic-be/model/dto"
	"elastic-be/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic"
)

type ElasticServiceInterface struct {
}

func getPost(page int, count int) []model.Posting {

	var postings []model.Posting
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
		var posting model.Posting
		_ = json.Unmarshal(hit.Source, &posting)
		postings = append(postings, posting)
	}

	for _, p := range postings {
		fmt.Printf("Posting found Name: %s,  like:%v, dislike: %v , time : %v \n", p.Post, p.Like, p.Dislike, p.InsertDate)
	}

	return postings

}
func addPosting(post *model.Posting) {
	// var esClient elastic.Client
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	urut := utils.MakeTimesLong() + post.RestoId
	posting := model.Posting{
		Post:       post.Post,
		Like:       0,
		Dislike:    0,
		Pict:       post.Pict,
		InsertDate: time.Now(),
		RestoId:    post.RestoId,
		Urut:       urut,
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

	fmt.Println("[Elastic][InsertProduct]Insertion Successful")

}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err
}

func InitializeServiceInterface() *ElasticServiceInterface {

	return &ElasticServiceInterface{}
}

func (service *ElasticServiceInterface) AddPosting(posting *model.Posting) model.Response {
	var res model.Response

	addPosting(posting)

	res.Rc = constants.ERR_CODE_00
	res.Msg = constants.ERR_CODE_00_MSG
	res.Data = "tes add "
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
