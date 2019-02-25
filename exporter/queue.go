package exporter

import (
	"container/list"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/steffenmllr/prometheus-webpagetest-exporter/webpagetest"
)

type TestRun struct {
	Id        string
	Url       string
	RemoteUrl string
	Status    *webpagetest.TestStatus
	Result    webpagetest.ResultData
	CreatedAt time.Time
}

func (lq *ListQueue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./exporter/index.html"))
	var testRuns []TestRun
	for e := lq.container.Front(); e != nil; e = e.Next() {
		testRuns = append(testRuns, e.Value.(TestRun))
	}
	tmpl.Execute(w, testRuns)
}

func (lq *ListQueue) UpdateStatus(testID string, status *webpagetest.TestStatus) {
	found := false
	for e := lq.container.Front(); e != nil; e = e.Next() {
		if e.Value.(TestRun).Id == testID {
			e.Value = TestRun{Id: testID, Status: status, CreatedAt: e.Value.(TestRun).CreatedAt}
			found = true
		}
	}

	if !found {
		lq.Add(TestRun{Id: testID, Status: status, CreatedAt: time.Now()})
	}
}

func (lq *ListQueue) AddTestResult(Result *webpagetest.ResultData, host string) {
	for e := lq.container.Front(); e != nil; e = e.Next() {
		testRun := e.Value.(TestRun)
		if testRun.Id == Result.Id {
			RemoteUrl := fmt.Sprintf("%v/results.php?test=%v", host, testRun.Id)
			e.Value = TestRun{Id: testRun.Id, Status: testRun.Status, Result: *Result, RemoteUrl: RemoteUrl, CreatedAt: testRun.CreatedAt}
		}
	}
}

type ListQueue struct {
	container *list.List
	maxSize   int
}

func NewListQueue() *ListQueue {
	return &ListQueue{container: list.New(), maxSize: 50}
}

func (lq *ListQueue) Add(item TestRun) {
	lq.container.PushFront(item)
	oversized := lq.container.Len() - lq.maxSize
	if oversized > 0 {
		for i := 0; i <= oversized; i++ {
			el := lq.container.Back()
			if el != nil {
				lq.container.Remove(el)
			}
		}
	}
}

func (lq *ListQueue) Empty() bool {
	return lq.container.Len() == 0
}

func (lq *ListQueue) Clear() {
	if !lq.Empty() {
		lq.container = list.New()
	}
}
