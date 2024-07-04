package clash

import (
	"fmt"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/metacubex/mihomo/constant"
)

func urlToMetadata(rawURL string, network constant.NetWork) (addr constant.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}
	uintPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return
	}

	addr = constant.Metadata{
		NetWork: network,
		Host:    u.Hostname(),
		DstIP:   netip.Addr{},
		DstPort: uint16(uintPort),
	}
	return
}

//func urlToMetadata(rawURL string, network constant.NetWork) (addr constant.Metadata, err error) {
//	u, err := url.Parse(rawURL)
//	if err != nil {
//		return
//	}
//
//	port := u.Port()
//	if port == "" {
//		switch u.Scheme {
//		case "https":
//			port = "443"
//		case "http":
//			port = "80"
//		default:
//			return
//		}
//	}
//	addr = constant.Metadata{
//		NetWork: network,
//		Host:    u.Hostname(),
//		DstIP:   nil,
//		DstPort: port,
//	}
//	return
//}
