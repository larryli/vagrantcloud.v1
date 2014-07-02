package vagrantcloud

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// Boxes are the primary resource on Vagrant Cloud.
// Before creating versions with attached providers, you'll need to create a box.
type Box struct {
	api                 *Api
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Tag                 string    `json:"tag"`
	Name                string    `json:"name"`
	ShortDescription    string    `json:"short_description"`
	DescriptionHtml     string    `json:"description_html"`
	DescriptionMarkdown string    `json:"description_markdown"`
	Username            string    `json:"username"`
	Private             bool      `json:"private"`
	CurrentVersion      Version   `json:"current_version"`
	Versions            []Version `json:"versions"`
}

func (a *Api) Box(username, name string) *Box {
	b := &Box{
		Username: username,
		Name:     name,
	}
	b.init(a)
	return b
}

func (b *Box) init(a *Api) {
	b.api = a
	b.CurrentVersion.init(b)
	for n := range b.Versions {
		(&b.Versions[n]).init(b)
	}
}

func (b *Box) parseBody(body []byte) error {
	err := json.Unmarshal(body, b)
	if err != nil {
		return err
	}
	b.init(b.api)
	return nil
}

func (b *Box) Uri() string {
	return "/box/" + b.Username + "/" + b.Name
}

// RETRIEVE A BOX
func (b *Box) Get() error {
	body, err := b.api.Get(b.Uri())
	if err != nil {
		return err
	}
	return b.parseBody(body)
}

// CREATE A BOX
//
// 	Name (required)
// 		The name of the box, used to identify it.
//		The name makes up the latter half of the tag.
//		It has a maximum length of 36 characters and must contain only letters,
//		numbers, dashes, underscores or periods.
// 	Username
//		The username to assign the box to.
//		You must be a member of the organization and have the ability to create boxes.
//		Defaults to the users username that is making the API request.
// 	ShortDescription
//		The short description is used on small box previews,
//		in search results and other places where displaying markdown isn't functional.
//		It has a maxium length of 120 characters.
// 	DescriptionMarkdown
//		Markdown text used as a full-length and in-depth description of the box.
//		Markdown is parsed according to GitHub Flavored Markdown.
//		There is no maximum length.
// 	Private
//		A boolean if the box should be private or not.
func (b *Box) New() error {
	params := url.Values{}
	params.Add("box[name]", b.Name)
	if b.Username != "" {
		params.Add("box[username]", b.Username)
	}
	if b.ShortDescription != "" {
		params.Add("box[short_description]", b.ShortDescription)
	}
	if b.DescriptionMarkdown != "" {
		params.Add("box[description]", b.DescriptionMarkdown)
	}
	if b.Private {
		params.Add("box[is_private]", strconv.FormatBool(b.Private))
	}
	body, err := b.api.Post("/boxes", params)
	if err != nil {
		return err
	}
	return b.parseBody(body)
}

// UPDATE A BOX
//
// 	Name (required)
// 		The name of the box, used to identify it.
//		The name makes up the latter half of the tag.
//		It has a maximum length of 36 characters and must contain only letters,
//		numbers, dashes, underscores or periods.
// 	Username (required)
//		The username to assign the box to.
//		You must be a member of the organization and have the ability to create boxes.
//		Defaults to the users username that is making the API request.
// 	ShortDescription
//		The short description is used on small box previews,
//		in search results and other places where displaying markdown isn't functional.
//		It has a maxium length of 120 characters.
// 	DescriptionMarkdown
//		Markdown text used as a full-length and in-depth description of the box.
//		Markdown is parsed according to GitHub Flavored Markdown.
//		There is no maximum length.
// 	Private
//		A boolean if the box should be private or not.
func (b *Box) Set() error {
	params := url.Values{}
	params.Add("box[short_description]", b.ShortDescription)
	params.Add("box[description]", b.DescriptionMarkdown)
	params.Add("box[is_private]", strconv.FormatBool(b.Private))
	body, err := b.api.Put(b.Uri(), params)
	if err != nil {
		return err
	}
	return b.parseBody(body)
}

// DESTROY A BOX
//
// 	Name (required)
// 		The name of the box, used to identify it.
//		The name makes up the latter half of the tag.
//		It has a maximum length of 36 characters and must contain only letters,
//		numbers, dashes, underscores or periods.
// 	Username (required)
//		The username to assign the box to.
//		You must be a member of the organization and have the ability to create boxes.
//		Defaults to the users username that is making the API request.
func (b *Box) Delete() error {
	body, err := b.api.Delete(b.Uri())
	if err != nil {
		return err
	}
	return b.parseBody(body)
}
