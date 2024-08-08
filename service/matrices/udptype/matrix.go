package udptype

import (
	"github.com/AiportR/miaospeed/interfaces"
	"github.com/AiportR/miaospeed/service/macros/udp"
)

type UDPType struct {
	interfaces.UDPTypeDS
}

func (m *UDPType) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixUDPType
}

func (m *UDPType) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroUDP
}

func (m *UDPType) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*udp.Udp); ok {
		m.Value = mac.NATType
	}
}
