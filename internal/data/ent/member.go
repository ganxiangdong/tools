// Code generated by ent, DO NOT EDIT.

package ent

import (
	_ "fmt"

	_ "strings"
)

// Member 测试文件
type Member struct {
	// ID of the ent.
	// 编号
	ID int `json:"id,omitempty"`
	// RoleID holds the value of the "role_id" field.
	RoleID int `json:"role_id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// UnionID holds the value of the "UnionID" field.
	UnionID string `json:"UnionID,omitempty"`
	// 状态
	Status int `json:"status,omitempty"`
	// 创建时间
	CreateTime int `json:"create_time,omitempty"`
	// 更新时间
	UpdateTime int `json:"update_time,omitempty"`
	// 最后登录时间
	LastLoginTime int `json:"last_login_time,omitempty"`
	// 最后登录IP
	LastLoginIP  int `json:"last_login_ip,omitempty"`
}