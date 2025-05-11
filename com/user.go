package com

import "go-microservices/dao/mysql/model"

// CheckUser 检测用户
func (c *ComFunc) CheckUser(user *model.User) bool {
	return user.Id > 0
}
