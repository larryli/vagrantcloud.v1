package vagrantcloud

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type VersionStatus string

const (
	VersionUnreleased VersionStatus = "unreleased"
	VersionActive     VersionStatus = "active"
	VersionRevoked    VersionStatus = "revoked"
)

// Versions represent new releases for boxes,
// and contain both information about the changes in that version as well as being the parent object for providers,
// which contain the eventual box file.
type Version struct {
	api                 *Api
	box                 *Box
	Version             string        `json:"version"`
	Status              VersionStatus `json:"status"`
	DescriptionHtml     string        `json:"description_html"`
	DescriptionMarkdown string        `json:"description_markdown"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	Number              int           `json:"number"`
	Downloads           int           `json:"downloads"`
	ReleaseUrl          string        `json:"release_url"`
	RevokeUrl           string        `json:"revoke_url"`
	Providers           []Provider    `json:"providers"`
}

func (b *Box) Version(number int) *Version {
	v := &Version{
		Number: number,
	}
	v.init(b)
	return v
}

func (v *Version) init(b *Box) {
	v.api = b.api
	v.box = b
	for n := range v.Providers {
		(&v.Providers[n]).init(v)
	}
}

func (v *Version) parseBody(body []byte) error {
	err := json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	v.init(v.box)
	return nil
}

func (v *Version) Uri() string {
	return v.box.Uri() + "/version/" + strconv.Itoa(v.Number)
}

// RETRIEVE A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
func (v *Version) Get() error {
	body, err := v.api.Get(v.Uri())
	if err != nil {
		return err
	}
	return v.parseBody(body)
}

// CREATE A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
//	DescriptionMarkdown
//		Markdown text used as a full-length and in-depth description of the version,
//		typically for denoting changes introduced. Markdown is parsed according to GitHub Flavored Markdown.
//		There is no maximum length.
//
// When a version is first created, its status is set to unreleased.
func (v *Version) New() error {
	params := url.Values{}
	params.Add("version[version]", v.Version)
	if v.DescriptionMarkdown != "" {
		params.Add("version[description]", v.DescriptionMarkdown)
	}
	body, err := v.api.Post(v.box.Uri()+"/versions", params)
	if err != nil {
		return err
	}
	return v.parseBody(body)
}

// UPDATE A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
//	DescriptionMarkdown
//		Markdown text used as a full-length and in-depth description of the version,
//		typically for denoting changes introduced. Markdown is parsed according to GitHub Flavored Markdown.
//		There is no maximum length.
//
// You cannot modify the status attribute directly,
// so their are seperate endpoints to revoke and release versions.
func (v *Version) Set() error {
	params := url.Values{}
	params.Add("version[description]", v.DescriptionMarkdown)
	body, err := v.api.Put(v.Uri(), params)
	if err != nil {
		return err
	}
	return v.parseBody(body)
}

// DESTROY A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
func (v *Version) Delete() error {
	body, err := v.api.Delete(v.Uri())
	if err != nil {
		return err
	}
	return v.parseBody(body)
}

// RELEASE A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
//
// This allows you to update the providers and description prior to release.
// Once a version is ready to release, you can then make a request to move it to an active state.
func (v *Version) Release() error {
	params := url.Values{}
	body, err := v.api.Put(v.Uri()+"/release", params)
	if err != nil {
		return err
	}
	return v.parseBody(body)
}

// REVOKE A VERSION
//
// 	Version (required)
//		The version number, typically incrementing a previous version.
//		We validate this version string based on Semantic Versioning.
//		We only require that the string matches a pattern that could be semver,
//		and don't validate that the version comes after your previous versions, and so on.
//
// Versions that have been "released" can then no longer be deleted.
// Instead, you should "revoke" a released version, setting the status as revoked.
// This stops access to the version from Vagrant, but maintains the history of the version.
func (v *Version) Revoke() error {
	params := url.Values{}
	body, err := v.api.Put(v.Uri()+"/revoke", params)
	if err != nil {
		return err
	}
	return v.parseBody(body)
}
