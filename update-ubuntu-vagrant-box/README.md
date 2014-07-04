Update Ubuntu Vagrant Box
=========================

update https://cloud-images.ubuntu.com/vagrant/ to https://vagrantcloud.com/larryli

Get:

	go get -ldflags "-s -w" github.com/larryli/vagrantcloud.v1/update-ubuntu-vagrant-box

Test:

	update-ubuntu-vagrant-box --username="yourname" --test

Update:

	update-ubuntu-vagrant-box --username="yourname" --token="--replace-your-access-token--"


