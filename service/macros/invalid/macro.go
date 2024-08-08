package invalid

import "github.com/airportr/miaospeed/interfaces"

type Invalid struct{}

func (m *Invalid) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroInvalid
}

func (m *Invalid) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	return nil
}
