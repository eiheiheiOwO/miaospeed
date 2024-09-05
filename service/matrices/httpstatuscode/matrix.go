package httpstatuscode

import (
	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/service/macros/ping"
)

type HTTPStatusCode struct {
	interfaces.HTTPStatusCodeDS
}

func (m *HTTPStatusCode) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixHTTPCode
}

func (m *HTTPStatusCode) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroPing
}

func (m *HTTPStatusCode) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*ping.Ping); ok {
		m.Values = mac.StatusCodes
	}
}
