package main

import (
	"elastic-be/routers"
	"elastic-be/utils"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

var (
	port string
)

func main() {
	maxProc, _ := strconv.Atoi(utils.GetEnv("MAXPROCS", "1"))
	port = utils.GetEnv("PORT_ELASTIC", "8200")
	runtime.GOMAXPROCS(maxProc)

	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", port)

	log.Println("[info] start http server listening %s", endPoint)

	//server.ListenAndServe()

	routersInit.Run(":" + port)

}

// func sample() {
// 	ctx := context.Background()
// 	esclient, err := GetESClient()
// 	if err != nil {
// 		fmt.Println("Error initializing : ", err)
// 		panic("Client fail ")
// 	}

// posting := Posting{
// 	Post:       "posting from go 8",
// 	Like:       8,
// 	Dislike:    88,
// 	Pict:       "",
// 	InsertDate: time.Now(),
// 	RestoId:    79,
// }

// dataJSON, err := json.Marshal(posting)
// js := string(dataJSON)
// esclient.Index().
// 	Index("emenu_pos").
// 	BodyJson(js).
// 	Do(ctx)

// if err != nil {
// 	panic(err)
// }

// fmt.Println("[Elastic][InsertProduct]Insertion Successful")

// var postings []Posting

// searchSource := elastic.NewSearchSource()
// searchSource.Query(elastic.NewMatchQuery("restoId", 79))

// /* this block will basically print out the es query */
// queryStr, err1 := searchSource.Source()
// queryJs, err2 := json.Marshal(queryStr)

// if err1 != nil || err2 != nil {
// 	fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
// }
// fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
// /* until this block */

// searchService := esclient.Search().Index("emenu_pos").Sort("Post.keyword", true).From(0).Size(3).Pretty(true).SearchSource(searchSource)

// searchResult, err := searchService.Do(ctx)
// if err != nil {
// 	fmt.Println("[ProductsES][GetPIds]Error=", err)
// 	return
// }

// for _, hit := range searchResult.Hits.Hits {
// 	var posting Posting
// 	err := json.Unmarshal(hit.Source, &posting)
// 	if err != nil {
// 		fmt.Println("[Getting Posting][Unmarshal] Err=", err)
// 	}

// 	postings = append(postings, posting)
// }

// if err != nil {
// 	fmt.Println("Fetching posting fail: ", err)
// } else {
// 	for _, p := range postings {
// 		fmt.Printf("Posting found Name: %s,  like:%v, dislike: %v \n", p.Post, p.Like, p.Dislike)
// 	}
// }

// 	q := elastic.NewBoolQuery()
// 	q.Must(elastic.NewRangeQuery("Like").From(0).To(100))
// 	q.Must(elastic.NewTermQuery("restoId", 79))

// 	rows, errES := esclient.Search().Index("emenu_pos").Type("_doc").Query(q).Sort("InsertDate", false).From(0).Size(10).Do(ctx)
// 	if errES != nil {
// 		fmt.Println("err es :", errES.Error())
// 		return
// 	}

// 	var postings []model.Posting
// 	for _, hit := range rows.Hits.Hits {
// 		var posting model.Posting
// 		_ = json.Unmarshal(hit.Source, &posting)
// 		postings = append(postings, posting)
// 	}

// 	for _, p := range postings {
// 		fmt.Printf("Posting found Name: %s,  like:%v, dislike: %v , time : %v \n", p.Post, p.Like, p.Dislike, p.InsertDate)
// 	}

// }
