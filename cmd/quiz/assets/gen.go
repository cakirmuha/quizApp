// +build ignore

package main

import (
	"log"

	"quizApp/cmd/quiz/assets"
	"github.com/shurcooL/vfsgen"
)

//go:generate go run -tags=dev gen.go
func main() {
	err := vfsgen.Generate(assets.Assets, vfsgen.Options{
		Filename:     "assets_vfsdata.go",
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
