# === Synopsis:
#   RightScale Tagger (rs_tag) - (c) 2016 Millard Technical Services Ltd
#   Statically linked to avoid external deps.
#
#   Tagger allows listing, adding and removing tags on the current instance and
#   querying for instances with a given set of tags
#   This depends on RightLink10 on the box for Instance facing calls only.
#
# === Examples:
#   Retrieve all tags:
#     rs_tag --list
#     rs_tag -l
#
#   Add tag 'a_tag' to instance:
#     rs_tag --add a_tag
#     rs_tag -a a_tag
#
#   Remove tag 'a_tag':
#     rs_tag --remove a_tag
#     rs_tag -r a_tag
#
#   Retrieve instances with any of the tags in a set each tag is a separate argument:
#     rs_tag --query "a_tag" "b:machine=tag" "c_tag with space"
#     rs_tag -q "a_tag" "b:machine=tag" "c_tag with space"
#
# === Usage
#    rs_tag (--list, -l | --add, -a TAG | --remove, -r TAG | --query, -q TAG[s])
#
#    Options:
#      --list, -l           List current server tags
#      --add, -a TAG        Add tag named TAG
#      --remove, -r TAG     Remove tag named TAG
#      --query, -q TAG[s]   Query for instances that have any of the TAG[s]
#                           with TAG being quoted if it contains spaces in it's value
#      --die, -e            Exit with error if query/list fails
#      --format, -f FMT     Output format: json, yaml, text
#      --verbose, -v        Display debug information
#      --help:              Display help
#      --version:           Display version information
#      --timeout, -t SEC    Custom timeout (default 180 sec)
#
