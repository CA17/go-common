package app

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"

	"github.com/ca17/go-common/common"
	"github.com/ca17/go-common/conf"
	"github.com/ca17/go-common/log"
)

// 获取数据库连接，执行一次
func GetDatabase(config *conf.DBConfig) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Passwd,
		config.Host,
		config.Port,
		config.Name)
	pool, err := sqlx.Open("mysql", dsn)
	common.Must(err)
	pool.SetMaxOpenConns(config.MaxConn)
	pool.SetMaxIdleConns(config.MaxIdle)
	return pool
}

func GetGrpcConn(config *conf.GrpcConfig) (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(config.CertFile, config.Host)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// 上下文管理， 用户管理全局对象
type ContextManager interface {
	DBPool() *sqlx.DB
	GrpConn() *grpc.ClientConn
	GetAppConfig() interface{}
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
}

type AppContext struct {
	Context ContextManager
}

func (m *AppContext) Set(key string, val interface{}) {
	m.Context.Set(key, val)
}

func (m *AppContext) Get(key string) (interface{}, bool) {
	return m.Context.Get(key)
}

func NewAppContext(context ContextManager) *AppContext {
	return &AppContext{Context: context}
}

// Webix 表格列定义
type WebixTableColumn struct {
	Id         string      `json:"id,omitempty"`
	Header     interface{} `json:"header,omitempty"`
	Headermenu interface{} `json:"headermenu,omitempty"`
	Adjust     interface{} `json:"adjust,omitempty"`
	Hidden     interface{} `json:"hidden,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	Fillspace  interface{} `json:"fillspace,omitempty"`
	Css        string      `json:"css,omitempty"`
	Template   string      `json:"template,omitempty"`
	Width      int         `json:"width,omitempty"`
	Height     int         `json:"height,omitempty"`
}

// 分页对象
type PageResult struct {
	TotalCount int64       `json:"total_count,omitempty"`
	Pos        int64       `json:"pos"`
	Data       interface{} `json:"data"`
}

// 空分页对象
var EmptyPageResult = &PageResult{
	TotalCount: 0,
	Pos:        0,
	Data:       common.EmptyList,
}

// CRUD 定义

type CrudGet struct {
	Table     string
	Culumns   []string
	Filter    map[string]interface{}
	ResultRef interface{}
}

func NewCrudGet(table string, culumns []string, filter map[string]interface{}, resultRef interface{}) *CrudGet {
	return &CrudGet{Table: table, Culumns: culumns, Filter: filter, ResultRef: resultRef}
}

type CrudFilterLike struct {
	Names []string
	Value string
}

type CrudQuery struct {
	Table      string
	Culumns    []string
	Tags       string
	LikeNames  []string
	LikeValue  interface{}
	DateRange  DateRange
	DateColumn string
	Joins      []string
	LeftJoins  []string
	Eq         sq.Eq
	LtOrEq     sq.LtOrEq
	GtOrEq     sq.GtOrEq
	OrderBy    string
	Limit      uint64
	Wheres     []string
	Pager      bool
	PageSize   uint64
	PagePos    uint64
	ResultRef  interface{}
	ResultPage *PageResult
}

func (cq *CrudQuery) SetEqValue(key, val string) {
	if key != "" && val != "" {
		if cq.Eq == nil {
			cq.Eq = sq.Eq{}
		}
		cq.Eq[key] = val
	}
}

func NewCrudQuery(table string, culumns []string, resultRef interface{}) *CrudQuery {
	v := &CrudQuery{Table: table, Culumns: culumns, ResultRef: resultRef}
	return v
}

type CrudAdd struct {
	Table string
	Vals  []map[string]interface{}
}

func NewCrudAdd(table string, vals []map[string]interface{}) *CrudAdd {
	return &CrudAdd{Table: table, Vals: vals}
}

type CrudUpdate struct {
	tx     *sql.Tx
	Table  string
	Vals   map[string]interface{}
	Eq     sq.Eq
	LtOrEq sq.LtOrEq
	GtOrEq sq.GtOrEq
}

func NewCrudUpdate(table string, vals map[string]interface{}, filter map[string]interface{}) *CrudUpdate {
	return &CrudUpdate{Table: table, Vals: vals, Eq: filter}
}

// CRUD 获取单个对象
func (m *AppContext) DBGet2(table string, culumns []string, filter map[string]interface{}, resultRef interface{}) error {
	return m.DBGet(NewCrudGet(table, culumns, filter, resultRef))
}

func (m *AppContext) DBGet(cg *CrudGet) error {
	sql, args, _ := sq.
		Select(cg.Culumns...).
		From(cg.Table).
		Where(cg.Filter).Limit(1).
		ToSql()

	err := m.Context.DBPool().Get(cg.ResultRef, sql, args...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// CRUD 查询列表
func (m *AppContext) DBQuery(cq *CrudQuery) error {
	// 查询过滤
	var filterBuilder = func(b sq.SelectBuilder) sq.SelectBuilder {

		if cq.Tags != "" {
			orwheres := []string{}
			for _, tag := range strings.Split(cq.Tags, ",") {
				orwheres = append(orwheres, fmt.Sprintf("FIND_IN_SET(\"%s\", tags)", tag))
			}
			b = b.Where(fmt.Sprintf(" ( %s ) ", strings.Join(orwheres, " or ")))
		}

		if cq.DateColumn != "" {
			if cq.DateRange.End != "" {
				b = b.Where(sq.LtOrEq{cq.DateColumn: cq.DateRange.End})
			}
			if cq.DateRange.Start != "" {
				b = b.Where(sq.GtOrEq{cq.DateColumn: cq.DateRange.Start})
			}
		}

		if cq.Joins != nil {
			for _, join := range cq.Joins {
				b = b.Join(join)
			}
		}

		if cq.LeftJoins != nil {
			for _, join := range cq.LeftJoins {
				b = b.LeftJoin(join)
			}
		}

		if cq.Wheres != nil && len(cq.Wheres) > 0 {
			for _, where := range cq.Wheres {
				if where != "" {
					b = b.Where(where)
				}
			}
		}

		if cq.LikeValue != "" && cq.LikeNames != nil {
			like := fmt.Sprintf("%s%%", cq.LikeValue)
			ormap := sq.Or{}
			for _, name := range cq.LikeNames {
				ormap = append(ormap, sq.Like{name: like})
			}
			b = b.Where(ormap)
		}

		if cq.Eq != nil {
			b = b.Where(cq.Eq)
		}

		if cq.LtOrEq != nil {
			b = b.Where(cq.LtOrEq)
		}

		if cq.GtOrEq != nil {
			b = b.Where(cq.GtOrEq)
		}

		if cq.OrderBy != "" {
			b = b.OrderBy(cq.OrderBy)
		}

		return b
	}
	bs := sq.Select(cq.Culumns...).From(cq.Table)
	bs = filterBuilder(bs)

	// 设置分页查询参数
	if cq.Pager {
		cq.ResultPage = EmptyPageResult
		bs = bs.Offset(cq.PagePos).Limit(cq.PageSize)
	} else {
		if cq.Limit > 0 {
			bs = bs.Limit(cq.Limit)
		}
	}

	// 查询数据
	sql, args, _ := bs.ToSql()
	if log.IsDebug() {
		log.Debug(sql, args)
	}
	err := m.Context.DBPool().Select(cq.ResultRef, sql, args...)
	if err != nil {
		log.Error(err)
		return err
	}

	// 封装分页结果
	if cq.Pager {
		var total int64 = 0
		if cq.PagePos == 0 {
			bc := sq.Select("count(*)").From(cq.Table)
			bc = filterBuilder(bc)
			sqlbc, argsbc, _ := bc.ToSql()
			if log.IsDebug() {
				log.Debug(sql, args)
			}
			err := m.Context.DBPool().Get(&total, sqlbc, argsbc...)
			if err != nil {
				log.Error(err)
				cq.ResultPage = EmptyPageResult
				return err
			}
		}
		cq.ResultPage = &PageResult{Data: cq.ResultRef, Pos: int64(cq.PagePos), TotalCount: total}
	}

	return nil
}

// CRUD 增加数据对象
func (m *AppContext) DBInsert(table string, vals map[string]interface{}) error {
	return m.DBInsertWithTx(nil, table, vals)
}
func (m *AppContext) DBInsertWithTx(tx *sql.Tx, table string, vals map[string]interface{}) error {
	var cols []string
	var values []interface{}
	for k, v := range vals {
		cols = append(cols, k)
		values = append(values, v)
	}
	sql, args, err := sq.Insert(table).Columns(cols...).Values(values...).ToSql()
	if err != nil {
		log.Error(err)
		return err
	}
	if log.IsDebug() {
		log.Debug(sql, args)
	}
	if tx != nil {
		_, err = tx.Exec(sql, args...)
	} else {
		_, err = m.Context.DBPool().Exec(sql, args...)
	}

	if err != nil {
		return err
	}
	return nil
}

// CRUD 增加数据对象
func (m *AppContext) DBAdd(ca *CrudAdd) error {
	tx, err := m.Context.DBPool().Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, valmap := range ca.Vals {
		var cols []string
		var values []interface{}

		for k, v := range valmap {
			cols = append(cols, k)
			values = append(values, v)
		}
		sql, args, err := sq.Insert(ca.Table).Columns(cols...).Values(values...).ToSql()
		if err != nil {
			log.Error(err)
			return err
		}

		if log.IsDebug() {
			log.Debug(sql, args)
		}

		_, err = tx.Exec(sql, args...)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// CRUD 数据更新
func (m *AppContext) DBUpdate2(table string, vals map[string]interface{}, filter map[string]interface{}) error {
	cu := NewCrudUpdate(table, vals, filter)
	return m.DBUpdate(cu)
}

// CRUD 数据更新
func (m *AppContext) DBUpdate2WithTx(tx *sql.Tx, table string, vals map[string]interface{}, filter map[string]interface{}) error {
	cu := NewCrudUpdate(table, vals, filter)
	cu.tx = tx
	return m.DBUpdate(cu)
}

func (m *AppContext) DBUpdate(cu *CrudUpdate) error {
	b := sq.Update(cu.Table).SetMap(cu.Vals)
	if cu.Eq != nil {
		b = b.Where(cu.Eq)
	}
	if cu.LtOrEq != nil {
		b = b.Where(cu.LtOrEq)
	}

	if cu.GtOrEq != nil {
		b = b.Where(cu.GtOrEq)
	}
	sql, args, _ := b.ToSql()
	if log.IsDebug() {
		log.Debug(sql, args)
	}
	var err error
	if cu.tx != nil {
		_, err = cu.tx.Exec(sql, args...)
	} else {
		_, err = m.Context.DBPool().Exec(sql, args...)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 根据id删除表数据
func (m *AppContext) DBDelete(table string, ids []string) error {
	return m.DBDeleteWithTx(nil, table, ids)
}

func (m *AppContext) DBDeleteWithTx(tx *sql.Tx, table string, ids []string) error {
	for _, id := range ids {
		sql, args, _ := sq.Delete(table).Where(sq.Eq{"id": id}).ToSql()
		var err error
		if tx != nil {
			_, err = tx.Exec(sql, args...)
		} else {
			_, err = m.Context.DBPool().Exec(sql, args...)
		}
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}

func (m *AppContext) DBDeleteWithFilter(table string, filter map[string]interface{}) error {
	return m.DBDeleteWithFilterTx(nil, table, filter)
}

func (m *AppContext) DBDeleteWithFilterTx(tx *sql.Tx, table string, filter map[string]interface{}) error {
	sql, args, _ := sq.Delete(table).Where(filter).ToSql()
	var err error
	if tx != nil {
		_, err = tx.Exec(sql, args...)
	} else {
		_, err = m.Context.DBPool().Exec(sql, args...)
	}
	if err != nil {
		log.Error(err)
	}
	return nil
}

// 清空表
func (m *AppContext) DBTrucate(table string) error {
	_, err := m.Context.DBPool().Exec("TRUNCATE TABLE " + table)
	return err
}
