package model

import (
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-microservices/utils/utils"
	"strings"
)

var _ GlobalVariablesModel = (*customGlobalVariablesModel)(nil)

type (
	// GlobalVariablesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGlobalVariablesModel.
	GlobalVariablesModel interface {
		globalVariablesModel
		withSession(session sqlx.Session) GlobalVariablesModel
		// GetTableName 获取表名
		GetTableName() string
		// GetCount 根据条件获取数量
		GetCount(ctx context.Context, ex goqu.Expression) (int64, error)
		// FindList 根据条件获取列表，排序：map[string]int{"字段":0/1(0-升序(ASC)；1-降序(DESC))}；分页：[]uint{页码，每页条数}
		FindList(ctx context.Context, ex goqu.Expression, optionalParams ...any) (*[]GlobalVariables, error)
		// FindOnly 根据条件获取单条数据，0-升序(ASC)；1-降序(DESC)
		FindOnly(ctx context.Context, ex goqu.Expression, order ...map[string]int) (*GlobalVariables, error)
		// InsertOnly 插入单条数据
		InsertOnly(ctx context.Context, row *GlobalVariables, tx ...*sql.Tx) (sql.Result, error)
		// BatchInsert 批量插入
		BatchInsert(ctx context.Context, rows []*GlobalVariables, tx ...*sql.Tx) (sql.Result, error)
		// UpdateByEx 根据条件更新
		UpdateByEx(ctx context.Context, record goqu.Record, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error)
		// DeleteByEx 根据条件删除数据
		DeleteByEx(ctx context.Context, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error)
	}

	customGlobalVariablesModel struct {
		*defaultGlobalVariablesModel
	}
)

// NewGlobalVariablesModel returns a model for the database table.
func NewGlobalVariablesModel(conn sqlx.SqlConn) GlobalVariablesModel {
	return &customGlobalVariablesModel{
		defaultGlobalVariablesModel: newGlobalVariablesModel(conn),
	}
}

func (m *customGlobalVariablesModel) withSession(session sqlx.Session) GlobalVariablesModel {
	return NewGlobalVariablesModel(sqlx.NewSqlConnFromSession(session))
}

// GetTableName 获取表名
func (m *customGlobalVariablesModel) GetTableName() string {
	return utils.SetTable(m.table)
}

// GetCount 根据条件获取数量
func (m *customGlobalVariablesModel) GetCount(ctx context.Context, ex goqu.Expression) (int64, error) {
	query, _, err := goqu.Dialect("mysql").Select(goqu.COUNT(1)).From(utils.SetTable(m.table)).Where(ex).ToSQL()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.conn.QueryRowCtx(ctx, &resp, query)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return 0, err
	}
	return resp, nil
}

// FindList 根据条件获取列表，排序：map[string]int{"字段":0/1(0-升序(ASC)；1-降序(DESC))}；分页：[]uint{页码，每页条数}
func (m *customGlobalVariablesModel) FindList(ctx context.Context, ex goqu.Expression, optionalParams ...any) (*[]GlobalVariables, error) {
	sql := goqu.Dialect("mysql").Select(&GlobalVariables{}).From(utils.SetTable(m.table)).Where(ex)
	if len(optionalParams) > 0 {
		for _, param := range optionalParams {
			// 排序
			if v, ok := param.(map[string]int); ok {
				for key, value := range v {
					if value > 0 {
						sql = sql.OrderAppend(goqu.C(key).Desc())
					} else {
						sql = sql.OrderAppend(goqu.C(key).Asc())
					}
				}
			}
			// 分页
			if v, ok := param.([]uint); ok {
				if len(v) == 2 {
					sql = sql.Offset((v[0] - 1) * v[1]).Limit(v[1])
				}
			}
		}
	}
	query, _, err := sql.ToSQL()
	query = strings.ReplaceAll(query, `""`, `"`)
	if err != nil {
		return nil, err
	}
	var resp []GlobalVariables
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, err
	}
	return &resp, nil
}

// FindOnly 根据条件获取单条数据，0-升序(ASC)；1-降序(DESC)
func (m *customGlobalVariablesModel) FindOnly(ctx context.Context, ex goqu.Expression, order ...map[string]int) (*GlobalVariables, error) {
	sql := goqu.Dialect("mysql").Select(&GlobalVariables{}).From(utils.SetTable(m.table)).Where(ex)
	if len(order) > 0 {
		for key, value := range order[0] {
			if value > 0 {
				sql = sql.OrderAppend(goqu.C(key).Desc())
			} else {
				sql = sql.OrderAppend(goqu.C(key).Asc())
			}
		}
	}
	query, _, err := sql.Limit(1).ToSQL()
	query = strings.ReplaceAll(query, `""`, `"`)
	if err != nil {
		return nil, err
	}
	var resp GlobalVariables
	err = m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// InsertOnly 插入单条数据
func (m *customGlobalVariablesModel) InsertOnly(ctx context.Context, row *GlobalVariables, tx ...*sql.Tx) (sql.Result, error) {
	query, _, err := goqu.Dialect("mysql").Insert(utils.SetTable(m.table)).Rows(row).ToSQL()
	if err != nil {
		return nil, err
	}
	var result sql.Result
	if len(tx) > 0 {
		result, err = tx[0].ExecContext(ctx, query)
	} else {
		result, err = m.conn.ExecCtx(ctx, query)
	}
	return result, err
}

// BatchInsert 批量插入
func (m *customGlobalVariablesModel) BatchInsert(ctx context.Context, rows []*GlobalVariables, tx ...*sql.Tx) (sql.Result, error) {
	query, _, err := goqu.Dialect("mysql").Insert(utils.SetTable(m.table)).Rows(rows).ToSQL()
	if err != nil {
		return nil, err
	}
	var result sql.Result
	if len(tx) > 0 {
		result, err = tx[0].ExecContext(ctx, query)
	} else {
		result, err = m.conn.ExecCtx(ctx, query)
	}
	return result, err
}

// UpdateByEx 根据条件更新
func (m *customGlobalVariablesModel) UpdateByEx(ctx context.Context, record goqu.Record, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error) {
	query, _, err := goqu.Dialect("mysql").Update(utils.SetTable(m.table)).Set(record).Where(ex).ToSQL()
	if err != nil {
		return nil, err
	}
	var result sql.Result
	if len(tx) > 0 {
		result, err = tx[0].ExecContext(ctx, query)
	} else {
		result, err = m.conn.ExecCtx(ctx, query)
	}
	return result, err
}

// DeleteByEx 根据条件删除数据
func (m *customGlobalVariablesModel) DeleteByEx(ctx context.Context, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error) {
	query, _, err := goqu.Dialect("mysql").Delete(utils.SetTable(m.table)).Where(ex).ToSQL()
	if err != nil {
		return nil, err
	}
	var result sql.Result
	if len(tx) > 0 {
		result, err = tx[0].ExecContext(ctx, query)
	} else {
		result, err = m.conn.ExecCtx(ctx, query)
	}
	return result, err
}
