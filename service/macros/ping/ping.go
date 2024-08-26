package ping

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	urllib "net/url"
	"strings"
	"time"

	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/preconfigs"
	"github.com/airportr/miaospeed/utils"
	"github.com/airportr/miaospeed/utils/structs"
)

func pingViaTrace(ctx context.Context, p interfaces.Vendor, url string) (uint16, uint16, error) {
	transport := &http.Transport{
		//Dial: func(string, string) (net.Conn, error) {
		//	return p.DialTCP(ctx, url, interfaces.ROptionsTCP)
		//},
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return p.DialTCP(ctx, url, interfaces.ROptionsTCP)
		},
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			// for version prior to tls1.3, the handshake will take 2-RTTs,
			// plus, majority server supports tls1.3, so we set a limit here
			MinVersion: tls.VersionTLS13,
			RootCAs:    preconfigs.MiaokoRootCAPrepare(),
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	tlsStart := int64(0)
	tlsEnd := int64(0)
	writeStart := int64(0)
	writeEnd := int64(0)
	trace := &httptrace.ClientTrace{
		TLSHandshakeStart: func() {
			tlsStart = time.Now().UnixMilli()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsEnd = time.Now().UnixMilli()
			if err != nil {
				tlsEnd = 0
			}
		},
		GotFirstResponseByte: func() {
			writeEnd = time.Now().UnixMilli()
		},
		WroteHeaders: func() {
			writeStart = time.Now().UnixMilli()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	connStart := time.Now().UnixMilli()
	if resp, err := transport.RoundTrip(req); err != nil {
		return 0, 0, err
	} else {
		defer resp.Body.Close()
		connEnd := time.Now().UnixMilli()
		utils.DBlackhole(!strings.HasPrefix(url, "https:"), connEnd-writeEnd, writeEnd-tlsEnd, tlsEnd-tlsStart, tlsStart-connStart)
		if !strings.HasPrefix(url, "https:") {
			return uint16(writeStart - connStart), uint16(writeEnd - connStart), nil
		}
		if resp.TLS != nil && resp.TLS.HandshakeComplete {
			// use payload rtt
			return uint16(writeEnd - tlsEnd), uint16(writeEnd - connStart), nil
			// return uint16(tlsEnd - tlsStart), uint16(writeEnd - connStart), nil
		}
		// https rtt is also avalible
		if writeEnd != 0 {
			return 0, uint16(writeEnd - connStart), nil
		}
		return 0, 0, fmt.Errorf("cannot extract payload from response")
	}
}

func pingViaNetCat(ctx context.Context, p interfaces.Vendor, url string) (uint16, uint16, error) {
	purl, _ := urllib.Parse(url)
	data := purl.EscapedPath()
	if purl.RawQuery != "" {
		data += "?" + purl.Query().Encode()
	}

	data = structs.X(preconfigs.NETCAT_HTTP_PAYLOAD, data, purl.Hostname(), utils.VERSION)

	connStart := time.Now()
	conn, err := p.DialTCP(ctx, url, interfaces.ROptionsTCP)
	if err != nil || conn == nil {
		return 0, 0, fmt.Errorf("cannot dial remote address")
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(time.Second * 6))
	reader := bufio.NewReader(conn)

	// prewrite to ensure tcp conn is established
	// httpStartReq1 := time.Now().UnixMilli()
	if _, err := conn.Write([]byte(data)); err != nil {
		return 0, 0, fmt.Errorf("cannot write payload to remote")
	}
	_, _ = reader.ReadByte()
	rtt1 := time.Since(connStart).Milliseconds()
	//_, _, _ = reader.ReadLine()
	for reader.Buffered() > 0 {
		_, _, _ = reader.ReadLine()
	}

	httpStartReq2 := time.Now()
	if _, err := conn.Write([]byte(data)); err != nil {
		return 0, 0, fmt.Errorf("cannot write payload to remote")
	}
	_, _ = reader.ReadByte()
	//_, _, _ = reader.ReadLine()

	rtt2 := time.Since(httpStartReq2).Milliseconds()
	for reader.Buffered() > 0 {
		_, _, _ = reader.ReadLine()
	}
	utils.DBlackholef("http response time1: %d, %d", uint16(rtt2), uint16(rtt1))
	return uint16(rtt2), uint16(rtt1), nil
}

func ping(obj *Ping, p interfaces.Vendor, url string, withAvg uint16, timeout uint) {
	var totalMS []uint16
	var totalMSRTT []uint16
	if p == nil {
		obj.RTT = 0
		obj.Request = 0
		obj.RTTList = totalMSRTT
		obj.RequestList = totalMS
		obj.MaxRTT = 0
		obj.MaxRequest = 0
	}

	// 20 次足够了
	if withAvg < 1 || withAvg > 1000 {
		withAvg = 1
	}

	for len(totalMS) < int(withAvg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		delayRTT, delay := uint16(0), uint16(0)
		if strings.HasPrefix(url, "https:") {
			delayRTT, delay, _ = pingViaTrace(ctx, p, url)
		} else {
			delayRTT, delay, _ = pingViaNetCat(ctx, p, url)
		}
		obj.MaxRTT = structs.Max(obj.MaxRTT, delayRTT)
		obj.MaxRequest = structs.Max(obj.MaxRequest, delay)
		totalMSRTT = append(totalMSRTT, delayRTT)
		totalMS = append(totalMS, delay)
		cancel()
	}

	obj.RTT = calcAvgPing(totalMSRTT)
	obj.Request = calcAvgPing(totalMS)
	obj.RTTList = totalMSRTT
	obj.RequestList = totalMS
	obj.RTTSD = calcStdDevPing(totalMSRTT)
	obj.RequestSD = calcStdDevPing(totalMS)
}
