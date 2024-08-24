package matrices

import (
	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/utils/structs"

	"github.com/airportr/miaospeed/service/matrices/averagespeed"
	"github.com/airportr/miaospeed/service/matrices/httpping"
	"github.com/airportr/miaospeed/service/matrices/inboundgeoip"
	"github.com/airportr/miaospeed/service/matrices/invalid"
	"github.com/airportr/miaospeed/service/matrices/maxrttping"
	"github.com/airportr/miaospeed/service/matrices/maxspeed"
	"github.com/airportr/miaospeed/service/matrices/outboundgeoip"
	"github.com/airportr/miaospeed/service/matrices/persecondspeed"
	"github.com/airportr/miaospeed/service/matrices/rttping"
	"github.com/airportr/miaospeed/service/matrices/scripttest"
	"github.com/airportr/miaospeed/service/matrices/sdhttp"
	"github.com/airportr/miaospeed/service/matrices/sdrtt"
	"github.com/airportr/miaospeed/service/matrices/udptype"
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
	interfaces.MatrixMAXRTTPing: func() interfaces.SlaveRequestMatrix {
		return &maxrttping.MaxRttPing{}
	},
	interfaces.MatrixSDRTT:  func() interfaces.SlaveRequestMatrix { return &sdrtt.SDRTT{} },
	interfaces.MatrixSDHTTP: func() interfaces.SlaveRequestMatrix { return &sdhttp.SDHTTP{} },
	//interfaces.MatrixMAXHTTPPing: func() interfaces.SlaveRequestMatrix {
	//	return &maxrttping.MaxRttPing{}
	//},
	//interfaces.MatrixTotalRTTPing: func() interfaces.SlaveRequestMatrix {
	//	return &scripttest.ScriptTest{}
	//},
	//interfaces.MatrixTotalHTTPPing: func() interfaces.SlaveRequestMatrix {
	//	return &scripttest.ScriptTest{}
	//},
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
