# rl10-rs_tag
[RS_TAG](https://github.com/rightscale/right_link/blob/master/scripts/tagger.rb) command for RightLink10
-------
Requires RSC [https://github.com/rightscale/rsc]
This is installed as part of Rightlink10 also.

-------
    rs_agent:mime_include_url=https://rightlink.rightscale.com/rll/10.5.2/rightlink.boot.sh	
    rs_agent:type=right_link_lite

to compile

install go
----------
    wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz
    sudo tar -C /usr/local/ -xzf go1.7.1.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    
setup env
---------
uses best practices [docs](https://golang.org/doc/code.html#Organization)

    # install git (required for go get)
    sudo yum -y install git
    # set a working directory for our go project
    export GOPATH=$HOME/work
    # create project folder 
    mkdir -p work/src/github.com/user/
    cd work/src/github.com/user
    # install dependancy
    go get "github.com/rightscale/rsc/cm15"
    # clone this repo
    git clone https://github.com/neilmillard/rl10-rs_tag
    cd rl10-rs_tag/
    
to cross compile
----------------
    cd rl10-rs_tag/
    git pull
    #set GOOS and GOARCH to be the values for the target operating system and architecture.
    env GOOS=linux GOARCH=arm go build -v rs_tag
    
or just install
---------------
    #go install
    $GOPATH/bin/rs_tag