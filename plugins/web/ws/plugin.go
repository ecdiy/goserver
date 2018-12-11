package ws

import "github.com/ecdiy/goserver/plugins/web"
import (
	"github.com/gorilla/websocket"
	"time"
	"net/http"
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins/sp"
	"encoding/json"
	"github.com/ecdiy/goserver/plugins"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Bean struct {
	msgId, id string
	ws        *websocket.Conn
}

var Online []*Bean
//var OnlineUser = make(map[string]*websocket.Conn)

func init() {
	web.RegisterWebPlugin("WebSocket", func(ele *utils.Element) func(c *gin.Context) {
		ws := &Ws{verify: plugins.GetRef(ele, "Verify").(plugins.BaseFun),
			UserIdName: ele.Attr("SocketIdName", "UserId")}
		spName, spExt := ele.AttrValue("SpName")
		if spExt {
			ws.spName = spName
			ws.sp = &sp.WebSp{}
			ws.sp.Init(ele)
		}
		return ws.WsHandler
	})
}

type Ws struct {
	spName, UserIdName string
	sp                 *sp.WebSp
	verify             plugins.BaseFun
}

// 处理ws请求
func (ws *Ws) WsHandler(c *gin.Context) {
	ws.verify(utils.NewParam(c))
	userFace, uExt := c.Get(ws.UserIdName)
	if !uExt {
		seelog.Error("用户没有登录")
		return
	}
	userId := fmt.Sprint(userFace)
	pingTicker := time.NewTicker(wsupgrader.HandshakeTimeout - time.Second)
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	bean := &Bean{msgId: userId, ws: conn, id: userId + fmt.Sprint(time.Now().UnixNano())}
	Online = append(Online, bean)
	isOk := true
	go func() {
		for {
			<-pingTicker.C
			if isOk {
				conn.WriteMessage(websocket.PingMessage, []byte{})
			} else {
				pingTicker.Stop()
			}
		}
	}()
	for {
		// recieve
		_, bs, err := conn.ReadMessage()
		if err != nil {
			isOk = false
			for i, v := range Online {
				if v.id == bean.id {
					Online = append(Online[:i], Online[i+1:]...)
				}
			}
			seelog.Error("read message error", userId)
			break
		}
		wb := utils.NewParam(c)
		je := json.Unmarshal(bs, &wb.Param)
		if je != nil {
			seelog.Error("param error", je)
			//	web.Context.Set("param", web.Param)
		}
		ws.sp.SpExec(ws.spName, wb)
		if len(wb.Out) >= 1 {
			bs, err := json.Marshal(wb.Out)
			if err == nil {
				conn.WriteMessage(websocket.TextMessage, bs)
			} else {
				seelog.Error("输出数据错误，不能format JSON", wb.Out)
			}
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("{}"))
		}
	}
}

func WrMsg(msgId string, msg []byte) {
	for _, b := range Online {
		if b.msgId == msgId {
			b.ws.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
