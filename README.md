Vagrant Cloud API
=================

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
