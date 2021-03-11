package etc

// TOMLExampleTemplate todo
const TOMLExampleTemplate = `[app]
name = "{{.Name}}"
key  = "this is your app key"

[http]
host = "127.0.0.1"
port = "8848"

[grpc]
host = "127.0.0.1"
port = "18848"

[keyauth]
host = "127.0.0.1"
port = "18080"
client_id = ""
client_secret = ""

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
