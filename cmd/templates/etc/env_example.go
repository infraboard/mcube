package etc

// EnvExampleTemplate todo
const EnvExampleTemplate = `export MCUBE_APP_NAME="{{.Name}}"
export MCUBE_APP_HOST="127.0.0.1"
export MCUBE_APP_PORT=8050
export MCUBE_APP_KEY="defaut app key"
export MCUBE_LOG_LEVEL="info"
export MCUBE_LOG_PATH="logs"
export MCUBE_LOG_TO="stdout"
export MCUBE_MYSQL_HOST="127.0.0.1"
export MCUBE_MYSQL_PORT=3306
export MCUBE_MYSQL_USERNAME="{{.Name}}"
export MCUBE_MYSQL_PASSWORD=""
export MCUBE_MYSQL_DATABASE="{{.Name}}"`
