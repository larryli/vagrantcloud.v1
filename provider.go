package vagrantcloud

import (
	"encoding/json"
	"io"
	"net/url"
	"time"
)

type ProviderName string

const (
	ProviderVirtualbox    ProviderName = "virtualbox"
	ProviderVmwareDesktop ProviderName = "vmware_desktop"
	ProviderDigitalocean  ProviderName = "digitalocean"
	ProviderAws           ProviderName = "aws"
	ProviderRackspace     ProviderName = "rackspace"
	ProviderHyperv        ProviderName = "hyperv"
)

// Providers contain the pointers to the box files,
// be it a hosted or self-hosted box.
// Versions can have many providers,
// each which represents a Vagrant compatible provider,
// either from Vagrant Core as a 3rd party plugin.
type Provider struct {
	api         *Api
	box         *Box
	version     *Version
	Name        ProviderName `json:"name"`
	Hosted      bool         `json:"hosted"`
	HostedToken string       `json:"hosted_token"`
	OriginalUrl string       `json:"original_url"`
	UploadUrl   string       `json:"upload_url"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DownloadUrl string       `json:"download_url"`
}

func (v *Version) Provider(name ProviderName) *Provider {
	p := &Provider{
		Name: name,
	}
	p.init(v)
	return p
}

func (p *Provider) init(v *Version) {
	p.api = v.api
	p.box = v.box
	p.version = v
}

func (p *Provider) parseBody(body []byte) error {
	err := json.Unmarshal(body, p)
	if err != nil {
		return err
	}
	p.init(p.version)
	return nil
}

func (p *Provider) Uri() string {
	return p.version.Uri() + "/provider/" + string(p.Name)
}

// RETRIEVE A PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
func (p *Provider) Get() error {
	body, err := p.api.Get(p.Uri())
	if err != nil {
		return err
	}
	return p.parseBody(body)
}

// CREATE A PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
//	OriginalUrl
//		An HTTP URL to the box file.
//		This must be accessible at this URL from the machine where you expect a user to download the box by using Vagrant.
//		If ommitted, we assume you wish to host the provider with Vagrant Cloud.
//
// The provider API is used to host boxes.
// To create a hosted box, simply omit the URL parameter.
// You will then be able to use the upload endpoint to upload a box to us.
func (p *Provider) New() error {
	params := url.Values{}
	params.Add("provider[name]", string(p.Name))
	if p.OriginalUrl != "" {
		params.Add("provider[url]", p.OriginalUrl)
	}
	body, err := p.api.Post(p.version.Uri()+"/providers", params)
	if err != nil {
		return err
	}
	return p.parseBody(body)
}

// UPDATE A PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
//	OriginalUrl
//		An HTTP URL to the box file.
//		This must be accessible at this URL from the machine where you expect a user to download the box by using Vagrant.
//		If ommitted, we assume you wish to host the provider with Vagrant Cloud.
func (p *Provider) Set() error {
	params := url.Values{}
	params.Add("provider[url]", p.OriginalUrl)
	body, err := p.api.Put(p.Uri(), params)
	if err != nil {
		return err
	}
	return p.parseBody(body)
}

// DESTROY A PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
func (p *Provider) Delete() error {
	body, err := p.api.Delete(p.Uri())
	if err != nil {
		return err
	}
	return p.parseBody(body)
}

// UPLOAD A .BOX FOR PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
//  data (io.Reader, required)
//
// The upload path returns a URL that you can then PUT the boxes payload to.
// After streaming the box,
// you can confirm the upload by comparing the token provided in the UPLOAD response with the token returned from the providers GET route.
// When these tokens match, the upload has been successful.
func (p *Provider) Upload(data io.Reader) error {
	body, err := p.api.Upload(p.Uri()+"/upload", data)
	if err != nil {
		return err
	}
	return p.parseBody(body)
}

// DOWNLOAD A .BOX FOR PROVIDER
//
//	Name (required)
//		The name of the provider. Vagrant will use this to determine compatible boxes on the client.
//		Common providers include virtualbox, vmware_desktop, digitalocean, aws, rackspace, and hyperv.
func (p *Provider) Download(data io.Writer) error {
	err := p.api.Download("/"+p.box.Username+"/"+p.box.Name+"/version/"+p.version.Number+"/provider/"+string(p.Name)+".box", data)
	if err != nil {
		return err
	}
	return nil
}
