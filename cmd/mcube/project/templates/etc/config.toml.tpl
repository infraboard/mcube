[app]
name = "{{.Name}}"
key  = "this is your app key"

[http]
host = "127.0.0.1"
port = "8848"

[grpc]
host = "127.0.0.1"
port = "18848"

{{ if $.EnableKeyauth }}
[keyauth]
host = "{{.Keyauth.Host}}"
port = "{{.Keyauth.Port}}"
client_id = "{{.Keyauth.ClientID}}"
client_secret = "{{.Keyauth.ClientSecret}}"
{{ end }}

{{ if $.EnableMySQL }}
[mysql]
host = "{{.MySQL.Host}}"
port = "{{.MySQL.Port}}"
database = "{{.MySQL.Database}}"
username = "{{.MySQL.UserName}}"
password = "{{.MySQL.Password}}"
{{ end }}

{{ if $.EnableMongoDB }}
[mongodb]
endpoints = {{.MongoDB.Endpoints}}
username = "{{.MongoDB.UserName}}"
password = "{{.MongoDB.Password}}"
database = "{{.MongoDB.Database}}"
{{ end }}

[log]
level = "debug"
path = "logs"
format = "text"
to = "stdout"