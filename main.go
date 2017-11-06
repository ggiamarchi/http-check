package main

import (
	"github.com/ggiamarchi/http-check/api"
)

func main() {
	api.RunServer(api.CreateEngine())
}
