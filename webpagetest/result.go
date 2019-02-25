package webpagetest

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ResultData holds all info about test
type ResultData struct {
	Id         string
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
	Data       string `json:"data"`
}

// GetTestResult returns result of test with testID
// GetTestResult returns result of test with testID
func (c *Client) GetTestResult(testID string) (*ResultData, error) {
	query := url.Values{}
	query.Add("test", testID)
	query.Add("requests", "0")
	query.Add("average", "0")
	query.Add("standard", "0")

	body, err := c.query("/jsonResult.php", query)
	if err != nil {
		return nil, err
	}
	var responose struct {
		StatusCode int             `json:"statusCode"`
		StatusText string          `json:"statusText"`
		Data       json.RawMessage `json:"data"`
	}

	if err = json.Unmarshal(body, &responose); err != nil {
		return nil, err
	}

	if responose.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status %d: %v",
			responose.StatusCode, responose.StatusText)
	}

	data, err := json.Marshal(&responose.Data)
	if err != nil {
		return nil, err
	}

	var resultData = ResultData{
		Id:         testID,
		StatusCode: responose.StatusCode,
		StatusText: responose.StatusText,
		Data:       string(data),
	}

	return &resultData, nil
}
