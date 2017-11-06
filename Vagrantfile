# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV["LC_ALL"] = "en_US.UTF-8"

Vagrant.configure(2) do |config|

    config.vm.box = "ubuntu/trusty64"
    config.vm.hostname = 'http-check'

    config.vm.provider 'virtualbox' do |vb|
        vb.customize ['modifyvm', :id, '--memory', '1024']
    end

    config.vm.provision "shell", privileged: false, inline: <<-SHELL
        sudo rm -rf /usr/local/go $HOME/gopath
        sudo apt-get update
        sudo apt-get install -y git mercurial
        wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz 2>/dev/null
        sudo tar -C /usr/local -xzf go1.8.1.linux-amd64.tar.gz
        mkdir $HOME/gopath
        echo 'export GOROOT=/usr/local/go'               >> .profile
        echo 'export GOPATH=$HOME/gopath'                >> .profile
        echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> .profile
    SHELL

    config.vm.provision "shell", privileged: false, inline: <<-SHELL
        source $HOME/.profile

        go get gopkg.in/gin-gonic/gin.v1

        mkdir -p $GOPATH/src/github.com/ggiamarchi/
        cd $GOPATH/src/github.com/ggiamarchi/
        ln -s /vagrant http-check

        go build
    SHELL


    config.vm.network "forwarded_port", guest: 5000, host: 5000
end
