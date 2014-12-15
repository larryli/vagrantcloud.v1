package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/larryli/vagrantcloud.v1"
	"io/ioutil"
	"log"
	"strings"
)

type Arch struct {
	name, info, box string
}

type Release string

type Box vagrantcloud.Box

type Version vagrantcloud.Version

type CodeNames map[string]string

const (
	url      = "https://cloud-images.ubuntu.com/vagrant/"
	notFound = "404 Not Found"
	see      = "\n\nSee https://github.com/larryli/vagrantcloud.v1/tree/master/update-ubuntu-vagrant-box"
)

var (
	api      *vagrantcloud.Api
	names    = CodeNames{}
	username = flag.String("username", "larryli", "username")
	token    = flag.String("token", "", "access_token")
	test     = flag.Bool("test", false, "test, no effect")
	codename = flag.String("codename", "", "ubuntu code name file(json)")
	arches   = []Arch{
		{
			name: "64",
			info: "amd64",
			box:  "-server-cloudimg-amd64-vagrant-disk1.box",
		},
		// {
		// 	name: "64-juju",
		// 	info: "amd64 with Juju",
		// 	box:  "-server-cloudimg-amd64-juju-vagrant-disk1.box",
		// },
		{
			name: "32",
			info: "i386",
			box:  "-server-cloudimg-i386-vagrant-disk1.box",
		},
		// {
		// 	name: "32-juju",
		// 	info: "i386 with Juju",
		// 	box:  "-server-cloudimg-i386-juju-vagrant-disk1.box",
		// },
	}
)

func fatal(err error, a ...interface{}) {
	if err != nil {
		a = append(a, err)
		log.Fatalln(a...)
	}
}

func isChildren(s string) bool {
	return !strings.HasPrefix(s, "/") && !strings.HasPrefix(s, "?") && strings.HasSuffix(s, "/")
}

func initCodeNames() {
	if *codename != "" {
		text, err := ioutil.ReadFile(*codename)
		fatal(err, "read "+*codename)
		fatal(json.Unmarshal(text, &names), "unmarshal "+*codename)
	}
}

func fetchReleases() (releases []Release) {
	doc, err := goquery.NewDocument(url)
	fatal(err, "fetch "+url)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Attr("href"); ok {
			if isChildren(link) {
				releases = append(releases, Release(strings.Trim(link, "/")))
			}
		}
	})
	return
}

func (c CodeNames) get(name string) string {
	if str, ok := c[name]; ok {
		return str
	}
	return name
}

func (r Release) name() string {
	return string(r)
}

func (r Release) url() string {
	return url + r.name() + "/"
}

func (r Release) title(info, version string) string {
	ret := "Official Ubuntu Server " + strings.Title(names.get(r.name()))
	if info != "" {
		ret += " " + info
	}
	ret += " builds"
	if version != "" {
		ret += " (latest " + version + ")"
	}
	return ret
}

func (r Release) fetch() (versions []string) {
	doc, err := goquery.NewDocument(r.url())
	fatal(err, "fetch "+r.url())
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
			box.ShortDescription = r.title(t.info, "")
			box.DescriptionMarkdown = r.url() + see
			todo = fmt.Sprintf("add \"%s\"", box.Uri())
			if !(*test) {
				fatal(box.New(), todo)
			}
			log.Println(todo)
		} else {
			fatal(err, todo)
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
				box.ShortDescription = r.title(t.info, version)
				todo = fmt.Sprintf("update \"%s\": \"%s\"", box.Uri(), box.ShortDescription)
				if !(*test) {
					fatal(box.Set(), todo)
				}
				log.Println(todo)
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
		fatal(version.Delete(), todo)
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
	v := box.Version("0.0.1")
	v.Version = version
	v.DescriptionMarkdown = image + see
	todo := fmt.Sprintf("add \"%s\" Version: \"%s\"", v.Uri(), v.Version)
	if !(*test) {
		fatal(v.New(), todo)
	}
	p := v.Provider(vagrantcloud.ProviderVirtualbox)
	p.OriginalUrl = image
	todo = fmt.Sprintf("add \"%s\" Version: \"%s\"", p.Uri(), v.Version)
	if !(*test) {
		fatal(p.New(), todo)
	}
	todo = fmt.Sprintf("public \"%s\" Version: \"%s\" Url: \"%s\"", v.Uri(), v.Version, p.OriginalUrl)
	if !(*test) {
		fatal(v.Release(), todo)
	}
	log.Println(todo)
}

func main() {
	flag.Parse()
	if !(*test) && *token == "" {
		flag.Usage()
	} else {
		initCodeNames()
		api = vagrantcloud.New(*token)
		log.Println("start")
		for _, release := range fetchReleases() {
			release.scan()
		}
		log.Println("end")
	}
}
