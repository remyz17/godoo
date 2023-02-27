package client

import (
	"fmt"

	"github.com/kolo/xmlrpc"
)

type rpcClient struct {
	common *xmlrpc.Client
	object *xmlrpc.Client
	db     *xmlrpc.Client
}

func (rpc *rpcClient) commonCall(serviceMethod string, args interface{}) (interface{}, error) {
	resp, err := rpc.call(rpc.common, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *rpcClient) dbCall(serviceMethod string, args interface{}) (interface{}, error) {
	resp, err := rpc.call(rpc.db, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *rpcClient) objectCall(serviceMethod string, args interface{}) (interface{}, error) {
	resp, err := rpc.call(rpc.object, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *rpcClient) call(x *xmlrpc.Client, serviceMethod string, args interface{}) (interface{}, error) {
	var reply interface{}
	if err := x.Call(serviceMethod, args, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func getRpcClient(port int) (*rpcClient, error) {
	common, err := xmlrpc.NewClient(fmt.Sprint("http://localhost:", port, "/xmlrpc/2/common"), nil)
	if err != nil {
		return nil, err
	}
	object, err := xmlrpc.NewClient(fmt.Sprint("http://localhost:", port, "/xmlrpc/2/object"), nil)
	if err != nil {
		return nil, err
	}
	db, err := xmlrpc.NewClient(fmt.Sprint("http://localhost:", port, "/xmlrpc/2/db"), nil)
	if err != nil {
		return nil, err
	}
	return &rpcClient{common: common, object: object, db: db}, nil
}
