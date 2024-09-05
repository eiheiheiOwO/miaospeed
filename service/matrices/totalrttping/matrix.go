package totalrttping

import (
	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/service/macros/ping"
)

type TotalRTT struct {
	interfaces.TotalRTTDS
}

type TotalHTTP struct {
	interfaces.TotalHTTPDS
}

func (m *TotalRTT) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixTotalRTTPing
}

func (m *TotalRTT) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *TotalRTT) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*ping.Ping); ok {
		m.Values = mac.RTTList
	}
}

func (m *TotalHTTP) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixTotalHTTPPing
}

func (m *TotalHTTP) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *TotalHTTP) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*ping.Ping); ok {
		m.Values = mac.RequestList
	}
}
