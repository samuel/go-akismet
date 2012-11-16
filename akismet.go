package akismet

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const userAgent = "go-akismet/0.1"

const (
	AkismetDomain = "rest.akismet.com"
	TypePadDomain = "api.antispam.typepad.com"
)

type ErrUnexpectedHttpResponse struct {
	Response *http.Response
}

func (e *ErrUnexpectedHttpResponse) Error() string {
	return fmt.Sprintf("{Err unexpected http response '%s'}", e.Response)
}

type ErrUnexpectedResponse string

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("{Err unexpected response '%s'}", e)
}

type Akismet struct {
	key       string
	blog      string // The front page or home URL of the instance making the request. For a blog or wiki this would be the front page. Note: Must be a full URI, including http://.
	apiDomain string
}

func New(apiKey string, blog string, apiDomain string) (*Akismet, error) {
	if apiKey == "" {
		return nil, errors.New("akismet: apiKey must not be blank")
	}
	if blog == "" {
		return nil, errors.New("akismet: blog must not be blank")
	}
	return &Akismet{
		key:       apiKey,
		blog:      blog,
		apiDomain: apiDomain,
	}, nil
}

func (a *Akismet) post(cmd string, vals url.Values) (*http.Response, error) {
	var url string
	if cmd == "verify-key" {
		url = "http://" + a.apiDomain + "/1.1/" + cmd
	} else {
		url = "http://" + a.key + "." + a.apiDomain + "/1.1/" + cmd
	}
	body := bytes.NewReader([]byte(vals.Encode()))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(body.Len())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err == nil && res.StatusCode != 200 {
		err = &ErrUnexpectedHttpResponse{res}
	}
	return res, err
}

// [X-Akismet-Server:[192.168.7.149] X-Akismet-Debug-Help:[We were unable to parse your blog URI]
// X-Spam-Requestid:[L/duyRkst5Yha89r6qYnTKRO3/IgP+9tb0hq0ZUYD7Q=]]
func (a *Akismet) VerifyKey() (bool, error) {
	res, err := a.post("verify-key", url.Values{"key": {a.key}, "blog": {a.blog}})
	if err != nil {
		return false, err
	}
	b := make([]byte, 16)
	n, err := res.Body.Read(b)
	if err != nil {
		return false, err
	}
	return string(b[:n]) == "valid", nil
}

func (a *Akismet) CommentCheck(comment Comment) (bool, error) {
	vals := url.Values{}
	comment.ToValues(vals)

	res, err := a.post("comment-check", vals)

	b := make([]byte, 16)
	n, err := res.Body.Read(b)
	if err != nil {
		return false, err
	}
	resText := string(b[:n])
	switch resText {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, ErrUnexpectedResponse(resText)
}

// submit-spam / submit-ham -> "Thanks for making the web a better place."
