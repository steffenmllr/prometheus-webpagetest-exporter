// adapted from https://github.com/olegfedoseev/go-webpagetest
package webpagetest

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// TestInfo is info about test
type TestInfo struct {
	URL           string `json:"url"`
	Runs          int    `json:"runs"`
	FirstViewOnly int    `json:"fvonly"`
	Web10         int    `json:"web10"`     // Stop Test at Document Complete
	IgnoreSSL     int    `json:"ignoreSSL"` // Ignore SSL Certificate Errors
	Video         string `json:"video"`
	Label         string `json:"label"`
	Priority      int    `json:"priority"`
	Location      string `json:"location"`
	Browser       string `json:"browser"`

	Connectivity string `json:"connectivity"`
	BandwidthIn  int    `json:"bwIn"`
	BandwidthOut int    `json:"bwOut"`
	Latency      int    `json:"latency"`

	// It can be string or int
	RawPacketLossRate *json.RawMessage `json:"plr"`
	PacketLossRate    int

	Tcpdump      int `json:"tcpdump"`  // Capture network packet trace (tcpdump)
	Timeline     int `json:"timeline"` // Capture Dev Tools Timeline
	Trace        int `json:"trace"`    // Capture Chrome Trace (about://tracing)
	Bodies       int `json:"bodies"`
	NetLog       int `json:"netlog"`    // Capture Network Log
	Standards    int `json:"standards"` // Disable Compatibility View (IE Only)
	NoScript     int `json:"noscript"`  // Disable Javascript
	Pngss        int `json:"pngss"`
	ImageQuality int `json:"iq"`
	KeepUA       int `json:"keepua"` // Preserve original User Agent string
	Mobile       int `json:"mobile"`
	Scripted     int `json:"scripted"`
}

// TestStatus is status of a test
type TestStatus struct {
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`

	ID           string `json:"id"`
	TestID       string `json:"testId"`
	Location     string `json:"location"`
	StartTime    string `json:"startTime"`
	CompleteTime string `json:"completeTime"`

	Runs        int `json:"runs"`
	BehindCount int `json:"behindCount"`

	Remote         bool `json:"remote"` // Relay Test
	FirstViewOnly  int  `json:"fvonly"`
	Elapsed        int  `json:"elapsed"`
	ElapsedUpdate  int  `json:"elapsedUpdate"`
	TestsExpected  int  `json:"testsExpected"`
	TestsCompleted int  `json:"testsCompleted"`

	FirstViewRunsCompleted  int `json:"fvRunsCompleted"`
	RepeatViewRunsCompleted int `json:"rvRunsCompleted"`

	TestInfo TestInfo `json:"testInfo"`
}

type jsonTestStatus struct {
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`

	Data TestStatus `json:"data"`
}

// GetTestStatus will return status of test run by given testID
// StatusCode 200 indicates test is completed. 1XX means the test is still
// in progress. And 4XX indicates some error.
func (c *Client) GetTestStatus(testID string) (*TestStatus, error) {
	body, err := c.query("/testStatus.php", url.Values{"test": []string{testID}})
	if err != nil {
		return nil, err
	}

	var result jsonTestStatus
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.StatusCode > 200 {
		return nil, fmt.Errorf("%s", result.StatusText)
	}

	return &result.Data, nil
}
