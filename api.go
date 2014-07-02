package vagrantcloud

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseUrl = "https://vagrantcloud.com"
	apiUri  = "/api/v1"
)

type Api struct {
	token string
}

// All requests must be authenticated with an access_token and sent as a URL parameter.
// This token can be generated or revoked on the account tokens page.
// Your token will have access to all resources your account has access to.
func New(token string) *Api {
	a := &Api{
		token: token,
	}
	return a
}

func NewFromFile(fname string) (*Api, error) {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		return nil, err
	}
	return New(string(token)), nil
}

func (a *Api) buildUrl(url string) string {
	return baseUrl + apiUri + url
}

func (a *Api) contentType(header http.Header) {
	header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func (a *Api) response(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return body, NewError(resp.Status, string(body))
	}
	return body, nil
}

func (a *Api) Get(uri string) ([]byte, error) {
	u, err := url.ParseRequestURI(a.buildUrl(uri))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", a.token)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	return a.response(resp)
}

func (a *Api) Download(uri string, data io.Writer) error {
	u, err := url.ParseRequestURI(baseUrl + uri)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("access_token", a.token)
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(data, resp.Body)
	return err
}

func (a *Api) Post(uri string, params url.Values) ([]byte, error) {
	u, err := url.ParseRequestURI(a.buildUrl(uri))
	if err != nil {
		return nil, err
	}
	params.Set("access_token", a.token)
	resp, err := http.PostForm(u.String(), params)
	if err != nil {
		return nil, err
	}
	return a.response(resp)
}

func (a *Api) Put(uri string, params url.Values) ([]byte, error) {
	u, err := url.ParseRequestURI(a.buildUrl(uri))
	if err != nil {
		return nil, err
	}
	params.Set("access_token", a.token)
	req, err := http.NewRequest("PUT", u.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	a.contentType(req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return a.response(resp)
}

func (a *Api) Upload(uri string, data io.Reader) ([]byte, error) {
	u, err := url.ParseRequestURI(a.buildUrl(uri))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("access_token", a.token)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("PUT", u.String(), data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return a.response(resp)
}

func (a *Api) Delete(uri string) ([]byte, error) {
	u, err := url.ParseRequestURI(a.buildUrl(uri))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("access_token", a.token)
	req, err := http.NewRequest("DELETE", u.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	a.contentType(req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return a.response(resp)
}
