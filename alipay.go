/*
Copyright 2020 RS4
@Author: Weny Xu
@Date: 2021/01/07 4:20
*/

package alipay

import (
	"context"

	"github.com/WenyXu/better-alipay-go/config"
	m "github.com/WenyXu/better-alipay-go/m"
	"github.com/WenyXu/better-alipay-go/options"
)

type Service interface {
	Options() options.Options
	Request(method string, param m.M, response interface{}, opts ...options.Option) (err error)
	MakeParam(method string, param m.M, opts ...options.Option) (data m.M, err error)
}

type service struct {
	opts options.Options
}

// MakeParam
func (s service) MakeParam(method string, param m.M, opts ...options.Option) (data m.M, err error) {
	copyOpts := s.opts.Copy()
	// setup options
	for _, o := range opts {
		o(&copyOpts)
	}
	return m.CombineMakeMapEndpointFunc(
		config.SetMethod(method),
		config.SetPublicParam(copyOpts.Config),
		config.SetOptionalParam(copyOpts.Config),
		config.SetBizContent(param),
		config.SignParam(copyOpts.Config),
	)
}

// Do Request
func (s service) Request(method string, param m.M, response interface{}, opts ...options.Option) (err error) {
	copyOpts := s.opts.Copy()
	ctx, cancel := context.WithCancel(copyOpts.Context)
	defer cancel()
	// setup options
	for _, o := range opts {
		o(&copyOpts)
	}

	// make request
	req, err := copyOpts.MakeReq(ctx, method, param, copyOpts.Config)
	if err != nil {
		cancel()
		return err
	}

	// before hooks
	for _, f := range copyOpts.Before {
		ctx = f(ctx, req)
	}

	// do request
	resp, err := copyOpts.Transport.RoundTrip(req)
	if err != nil {
		cancel()
		return err
	}
	defer resp.Body.Close()

	// after hooks
	for _, f := range copyOpts.After {
		ctx = f(ctx, resp)
	}

	// dec response
	err = copyOpts.Dec(ctx, resp, response)
	if err != nil {
		cancel()
		return err
	}
	return
}

func (s service) Options() options.Options {
	return s.opts
}

// New return a empty service, with following global config
//
// you can use:
// options.SetDefaultTransport(f http.RoundTripper),
// options.SetDefaultDecFunc(f DecFunc),
// options.SetDefaultMakeReqFunc(f MakeReqFunc),
// options.SetDefaultLogger(f logger.Logger),
// options.SetDefaultLocation(f Option),
// to modify these global config.
//
// 	opt := Options{
//		Transport: DefaultTransport,
//		Context:   context.Background(),
//		MakeReq:   DefaultMakeReqFunc,
//		Dec:       DefaultDecFunc,
//		Logger:    DefaultLogger,
//	}
func New(opts ...options.Option) Service {
	return newService(func(s *service) {
		s.opts = options.NewOptions(opts...)
	})
}

// Default return a default service
//
//				s := Default()
// 				// same as
// 				s := New(
//					options.DefaultLocation,
//					options.DefaultCharset(),
//					options.DefaultFormat(),
//					options.DefaultVersion(),
//				)
func Default(opts ...options.Option) Service {
	return newService(func(s *service) {
		s.opts = options.DefaultOptions(
			append(
				append(
					[]options.Option{},
					options.DefaultLocation,
					options.DefaultCharset(),
					options.DefaultFormat(),
					options.DefaultVersion(),
				), opts...,
			)...,
		)
	})
}

func newService(setup func(s *service)) Service {
	service := service{}
	setup(&service)
	return service
}
