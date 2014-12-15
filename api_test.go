package vagrantcloud_test

import (
	"github.com/larryli/vagrantcloud.v1"
	"testing"
)

func TestApi(t *testing.T) {
	url := "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-i386-vagrant-disk1.box"
	a, err := vagrantcloud.NewFromFile("token.txt")
	if err != nil {
		t.Fatal(err)
	}
	b := a.Box("", "test")
	if err := b.New(); err != nil {
		t.Fatal(err)
	}
	b.DescriptionMarkdown = url
	if err := b.Set(); err != nil {
		t.Fatal(err)
	}
	if err := b.Get(); err != nil {
		t.Fatal(err)
	}

	v := b.Version("0.0.1")
	v.Version = "0.0.1"
	if err := v.New(); err != nil {
		t.Fatal(err)
	}
	v.DescriptionMarkdown = url
	if err := v.Set(); err != nil {
		t.Fatal(err)
	}
	if err := v.Get(); err != nil {
		t.Fatal(err)
	}

	p := v.Provider(vagrantcloud.ProviderVirtualbox)
	p.OriginalUrl = url
	if err := p.New(); err != nil {
		t.Fatal(err)
	}
	p.OriginalUrl = url
	if err := p.Set(); err != nil {
		t.Fatal(err)
	}
	if err := p.Get(); err != nil {
		t.Fatal(err)
	}

	if err := v.Release(); err != nil {
		t.Fatal(err)
	}
	if err := v.Revoke(); err != nil {
		t.Fatal(err)
	}

	if err := p.Delete(); err != nil {
		t.Fatal(err)
	}
	if err := v.Delete(); err != nil {
		t.Fatal(err)
	}
	if err := b.Delete(); err != nil {
		t.Fatal(err)
	}
}
