package hnjobs

import (
	"encoding/json"
	"fmt"
	"hnjobs/services/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	jobsURL = "https://hacker-news.firebaseio.com/v0/item/%s.json?"
)

type Response struct {
	ID       int          `json:"id" xorm:"integer '_id'"`
	ParentID int          `json:"parent" xorm:"integer '_parent'"`
	By       string       `json:"by" xorm:"varchar(32) '_by'"`
	Text     string       `json:"text" xorm:"varchar '_text'"`
	Time     JSONDateTime `json:"time" xorm:"timestamp '_time'"`
	Type     string       `json:"type" xorm:"varchar(12) '_type'"`
	Kids     []int        `json:"kids"`
}

func (r Response) TableName() string {
	return "hn_jobs"
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func fetchUpdates(story *Response) []int {
	engine := db.GetEngine()

	query := fmt.Sprintf(`SELECT * FROM unnest(array[%s]) __id WHERE
		NOT EXISTS (SELECT 1 FROM hn_jobs WHERE _id = __id AND _parent = %d)`,
		arrayToString(story.Kids, ","), story.ID)
	res, err := engine.Query(query)
	if err != nil {
		fmt.Printf("ERROR fetching kid info:: %s\n", err.Error())
		return nil
	}

	kids := make([]int, 0, len(res))
	for _, m := range res {
		id, _ := strconv.ParseUint(string(m["__id"]), 10, 64)
		kids = append(kids, int(id))
	}
	// log.Printf("RES :: %+v \n", kids)
	return kids
}

func FetchStory(id string) {
	story, err := fetch(id)
	if err != nil {
		log.Printf("Error fetching story:: %s\n", err.Error())
		return
	}

	kids := fetchUpdates(story)
	if len(kids) == 0 {
		return
	}

	var res *Response
	posts := make([]*Response, 0, len(kids)+1)
	// If returned list is same as story.Kids, that means it is new story (None of the kids are present in db),
	// so add the story as well.
	if len(kids) == len(story.Kids) {
		posts = append(posts, story)
	}
	for _, kidID := range kids {
		time.Sleep(time.Millisecond * 750)
		if res, err = fetch(strconv.FormatInt(int64(kidID), 10)); err == nil {
			posts = append(posts, res)
		} else {
			log.Printf("Error fetching comment:: %s\n", err.Error())
		}
	}

	engine := db.GetEngine()
	affected, err := engine.Omit("kids").Insert(posts)
	log.Printf("Affected: %d, Error: %+v", affected, err)
}

func fetch(id string) (*Response, error) {
	log.Printf("Fetching id %s\n", id)
	request, err := http.NewRequest("GET", fmt.Sprintf(jobsURL, id), nil)
	if err != nil {
		return nil, err
	}

	return executeQuery(request)
}

func executeQuery(request *http.Request) (*Response, error) {
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var res Response
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

type JSONDateTime time.Time

func (j *JSONDateTime) UnmarshalJSON(b []byte) error {
	sec, _ := strconv.ParseUint(string(b), 10, 64)
	t := time.Unix(int64(sec), 0)
	*j = JSONDateTime(t)
	return nil
}

func (j *JSONDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*j).UTC())
}
