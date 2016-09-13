# rl10-rs_tag
[RS_TAG](https://github.com/rightscale/right_link/blob/master/scripts/tagger.rb) command for RightLink10
-------
Requires RSC [https://github.com/rightscale/rsc]
This is installed as part of Rightlink10 also.

-------
    rs_agent:mime_include_url=https://rightlink.rightscale.com/rll/10.5.2/rightlink.boot.sh	
    rs_agent:type=right_link_lite

-----------------

install go
----------
    wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz
    sudo tar -C /usr/local/ -xzf go1.7.1.linux-amd64.tar.gz
    
setup env
---------
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/work
    mkdir -p work/src/github.com/user/
    cd work/src/github.com/user
    go get "github.com/rightscale/rsc/cm15"
    git clone https://github.com/neilmillard-msm/rl10-rs_tag
    cd rl10-rs_tag/
    
    git pull
    go install
    $GOPATH/bin/rl10-rs_tag