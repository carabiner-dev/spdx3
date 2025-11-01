package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/carabiner-dev/databom/internal/spdx3"
	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
)

func main() {
	//	r := internal.Renderer{}
	n := time.Now()
	c := &core.CreationInfo{
		SpecVersion: "3.0.1",
		CreatedBy: []types.Node{
			&core.Person{
				Node: core.Node{
					PreNode: base.PreNode{
						ID: "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
					},
				},
			},
		},
		Created: &n,
		PreNode: base.PreNode{
			ID:   "_:creationinfo",
			Type: "CreationInfo",
		},
		RootedNode: types.RootedNode{
			RootElement: []types.Node{},
		},
	}

	e := &spdx3.Envelope{
		Context: "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		Graph: spdx3.Graph{
			c,
		},
	}

	//fmt.Printf("%+v", e.Graph[0])
	d, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Printf("%s", string(d))
	//r.Render(e, os.Stdout)
}
