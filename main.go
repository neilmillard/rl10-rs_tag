// rs_tag
// written in go (c) Millard Technical Services Ltd
//   Tagger allows listing, adding and removing tags on the current instance and
//   querying for instances with a given set of tags
//
// === Examples:
//   Retrieve all tags:
//     rs_tag --list
//     rs_tag -l
//
//   Add tag 'a_tag' to instance:
//     rs_tag --add a_tag
//     rs_tag -a a_tag
//
//   Remove tag 'a_tag':
//     rs_tag --remove a_tag
//     rs_tag -r a_tag
//
//   Retrieve instances with any of the tags in a set each tag is a separate argument:
//     rs_tag --query "a_tag" "b:machine=tag" "c_tag with space"
//     rs_tag -q "a_tag" "b:machine=tag" "c_tag with space"
//
// === Usage
//    rs_tag (--list, -l | --add, -a TAG | --remove, -r TAG | --query, -q TAG[s])
//
//    Options:
//      --list, -l           List current server tags
//      --add, -a TAG        Add tag named TAG
//      --remove, -r TAG     Remove tag named TAG
//      --query, -q TAG[s]   Query for instances that have any of the TAG[s]
//                           with TAG being quoted if it contains spaces in it's value
//      --die, -e            Exit with error if query/list fails
//      --format, -f FMT     Output format: json, yaml, text
//      --verbose, -v        Display debug information
//      --help:              Display help
//      --version:           Display version information
//      --timeout, -t SEC    Custom timeout (default 180 sec)
//
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rightscale/rsc/cm15"
)

// for testing
var osStdout io.Writer = os.Stdout

func main() {
	// Create our RightLink10 client
	client, err := cm15.NewRL10()
	if err != nil {
		fail("Failed to Create a client: %v\n", err.Error())
	}
	if err := client.CanAuthenticate(); err != nil {
		fail("Unable to create connection to agent: %s", err)
	}
	// get the session for our RL10 reverse proxy connection
	sessionLocator := client.SessionLocator("/api/sessions/instance")
	// get our instance attached to the session (should be the calling instance i.e. this server)
	instanceEntry, err := sessionLocator.IndexInstanceSession()
	if err != nil {
		fail("Failed to retrieve session Instance: %v\n", err.Error())
	} else {
		fmt.Fprintf(osStdout, "Instance: %s\n", instanceEntry.Name)
	}
	// extract the HREF (api url) for this instance
	instanceHref := []string{getHref(instanceEntry)}
	// create a Locator for by_resource
	tagLocator := client.TagLocator("/api/tags/by_resource")
	// ByResource function expects an array of strings
	tags, err := tagLocator.ByResource(instanceHref)
	if err != nil {
		fail("Failed to retrieve TAGS Instance: %v\n", err.Error())
	}

	fmt.Fprintln(osStdout, "Tags:")
	for _, value := range tags {
		fmt.Fprintf(osStdout, "%v : %v\n", value)
	}

}

// Get the href of an audit entry from the Links attribute by inspecting the self link
func getHref(instance *cm15.Instance) string {
	var href string
	for _, link := range instance.Links {
		if link["rel"] == "self" {
			href = link["href"]
			break
		}
	}
	return href
}

// Print error message and exit with code 1
func fail(format string, v ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Println(fmt.Sprintf(format, v...))
	os.Exit(1)
}