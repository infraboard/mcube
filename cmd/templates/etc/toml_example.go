package etc

// TOMLExampleTemplate todo
const TOMLExampleTemplate = `[app]
name = "{{.Name}}"
host = "0.0.0.0"
port = "8050"
key  = "this is your app key"

[mysql]
host = "xxx"
port = "3306"
username = "{{.Name}}"
password = "xxxx"
database = "{{.Name}}"

[log]
level = "debug"
path = "logs"
format = "text"
to = "stdout"`
