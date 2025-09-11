package datasource

import (
	"context"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &dataSource{
	Provider:    PROVIDER_MYSQL,
	Host:        "127.0.0.1",
	Port:        3306,
	DB:          application.Get().Name(),
	AutoMigrate: false,
	Debug:       false,
	Trace:       true,

	SkipDefaultTransaction: false,
	DryRun:                 false,
	PrepareStmt:            true,
}

type dataSource struct {
	ioc.ObjectImpl
	Provider    PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"PROVIDER"`
	Host        string   `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port        int      `json:"port" yaml:"port" toml:"port" env:"PORT"`
	DB          string   `json:"database" yaml:"database" toml:"database" env:"DB"`
	Username    string   `json:"username" yaml:"username" toml:"username" env:"USERNAME"`
	Password    string   `json:"password" yaml:"password" toml:"password" env:"PASSWORD"`
	AutoMigrate bool     `json:"auto_migrate" yaml:"auto_migrate" toml:"auto_migrate" env:"AUTO_MIGRATE"`
	Debug       bool     `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG"`
	Trace       bool     `toml:"trace" json:"trace" yaml:"trace"  env:"TRACE"`

	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `toml:"skip_default_transaction" json:"skip_default_transaction" yaml:"skip_default_transaction"  env:"SKIP_DEFALT_TRANSACTION"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool `toml:"full_save_associations" json:"full_save_associations" yaml:"full_save_associations"  env:"FULL_SAVE_ASSOCIATIONS"`
	// DryRun generate sql without execute
	DryRun bool `toml:"dry_run" json:"dry_run" yaml:"dry_run"  env:"DRY_RUN"`
	// PrepareStmt executes the given query in cached statement
	PrepareStmt bool `toml:"prepare_stmt" json:"prepare_stmt" yaml:"prepare_stmt"  env:"PREPARE_STMT"`
	// DisableAutomaticPing
	DisableAutomaticPing bool `toml:"disable_automatic_ping" json:"disable_automatic_ping" yaml:"disable_automatic_ping"  env:"DISABLE_AUTOMATIC_PING"`
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool `toml:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating" yaml:"disable_foreign_key_constraint_when_migrating"  env:"DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING"`
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool `toml:"ignore_relationships_when_migrating" json:"ignore_relationships_when_migrating" yaml:"ignore_relationships_when_migrating"  env:"IGNORE_RELATIONSHIP_WHEN_MIGRATING"`
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool `toml:"disable_nested_transaction" json:"disable_nested_transaction" yaml:"disable_nested_transaction"  env:"DISABLE_NESTED_TRANSACTION"`
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool `toml:"allow_global_update" json:"allow_global_update" yaml:"allow_global_update"  env:"ALL_GLOBAL_UPDATE"`
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool `toml:"query_fields" json:"query_fields" yaml:"query_fields"  env:"QUERY_FIELDS"`
	// CreateBatchSize default create batch size
	CreateBatchSize int `toml:"create_batch_size" json:"create_batch_size" yaml:"create_batch_size"  env:"CREATE_BATCH_SIZE"`
	// TranslateError enabling error translation
	TranslateError bool `toml:"translate_error" json:"translate_error" yaml:"translate_error"  env:"TRANSLATE_ERROR"`

	db  *gorm.DB
	log *zerolog.Logger
}

func (m *dataSource) Name() string {
	return AppName
}

func (i *dataSource) Priority() int {
	return 699
}

func (m *dataSource) Init() error {
	m.log = log.Sub(m.Name())
	db, err := gorm.Open(m.Dialector(), &gorm.Config{
		SkipDefaultTransaction:                   m.SkipDefaultTransaction,
		FullSaveAssociations:                     m.FullSaveAssociations,
		DryRun:                                   m.DryRun,
		PrepareStmt:                              m.PrepareStmt,
		DisableAutomaticPing:                     m.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: m.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         m.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 m.DisableNestedTransaction,
		AllowGlobalUpdate:                        m.AllowGlobalUpdate,
		Logger:                                   newGormCustomLogger(m.log),
	})
	if err != nil {
		return err
	}

	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable gorm trace")
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			return err
		}
	}

	if m.Debug {
		db = db.Debug()
	}

	m.db = db
	return nil
}

// 关闭数据库连接
func (m *dataSource) Close(ctx context.Context) {
	if m.db == nil {
		return
	}

	d, err := m.db.DB()
	if err != nil {
		m.log.Error().Msgf("获取db error, %s", err)
		return
	}
	if err := d.Close(); err != nil {
		m.log.Error().Msgf("close db error, %s", err)
	}
}

// 从上下文中获取事物, 如果获取不到则直接返回 无事物的DB对象
func (m *dataSource) GetTransactionOrDB(ctx context.Context) *gorm.DB {
	db := GetTransactionFromCtx(ctx)
	if db != nil {
		return db.WithContext(ctx)
	}
	return m.db.WithContext(ctx)
}

func (m *dataSource) Dialector() gorm.Dialector {
	switch m.Provider {
	case PROVIDER_POSTGRES:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			m.Host,
			m.Username,
			m.Password,
			m.DB,
			m.Port,
		)
		return postgres.Open(dsn)
	case PROVIDER_SQLITE:
		return sqlite.Open(m.DB)
	default:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.Username,
			m.Password,
			m.Host,
			m.Port,
			m.DB,
		)
		return mysql.Open(dsn)
	}
}
