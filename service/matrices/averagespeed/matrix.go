package averagespeed

import (
	"github.com/AiportR/miaospeed/interfaces"
	"github.com/AiportR/miaospeed/service/macros/speed"
)

type AverageSpeed struct {
	interfaces.AverageSpeedDS
}

func (m *AverageSpeed) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixAverageSpeed
}

func (m *AverageSpeed) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroSpeed
}

func (m *AverageSpeed) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*speed.Speed); ok {
		m.Value = mac.AvgSpeed
	}
}
