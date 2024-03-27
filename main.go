package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jghiloni/custom-resource-type/framework"
	"github.com/jghiloni/github-tags-resource/resource"
)

var Version string = "0.0.0-local"
var BuildRef string = "local"

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile)
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		tw := tabwriter.NewWriter(os.Stdout, 7, 4, 2, '\t', 0)
		fmt.Fprintf(tw, "Version:\t%s\n", Version)
		fmt.Fprintf(tw, "Build Ref:\t%s\n", BuildRef)
		tw.Flush()
		return
	}

	impl := &resource.GithubTagsResourceType[resource.SourceInfo, resource.Version, any, any]{}

	rt := framework.NewResourceType[resource.SourceInfo, resource.Version, any, any](impl)
	if err := rt.Run(os.Args...); err != nil {
		log.Fatal(err)
	}
}
