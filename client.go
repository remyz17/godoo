package client

import (
	"strconv"

	"github.com/remyz17/godoo/internal/utils"
)

type Client struct {
	rpc    *rpcClient
	config *odooConfig
}

func GetClient() (*Client, error) {
	var err error
	var conf *odooConfig
	var rpc *rpcClient
	var httpPort int
	conf, err = GetOdooConfig()
	if err != nil {
		return nil, err
	}
	httpPort, err = conf.getHttpPort()
	if err != nil {
		return nil, err
	}
	rpc, err = getRpcClient(httpPort)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: rpc, config: conf}, nil
}

func (c *Client) Version() (float64, error) {
	var err error
	var reply interface{}
	var version float64
	reply, err = c.rpc.commonCall("version", nil)
	if err != nil {
		return 0, err
	}
	data := reply.(map[string]interface{})
	version, err = strconv.ParseFloat(data["server_serie"].(string), 64)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func (c *Client) ListDatabases() ([]string, error) {
	adminPasswd, err := c.config.getAdminPasswd()
	if err != nil {
		return nil, err
	}
	reply, err := c.rpc.dbCall("list", []interface{}{adminPasswd})
	if err != nil {
		return nil, err
	}
	resp := reply.([]interface{})
	dbs := make([]string, len(resp))
	for i, v := range resp {
		dbs[i] = v.(string)
	}
	return dbs, nil
}

func (c *Client) DatabaseExists(database string) (bool, error) {
	reply, err := c.ListDatabases()
	if err != nil {
		return false, err
	}
	_, found := utils.Find(reply, database)
	if !found {
		return false, nil
	}
	return true, nil
}

func (c *Client) DuplicateDatabase(from, to string) (bool, error) {
	adminPasswd, err := c.config.getAdminPasswd()
	if err != nil {
		return false, err
	}
	reply, err := c.rpc.dbCall("duplicate_database", []interface{}{adminPasswd, from, to})
	if err != nil {
		return false, err
	}
	return reply.(bool), nil
}

func (c *Client) DropDatabase(database string) (bool, error) {
	adminPasswd, err := c.config.getAdminPasswd()
	if err != nil {
		return false, err
	}
	reply, err := c.rpc.dbCall("drop", []interface{}{adminPasswd, database})
	if err != nil {
		return false, err
	}
	return reply.(bool), nil
}

func (c *Client) DumpDatabase(database string) (string, error) {
	adminPasswd, err := c.config.getAdminPasswd()
	if err != nil {
		return "", err
	}
	reply, err := c.rpc.dbCall("dump", []interface{}{adminPasswd, database, "zip"})
	if err != nil {
		return "", err
	}
	return reply.(string), nil
}

func (c *Client) RestoreDatabase(database string, dump string) (interface{}, error) {
	adminPasswd, err := c.config.getAdminPasswd()
	if err != nil {
		return "", err
	}
	reply, err := c.rpc.dbCall("restore", []interface{}{adminPasswd, database, dump})
	if err != nil {
		return "", err
	}
	return reply, nil
}
