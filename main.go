// rs_tag
// written in go (c) Millard Technical Services Ltd
//   Tagger allows listing, adding and removing tags on the current instance
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
// === Usage
//    rs_tag (--list, -l | --add, -a TAG | --remove, -r TAG)
//
//    Options:
//      --list, -l           List current server tags
//      --add, -a TAG        Add tag named TAG
//      --remove, -r TAG     Remove tag named TAG
//      --format, -f FMT     Output format: json, yaml, text
//	--verbose, -v        Display debug information
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
	"encoding/json"

	"github.com/rightscale/rsc/cm15"
	"gopkg.in/alecthomas/kingpin.v2"
)

// for testing
var (
	osStdout io.Writer = os.Stdout
	verbose  = kingpin.Flag("verbose", "Display debug information").Short('v').Bool()
	format   = kingpin.Flag("format", "Output format: json, text").Short('f').String()
	list     = kingpin.Flag("list", "List current server tags").Short('l').Bool()
	tagAdd   = kingpin.Flag("add","Add tag named TAG").Short('a').PlaceHolder("TAG").Bool()
	tagRem   = kingpin.Flag("remove","Remove tag named TAG").Short('r').PlaceHolder("TAG").Bool()
	tag	 = kingpin.Arg("TAG","TAG to be added or removed").String()
	Keys     []string
)

func main() {

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Neil Millard")
	kingpin.CommandLine.Help = `  Rightscale(tm) tagger (rs_tag)
	  Tagger allows listing, adding and removing tags on the current instance
	  This version written for RL10
   === Examples:
   Retrieve all tags:
     rs_tag --list
     rs_tag -l

   Add tag 'a_tag' to instance:
     rs_tag --add a_tag
     rs_tag -a a_tag

   Remove tag 'a_tag':
     rs_tag --remove a_tag
     rs_tag -r a_tag
     `
	kingpin.Parse()
	// check we have something to do
	action := string("")
	if *tagRem {
		action = "remove"
		checkTag(tag)
	} else if *tagAdd {
		action = "add"
		checkTag(tag)
	} else if *list {
		action = "list"
	} else {
		fail("Missing argument, rs_tag --help for additional information")
	}

	// Create our RightLink10 client
	client, err := cm15.NewRL10()
	if err != nil {
		fail("Failed to Create a client: %v\nTry elevating privilege (sudo/runas) before invoking this command.", err.Error())
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
		if *verbose {
			fmt.Fprintf(osStdout, "Instance: %s\n", instanceEntry.Name)
		}
	}
	// extract the HREF (api url) for this instance
	instanceHref := []string{getHref(instanceEntry)}
	// everything setup, lets see what we need to do.
	if *verbose {
		fmt.Fprintf(osStdout, "Action: %v\n",action)
	}
	switch action {
	case "remove":
		tags := []string{*tag}
		// create a Locator for multi_delete
		tagLocator := client.TagLocator("/api/tags/multi_delete")
		// Multi_add function expects an array of strings
		err := tagLocator.MultiDelete(instanceHref,tags)
		if err != nil {
			fail("Failed to remove TAGS from Instance: %v\n", err.Error())
		}
		fmt.Fprintf(osStdout, "Successfully removed tag %s\n",*tag)

	case "add":
		tags := []string{*tag}
		// create a Locator for multi_add
		tagLocator := client.TagLocator("/api/tags/multi_add")
		// Multi_add function expects an array of strings
		err := tagLocator.MultiAdd(instanceHref,tags)
		if err != nil {
			fail("Failed to add TAGS to Instance: %v\n", err.Error())
		}
		fmt.Fprintf(osStdout, "Successfully added tag %s\n",*tag)

	case "list":
		// create a Locator for by_resource
		tagLocator := client.TagLocator("/api/tags/by_resource")
		// ByResource function expects an array of strings
		tagData, err := tagLocator.ByResource(instanceHref)
		if err != nil {
			fail("Failed to retrieve TAGS Instance: %v\n", err.Error())
		}
		Keys = processTags(tagData)
	}

	if *verbose {
		fmt.Fprintf(osStdout, "Output: %v\n",*format)
		fmt.Fprintf(osStdout, "No Keys: %v\n",len(Keys))
	}
	if len(Keys) > 0 {
		switch *format {
		case "text":
			outputText(Keys)

		default:
			outputJson(Keys)
		}
	}
}

func checkTag(tag *string) {
	if len(*tag) < 3 {
		fail("Add tag failed: No tags supplied")
	}
	return
}

// text output,
func outputText(keys []string) {
	for tagentry := range keys {
		fmt.Fprintf(osStdout, "%v\n",keys[tagentry])
	}
}

// json output
func outputJson(keys []string) {
	tags, _ := json.MarshalIndent(keys,"","  ")
	fmt.Println(string(tags))
}

// processTags
// expects output from "github.com/rightscale/rsc/cm15" - tagLocator.ByResource
func processTags(tagData []map[string]interface{}) []string  {

	tags := tagData[0]["tags"].([]interface{})
	keys := make([]string, 0, len(tags))
	for _,value := range tags {
		for _,v := range value.(map[string]interface{}) {
			keys = append(keys,v.(string))
		}
	}
	return keys
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