// adapted from https://github.com/olegfedoseev/go-webpagetest
package webpagetest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is client of WebPageTest
type Client struct {
	Host string
}

// NewClient returns new ready to use Client
func NewClient(host string) (*Client, error) {
	validURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &Client{
		Host: validURL.String(),
	}, nil
}

func (c *Client) query(api string, params url.Values) ([]byte, error) {
	queryURL := c.Host + api + "?" + params.Encode()
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to GET \"%s\": %v", queryURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status is no OK: %v [%v]", resp.StatusCode, string(body))
	}

	return body, nil
}

// RunTest will submit given test to WPT server
func (c *Client) RunTest(settings TestSettings) (string, error) {
	resp, err := http.PostForm(c.Host+"/runtest.php", settings.GetFormParams())
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		StatusCode int    `json:"statusCode"`
		StatusText string `json:"statusText"`
		Data       struct {
			TestID  string `json:"testId"`
			UserURL string `json:"userUrl"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.StatusCode > 200 {
		return "", fmt.Errorf("StatusCode > 200: %v: %v", result.StatusCode, result.StatusText)
	}

	// fmt.Printf("Result URL for %v: %v\n", settings.URL, result.Data.UserURL)
	return result.Data.TestID, nil
}

// StatusCallback is helper type for function to be called while waiting for test to complete
type StatusCallback func(testID string, status *TestStatus)

// RunTestAndWait will start new WebPageTest test run with given TestSettings and will wait for it
// to complete. While it wait, it will poll status updates from server and will call StatusCallback with it
func (c *Client) RunTestAndWait(settings TestSettings, callback StatusCallback) (*ResultData, error) {
	testID, err := c.RunTest(settings)
	if err != nil {
		return nil, err
	}

	for {
		result, err := c.GetTestStatus(testID)
		if err != nil {
			return nil, err
		}
		// Call callback
		if callback != nil {
			go callback(testID, result)
		}
		if result.StatusCode < 200 {
			time.Sleep(10 * time.Second)
		}
		if result.StatusCode >= 200 {
			break
		}
	}

	testResult, err := c.GetTestResult(testID)
	if err != nil {
		return nil, err
	}
	return testResult, nil
}
