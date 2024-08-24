package service

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/airportr/miaospeed/interfaces"
	"github.com/airportr/miaospeed/preconfigs"
	"github.com/airportr/miaospeed/utils"
	"github.com/airportr/miaospeed/utils/structs"
	"github.com/gorilla/websocket"

	"github.com/airportr/miaospeed/service/matrices"
	"github.com/airportr/miaospeed/service/taskpoll"
)

type WsHandler struct {
	Serve func(http.ResponseWriter, *http.Request)
}

func (wh *WsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if wh.Serve != nil {
		wh.Serve(rw, r)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitServer() {
	if utils.GCFG.Binder == "" {
		utils.DErrorf("MiaoSpeed Server | Cannot listening the binder, bind=%s", utils.GCFG.Binder)
		os.Exit(1)
	}

	utils.DWarnf("MiaoSpeed Server | Start Listening, bind=%s", utils.GCFG.Binder)
	wsHandler := WsHandler{
		Serve: func(rw http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(rw, r, nil)
			if err != nil {
				utils.DErrorf("MiaoServer Test | Socket establishing error, error=%s", err.Error())
				return
			}
			defer conn.Close()
			// Verify the websocket path
			if !utils.GCFG.ValidateWSPath(r.URL.Path) {
				conn.WriteJSON(&interfaces.SlaveResponse{
					Error: "invalid websocket path",
				})
				utils.DWarnf("MiaoServer Test | websocket path error, error=%s", "invalid websocket path")
				return
			}
			var poll *taskpoll.TaskPollController

			batches := structs.NewAsyncMap[string, bool]()
			cancel := func() {
				if poll != nil {
					for id := range batches.ForEach() {
						poll.Remove(id, taskpoll.TPExitInterrupt)
					}
				}
			}

			defer cancel()
			for {
				sr := interfaces.SlaveRequest{}
				_, r, err := conn.NextReader()
				if err == nil {
					// unsafe, ensure jsoniter package receives frequently security updates.
					err = jsoniter.NewDecoder(r).Decode(&sr)
					// 原方案
					//err = json.NewDecoder(r).Decode(&sr)
					if err == io.EOF {
						// One value is expected in the message.
						err = io.ErrUnexpectedEOF
					}
				}

				if err != nil {
					if !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "reset by peer") {
						utils.DErrorf("MiaoServer Test | Task receiving error, error=%s", err.Error())
					}
					return
				}
				verified := utils.GCFG.VerifyRequest(&sr)
				utils.DLogf("MiaoServer Test | Receive Task, name=%s invoker=%v matrices=%v payload=%d verify=%v", sr.Basics.ID, sr.Basics.Invoker, sr.Options.Matrices, len(sr.Nodes), verified)

				// verify token
				if !verified {
					conn.WriteJSON(&interfaces.SlaveResponse{
						Error: "cannot verify the request, please check your token",
					})
					return
				}
				sr.Challenge = ""

				// verify invoker
				if !utils.GCFG.InWhiteList(sr.Basics.Invoker) {
					conn.WriteJSON(&interfaces.SlaveResponse{
						Error: "the bot id is not in the whitelist",
					})
					return
				}

				// find all matrices
				mtrx := matrices.FindBatchFromEntry(sr.Options.Matrices)

				// extra macro from the matrices
				macros := ExtractMacrosFromMatrices(mtrx)

				// select poll
				if structs.Contains(macros, interfaces.MacroSpeed) {
					if utils.GCFG.NoSpeedFlag {
						conn.WriteJSON(&interfaces.SlaveResponse{
							Error: "speedtest is disabled on backend",
						})
						return
					}
					poll = SpeedTaskPoll
				} else {
					poll = ConnTaskPoll
				}
				utils.DLogf("MiaoServer Test | Receive Task, name=%s poll=%s", sr.Basics.ID, poll.Name())

				// build testing item
				item := poll.Push((&TestingPollItem{
					id:       utils.RandomUUID(),
					name:     sr.Basics.ID,
					request:  &sr,
					matrices: sr.Options.Matrices,
					macros:   macros,
					onProcess: func(self *TestingPollItem, idx int, result interfaces.SlaveEntrySlot) {
						conn.WriteJSON(&interfaces.SlaveResponse{
							ID:               self.ID(),
							MiaoSpeedVersion: utils.VERSION,
							Progress: &interfaces.SlaveProgress{
								Record:  result,
								Index:   idx,
								Queuing: poll.AwaitingCount(),
							},
						})
					},
					onExit: func(self *TestingPollItem, exitCode taskpoll.TaskPollExitCode) {
						batches.Del(self.ID())
						conn.WriteJSON(&interfaces.SlaveResponse{
							ID:               self.ID(),
							MiaoSpeedVersion: utils.VERSION,
							Result: &interfaces.SlaveTask{
								Request: sr,
								Results: self.results.ForEach(),
							},
						})
					},
					// 计算权重
					calcWeight: func(self *TestingPollItem) uint {
						return 1
						//if poll.Name() == "SpeedPoll" {
						//	nodeNum := len(self.request.Nodes)
						//	w := nodeNum / 10
						//	if w == 0 {
						//		return 1
						//	} else {
						//		return uint(w)
						//	}
						//} else {
						//	return 1
						//}
					},
				}).Init())
				// onstart
				conn.WriteJSON(&interfaces.SlaveResponse{
					ID:               item.ID(),
					MiaoSpeedVersion: utils.VERSION,
					Progress: &interfaces.SlaveProgress{
						Queuing: poll.AwaitingCount(),
					},
				})
				batches.Set(item.ID(), true)
			}
		},
	}

	server := http.Server{
		Handler:   &wsHandler,
		TLSConfig: preconfigs.MakeSelfSignedTLSServer(),
	}

	if strings.HasPrefix(utils.GCFG.Binder, "/") {
		unixListener, err := net.Listen("unix", utils.GCFG.Binder)
		if err != nil {
			utils.DErrorf("MiaoServer Launch | Cannot listen on unixsocket %s, error=%s", utils.GCFG.Binder, err.Error())
			os.Exit(1)
		}
		server.Serve(unixListener)
	} else {
		netListener, err := net.Listen("tcp", utils.GCFG.Binder)
		if err != nil {
			utils.DErrorf("MiaoServer Launch | Cannot listen on socket %s, error=%s", utils.GCFG.Binder, err.Error())
			os.Exit(1)
		}
		if utils.GCFG.MiaoKoSignedTLS {
			server.ServeTLS(netListener, "", "")
		} else {
			server.Serve(netListener)
		}

	}
}

func CleanUpServer() {
	if strings.HasPrefix(utils.GCFG.Binder, "/") {
		os.Remove(utils.GCFG.Binder)
	}
}
