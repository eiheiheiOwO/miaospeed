package matrices

import (
	"github.com/AiportR/miaospeed/interfaces"
	"github.com/AiportR/miaospeed/utils/structs"

	"github.com/AiportR/miaospeed/service/matrices/averagespeed"
	"github.com/AiportR/miaospeed/service/matrices/httpping"
	"github.com/AiportR/miaospeed/service/matrices/inboundgeoip"
	"github.com/AiportR/miaospeed/service/matrices/invalid"
	"github.com/AiportR/miaospeed/service/matrices/maxspeed"
	"github.com/AiportR/miaospeed/service/matrices/outboundgeoip"
	"github.com/AiportR/miaospeed/service/matrices/persecondspeed"
	"github.com/AiportR/miaospeed/service/matrices/rttping"
	"github.com/AiportR/miaospeed/service/matrices/scripttest"
	"github.com/AiportR/miaospeed/service/matrices/udptype"
)

var registeredList = map[interfaces.SlaveRequestMatrixType]func() interfaces.SlaveRequestMatrix{
	interfaces.MatrixHTTPPing: func() interfaces.SlaveRequestMatrix {
		return &httpping.HTTPPing{}
	},
	interfaces.MatrixRTTPing: func() interfaces.SlaveRequestMatrix {
		return &rttping.RTTPing{}
	},
	interfaces.MatrixUDPType: func() interfaces.SlaveRequestMatrix {
		return &udptype.UDPType{}
	},
	interfaces.MatrixAverageSpeed: func() interfaces.SlaveRequestMatrix {
		return &averagespeed.AverageSpeed{}
	},
	interfaces.MatrixMaxSpeed: func() interfaces.SlaveRequestMatrix {
		return &maxspeed.MaxSpeed{}
	},
	interfaces.MatrixPerSecondSpeed: func() interfaces.SlaveRequestMatrix {
		return &persecondspeed.PerSecondSpeed{}
	},
	interfaces.MatrixInboundGeoIP: func() interfaces.SlaveRequestMatrix {
		return &inboundgeoip.InboundGeoIP{}
	},
	interfaces.MatrixOutboundGeoIP: func() interfaces.SlaveRequestMatrix {
		return &outboundgeoip.OutboundGeoIP{}
	},
	interfaces.MatrixScriptTest: func() interfaces.SlaveRequestMatrix {
		return &scripttest.ScriptTest{}
	},
}

func Find(matrixType interfaces.SlaveRequestMatrixType) interfaces.SlaveRequestMatrix {
	if newFn, ok := registeredList[matrixType]; ok && newFn != nil {
		return newFn()
	}

	return &invalid.Invalid{}
}

func FindBatch(macroTypes []interfaces.SlaveRequestMatrixType) []interfaces.SlaveRequestMatrix {
	return structs.Map(macroTypes, func(m interfaces.SlaveRequestMatrixType) interfaces.SlaveRequestMatrix {
		return Find(m)
	})
}

func FindBatchFromEntry(macroTypes []interfaces.SlaveRequestMatrixEntry) []interfaces.SlaveRequestMatrix {
	return structs.Map(macroTypes, func(m interfaces.SlaveRequestMatrixEntry) interfaces.SlaveRequestMatrix {
		return Find(m.Type)
	})
}
