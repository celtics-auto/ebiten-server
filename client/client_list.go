package client

import "net"

type ClientsMap map[net.Addr]Client

func NewMap() ClientsMap {
	return ClientsMap{}
}

func (cm ClientsMap) Add(c *Client) {
	cm[c.Address] = *c
}

func (cm ClientsMap) Disconnect(addr net.Addr) {
	cm[addr].conn.Close()
	delete(cm, addr)
}

func (cm ClientsMap) FindClient(addr net.Addr) (*Client, bool) {
	c, exists := cm[addr]
	return &c, exists
}
