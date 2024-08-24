package sdhttp

import (
	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/service/macros/ping"
)

type SDHTTP struct {
	interfaces.SDHTTPDS
}

func (m *SDHTTP) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixSDHTTP
}

func (m *SDHTTP) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *SDHTTP) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*ping.Ping); ok {
		m.Value = mac.RequestSD
	}
}
