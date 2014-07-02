Vagrant Cloud API
=================

[![GoDoc](https://godoc.org/github.com/larryli/vagrantcloud.v1?status.png)](https://godoc.org/github.com/larryli/vagrantcloud.v1)

This is the documentation and guide for the Vagrant Cloud API.
You can use it to create and update boxes, versions and providers.

	import "github.com/larryli/vagrantcloud.v1"

 	api := vagrantcloud.New("--replace-your-access-token--")
	box := api.Box("yourname", "boxname")
	if box.New() == nil {
		version := box.Version(0)
		version.Version = "0.0.1"
		if version.New() == nil {
			provider := version.Provider(vagrantcloud.ProviderVirtualbox)
			provider.OriginalUrl = "http://your.box.url"
			if provider.New() == nil {
				version.Release()
			}
		}
	}
