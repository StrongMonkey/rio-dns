package main

import (
	"github.com/rancher/norman/generator"
	"github.com/rancher/rio-dns/types/apis/some.api.group/v1"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := generator.DefaultGenerate(v1.Schemas, "github.com/rancher/rio-dns/types", false, nil); err != nil {
		logrus.Fatal(err)
	}
}
