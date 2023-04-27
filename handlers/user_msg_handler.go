package handlers

import (
	"deouy/wechatbot/config"
	"deouy/wechatbot/gpt"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

var isUserStart = false
var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {

	if !config.LoadConfig().OpenFriendMode {
		return nil
	}
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	if UserService.ClearUserSessionContext(sender.ID(), msg.Content) {
		_, err = msg.ReplyText("上下文已经清空了，你可以问下一个问题啦。")
		if err != nil {
			log.Printf("response user error: %v \n", err)
		}
		return nil
	}

	// 获取上下文，向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")

	if requestText == "start bot" {
		isUserStart = true
		_, err = msg.ReplyText("启动成功")
		if err != nil {
			log.Printf("response group error: %v \n", err)
		}
		return err
	}
	if requestText == "shutdown bot" {
		isUserStart = false
		_, err = msg.ReplyText("关闭成功")
		if err != nil {
			log.Printf("response group error: %v \n", err)
		}
		return err
	}
	if !isUserStart {
		log.Printf("没启动")
		return nil
	}
	messages := UserService.GetUserSessionContext(sender.ID())
	messages = append(messages, gpt.CreateMessage(requestText))
	reply, err := gpt.Completions(messages)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 设置上下文，回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	UserService.SetUserSessionContext(sender.ID(), requestText)
	UserService.SetUserSessionContext(sender.ID(), reply)
	reply = "本消息由 chatGPT Bot回复：\n" + reply
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
