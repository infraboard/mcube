package root

// GitIgnreTemplate todo
const GitIgnreTemplate = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
# Test binary, built with {{.Backquote}}go test -c{{.Backquote}}
*.test
# Output of the go coverage tool, specifically when used with LiteIDE
*.out
*.idea
*vendor

cover.out
coverage.txt

etc/{{.Name}}.toml
etc/{{.Name}}.env
etc/ssl
{{.Name}}`
