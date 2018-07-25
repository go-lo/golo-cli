package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

var (
	client oauthClient
)

type Job struct {
	Name     string `json:"name",yaml:"name"`
	Users    int    `json:"users",yaml:"users"`
	Duration int64  `json:"duration",yaml:"duration",`
	Binary   string `json:"binary"`
}

type queueResponse struct {
	Queued bool `json:"queued"`
}

type uploadResponse struct {
	Binary string `json:"binary"`
}

func QueueJob(hb HostBinary, job Job) (err error) {
	job.Binary = hb.Binary

	b, err := json.Marshal(job)
	if err != nil {
		return
	}

	body := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", setURL(hb.Host, "/queue"), body)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Received %q from %s, %s", resp.Status, hb.Host, string(respBody))
	}

	uploadResp := new(queueResponse)
	err = json.Unmarshal(respBody, uploadResp)
	if err != nil {
		return
	}

	if !uploadResp.Queued {
		err = fmt.Errorf("Could not queue job, received: %+v", respBody)
	}

	return
}

func UploadSchedule(f, addr string) (hb HostBinary, err error) {
	hb.Host = addr

	r, err := os.Open(f)
	if err != nil {
		return
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", r.Name())
	if err != nil {
		return
	}

	if _, err = io.Copy(fw, r); err != nil {
		return
	}

	w.Close()

	req, err := http.NewRequest("POST", setURL(addr, "/upload"), &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	uploadResp := new(uploadResponse)
	err = json.Unmarshal(body, uploadResp)
	hb.Binary = uploadResp.Binary

	return
}

func setURL(addr, path string) string {
	u := new(url.URL)
	u.Host = fmt.Sprintf("%s:8081", addr)
	u.Path = path
	u.Scheme = "http"

	return u.String()
}
