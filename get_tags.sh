#!/usr/bin/env bash

#check if we are running with root privileges
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# check rsc is installed
if [ ! -f /usr/local/bin/rsc ]; then
  echo "RSC is missing or not installed"
  exit 2
fi

# find ourselves
SELF=$(/usr/local/bin/rsc --rl10 --x1 'object:has(.rel:val("self")).href' cm15 index_instance_session sessions/instance)

# get the tags
/usr/local/bin/rsc --rl10 --xm '.name' cm15 by_resource /api/tags/by_resource resource_hrefs[]=${SELF}
