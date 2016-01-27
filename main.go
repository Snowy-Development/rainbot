package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/RyanPrintup/nimbus"
	"github.com/wolfchase/rainbot/lib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rainbot \"config\"")
		os.Exit(1)
	}

	nimConfig, rainConfig, err := rainbot.GetConfigs(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bot := &rainbot.Bot{
		/* Client      */ nimbus.NewClient(rainConfig.Host,
							rainConfig.Nick, *nimConfig),
		/* Version     */ "Alpha 0.3.0 (Steeljack)",
		/* ModuleNames */ rainConfig.GoModules,
		/* Channels    */ make(map[string]*rainbot.Channel),
		/* Parser      */ rainbot.NewParser(rainConfig.CmdPrefix),
		/* Handler     */ rainbot.NewHandler(),
		/* Mutex       */ sync.Mutex{},
	}

	fmt.Print("Connecting... ")

	bot.Connect(func(e error) {
		if e != nil {
			fmt.Println(e)
			return
		}

		fmt.Println("Done")

		bot.LoadModules()

		bot.Handler.AddInternalCommand("m", &rainbot.Command{
			Help: "The RainBot Module Manager helps manage modules.",
			Fun: func (msg *nimbus.Message, args []string) {
				bot.Say(msg.Args[0], "Not exactly working yet")
			},
			PM: true,
			CM: true,
		})

		// Commands listener
		bot.AddListener(nimbus.PRIVMSG, func(msg *nimbus.Message) {
			if bot.Parser.IsCommand(msg.Trailing) {
				command, args := bot.Parser.ParseCommand(msg.Trailing)
				bot.Handler.Invoke(msg, rainbot.CommandName(command), args)
			}
		})

		bot.AddListener(nimbus.PRIVMSG, func(msg *nimbus.Message) {
			text := msg.Trailing
			if text == "Hello, "+bot.Client.Nick {
				bot.Client.Say(msg.Args[0], "Hello there!")
			}
		})

		// Add/Update channel when topic is received
		bot.AddListener(nimbus.RPL_TOPIC, func(msg *nimbus.Message) {
			bot.Mu.Lock()

			name  := msg.Args[1]
			topic := msg.Trailing

			bot.Channels[strings.ToLower(name)].Topic = topic

			bot.Mu.Unlock()
		})

		// Update users for channel
		bot.AddListener(nimbus.RPL_NAMEREPLY, func(msg *nimbus.Message) {
			bot.Mu.Lock()

			channel := bot.Channels[strings.ToLower(msg.Args[2])]
			users   := strings.Split(strings.Trim(msg.Trailing, " "), " ")

			for _, user := range users {
				var name, rank string

				if strings.ContainsAny(string(user[0]), "+ & ~ & @ & &") {
					name, rank = user[1:], string(user[0])
				} else {
					name, rank = user, ""
				}

				channel.Users[name] = rank
			}

			bot.Mu.Unlock()
		})

		// Update on user Join
		bot.AddListener(nimbus.JOIN, func (msg *nimbus.Message) {
			bot.Mu.Lock()
			defer bot.Mu.Unlock()

			who, _ := bot.Parser.ParsePrefix(msg.Prefix)
			where  := msg.Args[0][1:]

			if who == bot.Nick {
				channel := rainbot.NewChannel(where)
				bot.Channels[strings.ToLower(where)] = channel
				return
			}

			channel := bot.Channels[strings.ToLower(where)]
			channel.Users[who] = ""
		})

		// Update on user Kick
		bot.AddListener(nimbus.KICK, func (msg *nimbus.Message) {
			bot.Mu.Lock()

			who, _ := bot.Parser.ParsePrefix(msg.Prefix)
			where  := msg.Args[0]

			bot.RemoveUser(who, where)

			bot.Mu.Unlock()
		})

		// Update on user Kill
		bot.AddListener(nimbus.KILL, func (msg *nimbus.Message) {
			bot.Mu.Lock()

			// Implement getInfo(msg) function?
			who, _ := bot.Parser.ParsePrefix(msg.Prefix)
			where  := msg.Args[0][1:]

			bot.RemoveUser(who, where)

			bot.Mu.Unlock()
		})

		// Update on user part
		bot.AddListener(nimbus.PART, func (msg *nimbus.Message) {
			bot.Mu.Lock()

			who, _ := bot.Parser.ParsePrefix(msg.Prefix)
			where  := msg.Args[0][1:]

			bot.RemoveUser(who, where)

			bot.Mu.Unlock()
		})

		// Update on user quit
		bot.AddListener(nimbus.QUIT, func (msg *nimbus.Message) {
			bot.Mu.Lock()

			who, _ := bot.Parser.ParsePrefix(msg.Prefix)
			where  := msg.Args[0][1:]

			if who == bot.Nick {
				delete(bot.Channels, strings.ToLower(where))
				return
				// Wait a minute...
			}

			delete(bot.Channels[strings.ToLower(where)].Users, who)

			bot.Mu.Unlock()
		})

		// Update on nick change
		bot.AddListener(nimbus.NICK, func (msg *nimbus.Message) {
			bot.Mu.Lock()

			fmt.Println(msg)

			bot.Mu.Unlock()
		})

		bot.Listen()
		result := <- bot.Quit

		fmt.Println(result)
	})
}
