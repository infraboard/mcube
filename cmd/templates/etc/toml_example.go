package etc

// TOMLExampleTemplate todo
const TOMLExampleTemplate = `[app]
name = "{{.Name}}"
host = "127.0.0.1"
port = "8080"
key  = "this is your app key"

[mysql]
host = "127.0.0.1"
port = "3306"
username = "{{.Name}}"
password = "xxxx"
database = "{{.Name}}"

[mongodb]
endpoints = ["127.0.0.1:27017"]
username = "{{.Name}}"
password = ""
database = "{{.Name}}"

[log]
level = "debug"
path = "logs"
format = "text"
to = "stdout"`
