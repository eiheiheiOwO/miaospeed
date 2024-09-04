package clash

import (
	"context"
	"fmt"
	"github.com/airportr/miaospeed/utils"
	"net"

	"github.com/airportr/miaospeed/interfaces"
	"github.com/metacubex/mihomo/component/resolver"
	"github.com/metacubex/mihomo/constant"
)

type Clash struct {
	proxy constant.Proxy
}

func init() {
	if resolver.DisableIPv6 {
		resolver.DisableIPv6 = false
	}
}

func (c *Clash) Proxy() constant.Proxy {
	return c.proxy
}

func (c *Clash) Type() interfaces.VendorType {
	return interfaces.VendorClash
}

func (c *Clash) Status() interfaces.VendorStatus {
	if c == nil || c.proxy == nil {
		return interfaces.VStatusNotReady
	}

	return interfaces.VStatusOperational
}

func (c *Clash) Build(proxyName string, proxyInfo string) interfaces.Vendor {
	if c == nil {
		c = &Clash{}
	}
	c.proxy = extractFirstProxy(proxyName, proxyInfo)
	return c
}

func (c *Clash) DialTCP(ctx context.Context, url string, network interfaces.RequestOptionsNetwork) (net.Conn, error) {
	if c == nil || c.proxy == nil {
		return nil, fmt.Errorf("should call Build() before run")
	}

	addr, err := urlToMetadata(url, constant.TCP)
	if err != nil {
		return nil, fmt.Errorf("cannot build tcp context")
	}
	conn, err := c.proxy.DialContext(ctx, &addr)
	if err != nil {
		utils.DLogf("cannot dialTCP: %s | proxy=%s | vendor=Clash | err=%s", url, c.proxy.Name(), err.Error())
	}
	return conn, err
}

func (c *Clash) DialUDP(ctx context.Context, url string) (net.PacketConn, error) {
	if c == nil || c.proxy == nil {
		return nil, fmt.Errorf("should call Build() before run")
	}

	addr, err := urlToMetadata(url, constant.UDP)
	if err != nil {
		return nil, fmt.Errorf("cannot build udp context")
	}
	conn, err := c.proxy.DialUDP(&addr)
	if err != nil {
		utils.DLogf("cannot dialUDP: %s | proxy=%s | vendor=Clash | err=%s", url, c.proxy.Name(), err.Error())
	}
	return conn, err

}
func (c *Clash) ProxyInfo() interfaces.ProxyInfo {
	if c == nil || c.proxy == nil {
		return interfaces.ProxyInfo{}
	}

	return interfaces.ProxyInfo{
		Name:    c.proxy.Name(),
		Address: c.proxy.Addr(),
		Type:    interfaces.Parse(c.proxy.Type().String()),
	}
}
