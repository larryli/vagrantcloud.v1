package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/larryli/vagrantcloud.v1"
	"log"
	"strings"
)

type Arch struct {
	name, info, box string
}

type Release string

type Box vagrantcloud.Box

type Version vagrantcloud.Version

const (
	url      = "https://cloud-images.ubuntu.com/vagrant/"
	notFound = "404 Not Found"
	see      = "\n\nSee https://github.com/larryli/vagrantcloud.v1/tree/master/update-ubuntu-vagrant-box"
)

var (
	api      *vagrantcloud.Api
	username = flag.String("username", "larryli", "username")
	token    = flag.String("token", "", "access_token")
	test     = flag.Bool("test", false, "test, no effect")
	arches   = []Arch{
		{
			name: "64",
			info: "amd64",
			box:  "-server-cloudimg-amd64-vagrant-disk1.box",
		},
		{
			name: "64juju",
			info: "amd64 with juju",
			box:  "-server-cloudimg-amd64-juju-vagrant-disk1.box",
		},
		{
			name: "32",
			info: "i386",
			box:  "-server-cloudimg-i386-vagrant-disk1.box",
		},
		{
			name: "32juju",
			info: "i386 with juju",
			box:  "-server-cloudimg-i386-juju-vagrant-disk1.box",
		},
	}
)

func Fatal(err error, a ...interface{}) {
	if err != nil {
		a = append(a, err)
		log.Fatalln(a...)
	}
}

func isChildren(s string) bool {
	return !strings.HasPrefix(s, "/") && !strings.HasPrefix(s, "?")
}

func fetchReleases() (releases []Release) {
	doc, err := goquery.NewDocument(url)
	Fatal(err, "fetch "+url)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Attr("href"); ok {
			if isChildren(link) {
				releases = append(releases, Release(strings.Trim(link, "/")))
			}
		}
	})
	return
}

func (r Release) name() string {
	return string(r)
}

func (r Release) url() string {
	return url + r.name() + "/"
}

func (r Release) fetch() (versions []string) {
	doc, err := goquery.NewDocument(r.url())
	Fatal(err, "fetch "+r.url())
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Attr("href"); ok {
			if isChildren(link) && link != "current/" {
				versions = append(versions, strings.Trim(link, "/"))
			}
		}
	})
	return
}

func (r Release) scan() {
	versions := r.fetch()
	for _, t := range arches {
		box := api.Box(*username, r.name()+t.name)
		todo := fmt.Sprintf("fetch \"%s\"", box.Uri())
		err := box.Get()
		if err != nil && strings.EqualFold(err.Error(), notFound) {
			box.ShortDescription = "Ubuntu " + strings.Title(r.name()) + " " + t.info
			box.DescriptionMarkdown = r.url() + see
			todo = fmt.Sprintf("add \"%s\"", box.Uri())
			if !(*test) {
				Fatal(box.New(), todo)
			}
			log.Println(todo)
		} else {
			Fatal(err, todo)
		}
		for _, version := range box.Versions {
			v := Version(version)
			if !v.exists(versions) {
				v.delete()
			}
		}
		for _, version := range versions {
			b := (*Box)(box)
			if !b.find(version) {
				image := r.url() + version + "/" + r.name() + t.box
				b.add(version, image)
			}
		}
	}
}

func (v *Version) exists(versions []string) (found bool) {
	for _, test := range versions {
		if test == v.Version {
			found = true
			return
		}
	}
	return
}

func (v *Version) delete() {
	version := (*vagrantcloud.Version)(v)
	todo := fmt.Sprintf("delete \"%s\" Version: \"%s\"", version.Uri(), version.Version)
	if !(*test) {
		Fatal(version.Delete(), todo)
	}
	log.Println(todo)
}

func (b *Box) find(version string) (found bool) {
	for _, test := range b.Versions {
		if test.Version == version {
			found = true
			return
		}
	}
	return
}

func (b *Box) add(version, image string) {
	box := (*vagrantcloud.Box)(b)
	v := box.Version(0)
	v.Version = version
	v.DescriptionMarkdown = image + see
	todo := fmt.Sprintf("add \"%s\" Version: \"%s\"", v.Uri(), v.Version)
	if !(*test) {
		Fatal(v.New(), todo)
	}
	p := v.Provider(vagrantcloud.ProviderVirtualbox)
	p.OriginalUrl = image
	todo = fmt.Sprintf("add \"%s\" Version: \"%s\"", p.Uri())
	if !(*test) {
		Fatal(p.New(), todo)
	}
	todo = fmt.Sprintf("public \"%s\" Version: \"%s\" Url: \"%s\"", v.Uri(), v.Version, p.OriginalUrl)
	if !(*test) {
		Fatal(v.Release(), todo)
	}
	log.Println(todo)
}

func main() {
	flag.Parse()
	api = vagrantcloud.New(*token)
	log.Println("start")
	for _, release := range fetchReleases() {
		release.scan()
	}
	log.Println("end")
}
