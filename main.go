package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/carabiner-dev/databom/internal"
)

func main() {
	//	r := internal.Renderer{}
	n := time.Now()
	c := &internal.CreationInfo{
		SpecVersion: "3.0.1",
		CreatedBy: []string{
			"https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
		},
		Created: &n,
		PreNode: internal.PreNode{
			ID:   "_:creationinfo",
			Type: "CreationInfo",
		},
		RootedNode: internal.RootedNode{
			RootElement: []string{},
		},
	}

	e := &internal.Envelope{
		Context: "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		Graph: internal.Graph{
			c,
		},
	}

	//fmt.Printf("%+v", e.Graph[0])
	d, err := json.Marshal(e)
	if err != nil {
		fmt.Println("Err: %v", err)
	}
	fmt.Printf("%s", string(d))
	//r.Render(e, os.Stdout)
}
