package root

// MainTemplate todo
const MainTemplate = `package main

import (
	"{{.PKG}}/cmd"
)

func main() {
	cmd.Execute()
}`
