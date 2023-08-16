package qiniu

import "encoding/json"

func ParsePfopResponse(reqBody []byte) (*PfopResponse, error) {
	var resp PfopResponse
	if err := json.Unmarshal(reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
