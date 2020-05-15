package app

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

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

type ContextManager interface {
	DBPool() *sqlx.DB
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
}

// Model 管理， 用户方法扩展
type AppContext struct {
	DB      *sqlx.DB
	Context ContextManager
}

func NewAppContext(DB *sqlx.DB, context ContextManager) *AppContext {
	return &AppContext{DB: DB, Context: context}
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
	LikeNames  []string
	LikeValue  interface{}
	Joins      []string
	LeftJoins  []string
	Eq         sq.Eq
	LtOrEq     sq.LtOrEq
	GtOrEq     sq.GtOrEq
	Wheres     []string
	Pager      bool
	PageSize   uint64
	PagePos    uint64
	ResultRef  interface{}
	ResultPage *PageResult
}

func NewCrudQuery(table string, culumns []string, resultRef interface{}) *CrudQuery {
	v := &CrudQuery{Table: table, Culumns: culumns, ResultRef: resultRef}
	return v
}

type CrudAdd struct {
	Table string
	vals  []map[string]interface{}
}

func NewCrudAdd(table string, vals []map[string]interface{}) *CrudAdd {
	return &CrudAdd{Table: table, vals: vals}
}

type CrudUpdate struct {
	Table  string
	vals   map[string]interface{}
	Eq     sq.Eq
	LtOrEq sq.LtOrEq
	GtOrEq sq.GtOrEq
}

func NewCrudUpdate(table string, vals map[string]interface{}, filter map[string]interface{}) *CrudUpdate {
	return &CrudUpdate{Table: table, vals: vals, Eq: filter}
}

// CRUD 获取单个对象
func (m *AppContext) DBGet(cg *CrudGet) error {
	sql, args, _ := sq.
		Select(cg.Culumns...).
		From(cg.Table).
		Where(cg.Filter).Limit(1).
		ToSql()

	err := m.DB.Get(cg.ResultRef, sql, args...)
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

		if cq.Wheres != nil {
			for _, where := range cq.Wheres {
				if where != "" {
					b = b.Where(where)
				}
			}
		}

		if cq.LikeValue != "" {
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

		return b
	}
	bs := sq.Select(cq.Culumns...).From(cq.Table)
	bs = filterBuilder(bs)

	// 设置分页查询参数
	if cq.Pager {
		cq.ResultPage = EmptyPageResult
		bs = bs.Offset(cq.PagePos).Limit(cq.PagePos)
	}

	// 查询数据
	sql, args, _ := bs.ToSql()
	if log.IsDebug {
		log.Debug(sql, args)
	}
	err := m.DB.Select(cq.ResultRef, sql, args...)
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
			if log.IsDebug {
				log.Debug(sql, args)
			}
			err := m.DB.Get(&total, sqlbc, argsbc...)
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
func (m *AppContext) DBAdd(ca *CrudAdd) error {
	tx, err := m.DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, valmap := range ca.vals {
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

		if log.IsDebug {
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
func (m *AppContext) DBUpdate(cu *CrudUpdate) error {
	b := sq.Update(cu.Table).SetMap(cu.vals)
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
	if log.IsDebug {
		log.Debug(sql, args)
	}
	_, err := m.DB.Exec(sql, args...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 根据id删除表数据
func (m *AppContext) DBDelete(table string, ids []string) error {
	for _, id := range ids {
		sql, args, _ := sq.Delete(table).Where(sq.Eq{"id": id}).ToSql()
		_, err := m.DB.Exec(sql, args...)
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return nil
}

// 清空表
func (m *AppContext) DBTrucate(table string) error {
	_, err := m.DB.Exec("TRUNCATE TABLE " + table)
	return err
}
