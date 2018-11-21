package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	GERRIT_URL = "https://review.openstack.org/"
)

var (
	lastChange GerritChange
)

func createQueryUrl() string {
	return fmt.Sprintf("%s/changes/?q=status:open+project:openstack/stx-tools&n=1",
		GERRIT_URL)
}

type GerritChange struct {
	Id          string            `json:"id"`
	Project     string            `json:"project"`
	Branch      string            `json:"branch"`
	Topic       string            `json:"topic"`
	Hashtags    []string          `json:"hashtags"`
	ChangeId    string            `json:"change_id"`
	Subject     string            `json:"subject"`
	Status      string            `json:"status"`
	Created     string            `json:"created"`
	Updated     string            `json:"updated"`
	SubmitType  string            `json:"submit_type"`
	Mergeable   bool              `json:"mergeable"`
	Submittable bool              `json:"submittable"`
	Insertions  int               `json:"insertions"`
	Deletions   int               `json:"deletions"`
	Number      int               `json:"_number"`
	Owner       GerritChangeOwner `json:"owner"`
	MoreChanges bool              `json:"_more_changes"`
}

type GerritChangeOwner struct {
	Id int `json:"_account_id"`
}

// sendNotification creates and run the notify-send command to launch the
// notification.
func sendNotification(g GerritChange) (err error) {
	msg := fmt.Sprintf("Title: %s\nUrl: %s#/c/%d/",
		g.Subject, GERRIT_URL, g.Number)
	cmd := exec.Command("notify-send",
		"--icon=dialog-information",
		"New gerrit change",
		msg)

	_, err = cmd.Output()
	return err
}

func queryAndSend() {
	url := createQueryUrl()
	log.Info("Querying data")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// This is really ugly with the Gerrit REST API always returns a )]}'
	// at the beginning of the data.

	var r []GerritChange
	err = json.Unmarshal(responseData[4:], &r)
	if err != nil {
		log.Error(err)
	}

	if len(r) > 0 {
		if lastChange.Id != r[0].Id {
			lastChange = r[0]
			sendNotification(lastChange)
		}
	} else {
		log.Fatalf("Response had invalid length: %d", len(r))
	}

}

func main() {
	timer := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-timer.C:
			queryAndSend()
		}
	}
}
