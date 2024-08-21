package maxrttping

import (
	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/service/macros/ping"
)

type MaxRttPing struct {
	interfaces.MaxRTTDS
}

func (m *MaxRttPing) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixMAXRTTPing
}

func (m *MaxRttPing) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *MaxRttPing) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*ping.Ping); ok {
		m.Value = mac.MaxRTT
	}
}
