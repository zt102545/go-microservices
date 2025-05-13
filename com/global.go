package com

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"go-microservices/dao/mysql/model"
)

// region 结构体定义

type AvatarConfig struct {
	FilePath   string
	HttpPath   string
	CheckPath  bool
	NoImageNum int64
}

// endregion

// region 接口实现

// GetGlobalVariable 根据id获取全局配置
func (c *ComFunc) GetGlobalVariable(ctx context.Context, id int64) (*model.GlobalVariables, error) {

	key := fmt.Sprintf("global_variable_%v", id)
	v, exist := c.Cache.Get(key)
	if exist {
		value, ok := v.(string)
		if ok {
			one := &model.GlobalVariables{}
			err := jsonx.UnmarshalFromString(value, one)
			if err == nil {
				return one, nil
			}
		}
	}

	one, err := c.GlobalVariablesModel.FindOne(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	str, err := jsonx.MarshalToString(one)
	if err != nil {
		return nil, err
	}
	c.Cache.Set(key, str)

	return one, nil
}

// endregion

// region 私有方法
// endregion
