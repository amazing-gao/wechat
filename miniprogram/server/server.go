package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/amazing-gao/wechat/v2/miniprogram/context"
	"github.com/amazing-gao/wechat/v2/miniprogram/message"
	"github.com/amazing-gao/wechat/v2/util"
)

// Server struct
type (
	Server struct {
		*context.Context
		messageHandler MessageHandler
	}

	MessageHandler func(mixMessage *message.MiniProgramMixMessage) *message.Reply
)

// miniProgramMixMessage 小程序回调的消息结构
type miniProgramMixMessage struct {
	message.CommonToken

	MsgID int64 `xml:"MsgId"`

	// 文本消息
	Content string `xml:"Content"`

	// 图片消息
	PicURL  string `xml:"PicUrl"`
	MediaID string `xml:"MediaId"`

	// 小程序卡片消息
	Title        string `xml:"Title"`
	AppID        string `xml:"AppId"`
	PagePath     string `xml:"PagePath"`
	ThumbURL     string `xml:"ThumbUrl"`
	ThumbMediaID string `xml:"ThumbMediaId"`

	// 进入会话事件
	Event       message.EventType `xml:"Event"`
	SessionFrom string            `xml:"SessionFrom"`

	// 用户操作订阅通知弹窗消息回调
	List []message.SubscribeMessageList `xml:"-" json:"List"`

	// 用户操作订阅通知弹窗消息回调
	SubscribeMsgPopupEvent struct {
		List []message.SubscribeMessageList `xml:"List"`
	} `xml:"SubscribeMsgPopupEvent"`

	// 用户管理订阅通知回调
	SubscribeMsgChangeEvent struct {
		List []message.SubscribeMessageList `xml:"List"`
	} `xml:"SubscribeMsgChangeEvent"`

	// 用户发送订阅通知回调
	SubscribeMsgSentEvent struct {
		List []message.SubscribeMessageList `xml:"List"`
	} `xml:"SubscribeMsgSentEvent"`
}

func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

func (srv *Server) SetMessageHandler(messageHandler MessageHandler) {
	srv.messageHandler = messageHandler
}

// ServeHTTP 小程序消息处理中间件
// GET 验证消息的确来自微信服务器
// POST 处理客服消息
func (srv *Server) ServeHTTP(request *http.Request, writer http.ResponseWriter) {
	if request.Method == "GET" {
		srv.messageHandleValid(request, writer)
	} else if request.Method == "POST" {
		srv.messageHandle(request, writer)
	} else {
		srv.messageHandleNotSupport(request, writer)
	}
}

// messageHandleValid 消息校验
func (srv *Server) messageHandleValid(request *http.Request, writer http.ResponseWriter) {
	var (
		query     = request.URL.Query()
		nonce     = query.Get("nonce")
		echostr   = query.Get("echostr")
		signature = query.Get("signature")
		timestamp = query.Get("timestamp")
	)

	if signature == util.Signature(srv.Token, timestamp, nonce) {
		writer.Write([]byte(echostr))
	} else {
		writer.Write([]byte("invalid"))
	}
}

func (srv *Server) messageHandle(request *http.Request, writer http.ResponseWriter) {
	var (
		err          error
		contentType  = request.Header.Get("Content-Type")
		query        = request.URL.Query()
		nonce        = query.Get("nonce")
		timestamp    = query.Get("timestamp")
		encryptType  = query.Get("encrypt_type")
		msgSignature = query.Get("msg_signature")
	)

	for ok := true; ok; ok = !ok {
		var (
			encryptData          []byte
			encryptMsg           message.EncryptedMsg
			mixMessage           miniProgramMixMessage
			mixMessageReader     io.Reader
			subscribeMessageList []message.SubscribeMessageList
		)

		if encryptType == "aes" {
			// 解析密文消息
			// 校验消息是否合法
			// 解密密文消息
			if err = Decode(contentType, request.Body, &encryptMsg); err != nil {
				break
			} else if util.Signature(srv.Token, timestamp, nonce, encryptMsg.EncryptedMsg) != msgSignature {
				err = fmt.Errorf("invalid message")
				break
			} else if _, encryptData, err = util.DecryptMsg(srv.AppID, encryptMsg.EncryptedMsg, srv.EncodingAESKey); err != nil {
				break
			}

			mixMessageReader = bytes.NewBuffer(encryptData)
		} else {
			mixMessageReader = request.Body
		}

		// 解析到明文结构
		if err = Decode(contentType, mixMessageReader, &mixMessage); err != nil {
			break
		}

		// 处理订阅消息json和xml格式不同的情况
		if IsXML(contentType) {
			switch mixMessage.Event {
			case message.EventSubscribeSent:
				subscribeMessageList = mixMessage.SubscribeMsgSentEvent.List
			case message.EventSubscribePopup:
				subscribeMessageList = mixMessage.SubscribeMsgPopupEvent.List
			case message.EventSubscribeChange:
				subscribeMessageList = mixMessage.SubscribeMsgChangeEvent.List
			}
		} else {
			subscribeMessageList = mixMessage.List
		}

		// 处理消息
		srv.messageHandler(&message.MiniProgramMixMessage{
			CommonToken:  mixMessage.CommonToken,
			MsgID:        mixMessage.MsgID,
			Content:      mixMessage.Content,
			PicURL:       mixMessage.PicURL,
			MediaID:      mixMessage.MediaID,
			Title:        mixMessage.Title,
			AppID:        mixMessage.AppID,
			PagePath:     mixMessage.PagePath,
			ThumbURL:     mixMessage.ThumbURL,
			ThumbMediaID: mixMessage.ThumbMediaID,
			Event:        mixMessage.Event,
			SessionFrom:  mixMessage.SessionFrom,
			List:         subscribeMessageList,
		})
	}

	if err != nil {
		log.Printf("miniprogram.Message.Handle.Error: %s\n", err)
	}

	writer.Write([]byte("success"))
}

// messageHandleNotSupport 不支持的消息
func (srv *Server) messageHandleNotSupport(request *http.Request, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusMethodNotAllowed)
	writer.Write([]byte("Method Not Allowed"))
}
