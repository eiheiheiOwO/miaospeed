package ping

import (
	"github.com/airportr/miaospeed/interfaces"
)

type Ping struct {
	RTT     uint16
	Request uint16
}

func (m *Ping) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *Ping) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	//var pingFunc pingFuncType = nil
	//if r.Configs.UseClashPing {
	//	pingFunc = pingViaClash
	//}
	m.RTT, m.Request = ping(proxy, r.Configs.PingAddress, r.Configs.PingAverageOver, int(r.Configs.TaskRetry), r.Configs.TaskTimeout)
	return nil
}
