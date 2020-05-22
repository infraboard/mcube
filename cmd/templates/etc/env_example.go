package etc

// EnvExampleTemplate todo
const EnvExampleTemplate = `export D_MYSQL_HOST="127.0.0.1"
export D_MYSQL_PORT=3306
export D_MYSQL_USERNAME="{{.Name}}"
export D_MYSQL_PASSWORD="xxxx"
export D_MYSQL_DATABASE="{{.Name}}"`
