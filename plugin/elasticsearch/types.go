package elasticsearch

type QueryResponse struct {
	Took         int                    `json:"took"`
	TimeOut      bool                   `json:"time_out"`
	Shards       *Shards                `json:"_shards"`
	Hits         *Hits                  `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits struct {
	Total    *Total     `json:"total"`
	MaxScore float64    `json:"max_score"`
	Hits     []*SubHits `json:"hits"`
}

type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type SubHits struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	Id     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}
