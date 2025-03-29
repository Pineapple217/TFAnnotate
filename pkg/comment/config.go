package comment

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const configFile = "annotate.hcl"

type Config struct {
	Annotations []Annotation `hcl:"annotation,block"`
}

type Annotation struct {
	Name    string  `hcl:",label"`
	Module  *Module `hcl:"module,block"`
	Values  []Value `hcl:"value,block"`
	Comment string  `hcl:"comment"`
}

type Value struct {
	Name   string `hcl:",label"`
	Target string `hcl:"target"`
}

type Module struct {
	Source string `hcl:"source"`
}

func GetConfig(dir string) Config {
	var config Config
	err := hclsimple.DecodeFile(dir+configFile, nil, &config)
	if err != nil {
		panic(err)
	}

	return config
}
