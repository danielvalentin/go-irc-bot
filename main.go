package main

import (
	"net"
	"fmt"
	"bufio"
	"net/textproto"
	"strings"
	"errors"
)

type Bot struct {
	server string
	port string
	nick string
	user string
	channel string
	pass string
	pread, pwrite chan string
	conn net.Conn
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		panic("Couldnt connect")
	}
	bot.conn = conn
	print("Connected to server");
	return bot.conn, nil
}

func (bot *Bot) GetReply(query string) (answer string, err error) {
	switch(query) {
		case "tester":
			return "teeeest", nil
	}
	return "", errors.New("wtf")
}

func main() {
	
	bot := &Bot{
		server:"irc.freenode.net",
		port:"6667",
		nick:"Saktobot",
		channel:"#saktobottest",
		pass:"",
		conn:nil,
		user:"Dontknow",
	}
	conn, _ := bot.Connect()
	fmt.Fprintf(conn, "USER %s 8 * :%s\r\n", bot.nick, bot.nick)
	fmt.Fprintf(conn, "NICK %s\r\n", bot.nick)
	fmt.Fprintf(conn, "JOIN %s\r\n", bot.channel)
	defer conn.Close()
	
	reader := bufio.NewReader(conn)
	tr := textproto.NewReader(reader)
	for {
		line, err := tr.ReadLine()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", line)
		if strings.Contains(line, "PING ") {
			// PING
			pongdata := strings.Split(line, "PING ")
			fmt.Fprintf(bot.conn, "PONG %s\r\n", pongdata[1])
			print("PONGing "+pongdata[1]+"\n")
		} else if strings.Contains(line, "PRIVMSG "+bot.nick) {
			// PRIVATE CHAT MESSAGE
			endofname := strings.Index(line, "!")
			sender := line[1:endofname]
			parts := strings.Split(line, bot.nick+" :")
			msg := parts[1]
			answer, err := bot.GetReply(msg)
			if err == nil {
				fmt.Fprintf(bot.conn, "PRIVMSG "+sender+" :"+answer+"\r\n")
			}
		} else if strings.Contains(line, "PRIVMSG "+bot.channel) {
			// GENERAL CHAT MESSAGE
			parts := strings.Split(line, bot.channel+" :")
			msg := parts[1]
			answer, err := bot.GetReply(msg)
			print(err)
			if err == nil {
				fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :"+answer+"\r\n")
			}
		}
	}
	
}
