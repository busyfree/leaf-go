package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/busyfree/leaf-go/models/schema"
	"github.com/busyfree/leaf-go/util/ctxkit"
	"github.com/busyfree/leaf-go/util/db"
)

type LeafAllocDao struct {
	schema.LeafAlloc
}

func NewLeafAllocDao() *LeafAllocDao {
	return new(LeafAllocDao)
}

func (dao *LeafAllocDao) UpdateMaxId(ctx context.Context, tag string) (err error) {
	dao.BeforeInsert()
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlInsert := fmt.Sprintf("UPDATE %s SET max_id = max_id + step, update_time=? WHERE biz_tag =?", dao.TableName())
	q := db.SQLInsert(dao.TableName(), sqlInsert)
	_, err = c.ExecContext(
		ctx,
		q,
		dao.UpdatedAt,
		tag)
	return
}

func (dao *LeafAllocDao) UpdateMaxIdAndGetLeafAlloc(ctx context.Context, tag string) (err error) {
	dao.BeforeInsert()
	err = dao.UpdateMaxId(ctx, tag)
	if err != nil {
		return
	}
	err = dao.GetLeafAlloc(ctx, tag)
	if err != nil {
		return
	}
	return
}

func (dao *LeafAllocDao) UpdateMaxIdByCustomStepAndGetLeafAlloc(ctx context.Context, oldDao *LeafAllocDao) (err error) {
	dao.BeforeInsert()
	err = dao.UpdateMaxIdByCustomStep(ctx, oldDao.Step, oldDao.BizTag)
	if err != nil {
		return
	}
	err = dao.GetLeafAlloc(ctx, oldDao.BizTag)
	if err != nil {
		return
	}
	return
}

func (dao *LeafAllocDao) UpdateMaxIdByCustomStep(ctx context.Context, step int, tag string) (err error) {
	dao.BeforeInsert()
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlInsert := fmt.Sprintf("UPDATE %s SET max_id = max_id + ?, update_time=? WHERE biz_tag =?", dao.TableName())
	q := db.SQLInsert(dao.TableName(), sqlInsert)
	_, err = c.ExecContext(
		ctx,
		q,
		step,
		dao.UpdatedAt,
		tag)
	return
}

func (dao *LeafAllocDao) GetAllTags(ctx context.Context) (array []string, err error) {
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlStr := "SELECT biz_tag FROM %s"
	sqlSelect := fmt.Sprintf(sqlStr, dao.TableName())
	q := db.SQLSelect(dao.TableName(), sqlSelect)
	var rows *sql.Rows
	rows, err = c.QueryContext(ctx, q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var bizTag string
		if err := rows.Scan(
			&bizTag,
		); err != nil {
			continue
		}
		array = append(array, bizTag)
	}
	err = rows.Close()
	if err != nil {
		return
	}
	err = rows.Err()
	return
}

func (dao *LeafAllocDao) GetLeafAlloc(ctx context.Context, tag string) (err error) {
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlSelect := fmt.Sprintf("SELECT biz_tag, max_id, step FROM  %s WHERE biz_tag = ? AND deleted_at=0", dao.TableName())
	q := db.SQLSelect(dao.TableName(), sqlSelect)
	result := c.QueryRowContext(ctx, q, tag)
	if result == nil {
		err = sql.ErrNoRows
		return
	}
	err = result.Scan(&dao.BizTag, &dao.MaxId, &dao.Step)
	return
}

func (dao *LeafAllocDao) GetAllLeafAllocs(ctx context.Context) (array []*LeafAllocDao, err error) {
	c := db.Get(ctx, ctxkit.GetProjectDBName(ctx))
	sqlStr := "SELECT biz_tag, max_id, step, update_time FROM %s WHERE deleted_at=0"
	sqlSelect := fmt.Sprintf(sqlStr, dao.TableName())
	q := db.SQLSelect(dao.TableName(), sqlSelect)
	var rows *sql.Rows
	rows, err = c.QueryContext(ctx, q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		arrDao := NewLeafAllocDao()
		if err := rows.Scan(
			&arrDao.BizTag,
			&arrDao.MaxId,
			&arrDao.Step,
			&arrDao.UpdatedAt,
		); err != nil {
			continue
		}
		arrDao.AfterLoad()
		array = append(array, arrDao)
	}
	err = rows.Close()
	if err != nil {
		return
	}
	err = rows.Err()
	return
}
