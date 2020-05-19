package conductrpc

import (
	"github.com/valyala/gorpc"
)

// Client of the conductor RPC server.
type Client struct {
	client *gorpc.Client
	dispc  *gorpc.DispatcherClient
}

// NewClient creates new client will be interacting
// with server with given address.
func NewClient(address string) (c *Client) {
	c = new(Client)
	c.client = gorpc.NewTCPClient(address)

	var disp = gorpc.NewDispatcher()
	disp.AddFunc("onViewChange", nil)
	disp.AddFunc("onAddMiner", nil)
	disp.AddFunc("onAddSharder", nil)
	disp.AddFunc("onMinerReady", nil)
	disp.AddFunc("onSharderReady", nil)
	c.dispc = disp.NewFuncClient(c.client)

	return
}

// ViewChange notification.
func (c *Client) ViewChange(viewChange ViewChange) (err error) {
	_, err = c.dispc.Call("onViewChange", viewChange)
	return
}

// AddMiner notification.
func (c *Client) AddMiner(minerID MinerID) (err error) {
	_, err = c.dispc.Call("onAddMiner", minerID)
	return
}

// AddSharder notification.
func (c *Client) AddSharder(sharderID SharderID) (err error) {
	_, err = c.dispc.Call("onAddSharder", sharderID)
	return
}

// MinerReady notification.
func (c *Client) MinerReady(minerID MinerID) (err error) {
	_, err = c.dispc.Call("onMinerReady", minerID)
	return
}

// SharderReady notification.
func (c *Client) SharderReady(sharderID SharderID) (err error) {
	_, err = c.dispc.Call("onSharderReady", sharderID)
	return
}