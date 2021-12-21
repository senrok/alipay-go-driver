/*
Copyright 2020 RS4
@Author: Weny Xu
@Date: 2021/01/07 5:19
*/

package entity

import "github.com/pkg/errors"

type ErrorResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (c ErrorResponse) Failed() bool {
	return c.Code != ""
}

type AlipayResponseCommon interface {
	GetCode() string
	GetMsg() string
	GetSubMsg() string
	GetSubCode() string
}

type Common struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (c Common) GetCode() string {
	return c.Code
}

func (c Common) GetMsg() string {
	return c.Msg
}

func (c Common) GetSubCode() string {
	return c.SubCode
}

func (c Common) GetSubMsg() string {
	return c.SubMsg
}

func (c Common) Success() bool {
	return c.Code == "10000"
}

func (c Common) ErrorWrap(err error) error {
	if err != nil || !c.Success() {
		if err != nil {
			return errors.Wrap(err, c.Msg)
		}
		return errors.New(c.Msg)
	}
	return nil
}
