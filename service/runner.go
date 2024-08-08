package service

import (
	"github.com/AiportR/miaospeed/interfaces"
	"github.com/AiportR/miaospeed/service/macros"
	"github.com/AiportR/miaospeed/utils/structs"
)

func ExtractMacrosFromMatrices(matrices []interfaces.SlaveRequestMatrix) []interfaces.SlaveRequestMacroType {
	macroTypes := structs.NewSet[interfaces.SlaveRequestMacroType]()
	for _, matrix := range matrices {
		macroTypes.Add(matrix.MacroJob())
	}
	return structs.Filter(macroTypes.Digest(), func(m interfaces.SlaveRequestMacroType) bool {
		return macros.Find(m).Type() != interfaces.MacroInvalid
	})
}
