package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"os/exec"
	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
    "log"
	"strings"
	"syscall"
	"reflect"
)

var(
    Token = "Bot <Bot_Token>"
	BotName = "<Cliant_ID_Bot>"
    stopBot = make(chan os.Signal, 1)
    vcsession *discordgo.VoiceConnection
    Prefixstring = "!"
    ChannelVoiceJoin = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
	buf []byte
	output string = "nil"
)

func main() {
	//Discordのセッションを作成
	discord, err := discordgo.New(Token)
	if err != nil {
        fmt.Println("Error logging in")
        fmt.Println(err)
	}else{
		fmt.Println("Hello")
	}

	discord.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
    // websocketを開いてlistening開始
    err = discord.Open()
    if err != nil {
        fmt.Println(err)
    }

	fmt.Println("Listening...")
	

signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
<-stopBot
err = discord.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
        log.Println("Error getting channel: ", err)
        return
	}
	
	fmt.Printf("%20s %20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Author.ID, m.Content)

    switch {
		case strings.HasPrefix(m.Content, fmt.Sprintf("%s", Prefixstring)): //Prefixがついていた時
			if m.Author.ID == "<Author_ID>" {
				outret, errorret := runCmdStr(m.Content[1:])
				if errorret != nil{
					sendMessage(s, c, "ERROR has occured")
				} else {
					sendMessage(s, c, outret)
				}
			}
    }
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
    _, err := s.ChannelMessageSend(c.ID, msg)

    log.Println(">>> " + msg)
    if err != nil {
		log.Println("Error sending message: ", reflect.TypeOf(err))
		sendMessage(s,c,"DiscordGo's ERROR (null message/over 4000 characters)")
    }
}

func runCmdStr(cmdstr string) (string,error) {
    // 文字列をコマンド、オプション単位でスライス化する
    c, err := shellwords.Parse(cmdstr)
    if err != nil {
        return output,err
	}
    switch len(c) {
    case 0:
		// 空の文字列が渡された場合
        return output,nil
    case 1:
        // コマンドのみを渡された場合
        buf, err = exec.Command(c[0]).Output()
    default:
        // コマンド+オプションを渡された場合
        // オプションは可変長でexec.Commandに渡す
        buf, err = exec.Command(c[0], c[1:]...).Output()
    }
    if err != nil {
        return output,err
	}
	output = string(buf)
	
    return output,nil
}