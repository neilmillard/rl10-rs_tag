#!/usr/bin/env bash

#check if we are running with root privileges
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# check rsc is installed - could also be in /opt/bin/rsc
if [ ! -f /usr/local/bin/rsc ]; then
  echo "RSC is missing or not installed"
  exit 2
fi

if ! /usr/local/bin/rsc rl10 actions 2>/dev/null | grep --ignore-case --quiet /rll/login/control; then
  echo "This script must be run on a RightLink 10.5 or newer instance"
  exit 1
fi

# find ourselves - we might be able to use $RS_SELF_HREF
SELF=$(/usr/local/bin/rsc --rl10 --x1 'object:has(.rel:val("self")).href' cm15 index_instance_session sessions/instance)

# get the tags
/usr/local/bin/rsc --rl10 --xm '.name' cm15 by_resource /api/tags/by_resource resource_hrefs[]=${SELF}

# add a tag "rs_login:state=user"
/usr/local/bin/rsc --rl10 cm15 multi_add /api/tags/multi_add resource_hrefs[]=$RS_SELF_HREF tags[]=rs_login:state=user

# Remove rs_login:state=user tag
/usr/local/bin/rsc --rl10 cm15 multi_delete /api/tags/multi_delete resource_hrefs[]=$RS_SELF_HREF tags[]=rs_login:state=user