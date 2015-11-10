package main

import (
	"crypto/tls"
	"log"
	"bytes"
	"os/exec"
	"strings"
	"path/filepath"
	"path"
	"strconv"
	"errors"
	"net/http"
	"flag"
	"os"

	"github.com/gorilla/mux"
	"github.com/thoj/go-ircevent"

	"github.com/espebra/gofer/app/config"
)

// Global configuration variables
var cfg = config.Configuration {}
var config_file string = "/etc/gofer/config.json"
var l = log.New(os.Stdout, "", log.LstdFlags)

type IRCPrivMsg struct {
	Target		string
	Message		string
	Action		string
}

// Global channels
var ch = make(chan IRCPrivMsg, 1024)

func init() {
	// Read path to configuration file
	flag.StringVar(&config_file, "config",
		config_file, "Gofer configuration file (JSON)")

	flag.Parse()
}

func main() {
	var err error

	err = cfg.Read(config_file)
	if err != nil {
		l.Fatal(err)
	}

	// Initialize
	i := irc.IRC(cfg.Nickname, cfg.Username)
	i.VerboseCallbackHandler = cfg.Verbose
	i.Debug = cfg.Debug
	i.UseTLS = true
	i.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.TLSSkipVerify}

	// Connect to the IRC server
	err = i.Connect(cfg.Server + ":" + strconv.Itoa(cfg.Port))
	if err != nil {
		l.Fatal(err)
	}

	// Join channels
	for c := range cfg.Channels {
		var channel = cfg.Channels[c].Name
		var key = cfg.Channels[c].Key
		i.AddCallback("001", func(e *irc.Event) {
			i.Join(channel + " " + key)
		})
	}

	// Set up a communications channel for sending Privmsg to users and
	// channels on IRC.
	i.AddCallback("002", func(e *irc.Event) {
		for n := range ch {
			if n.Action == "privmsg" {
				i.Privmsg(n.Target, n.Message)
			} else if n.Action == "action" {
				i.Action(n.Target, n.Message)
			}
		}
	})

	// Read messages said on IRC by other users
	i.AddCallback("PRIVMSG", func(e *irc.Event) {
		var message = e.Message()
		var nick = e.Nick
		var sender = e.Arguments[0]
		if sender == i.GetNick() {
			sender = "private message"
		}

		log.Print("Received [" + message + "] on [" +
			sender + "] from [" + nick + "] on IRC")

		// 2015/11/09 10:42:54 Received [bar foo zoo] on [#bar] from [someuser] on IRC
		// 2015/11/09 10:43:03 Received [meh meh] on [private message] from [someuser] on IRC

		// Check if the message begins with !, which is a command
		if string([]rune(message)[0]) == "!" && sender[0:1] == "#" {

			// Split the command and the arguments
			var command =  strings.Fields(message)[0]
			var args = strings.Fields(message)[1:]

			// Remove the ! from the command
			command = command[1:]

			// Execute the command
			out, err := Execute(command, args, sender[1:])
			if err != nil {
				l.Print("Unable to execute command: ", err)
			}
			if out != "" {
				l.Print("Result: " + out)
				i.Privmsg(sender, out)
			}
		}
	})

	go router()

	// Make sure we reconnect if disconnected. Not sure if this needs to
	// be a goroutine.
	i.Loop()

}

func router() {
	router := mux.NewRouter()
	http.Handle("/", httpInterceptor(router))
	router.HandleFunc("/", reqHandler(APIIndex)).Methods("HEAD", "GET")
	router.HandleFunc("/{type}/{target}/{action}", reqHandler(APIHandler)).Methods("POST")

	l.Print("Starting HTTP interface at " + cfg.HTTP.Host + ":" +
		strconv.Itoa(cfg.HTTP.Port))

	err := http.ListenAndServe(cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port), nil)
	if err != nil {
		l.Fatal(err.Error())
	}
}

func httpInterceptor(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})
}

func reqHandler(fn func (http.ResponseWriter, *http.Request, chan IRCPrivMsg)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ch)
	}
}

// Function that handles each of the incoming HTTP requests
func APIHandler(w http.ResponseWriter, r *http.Request, ch chan IRCPrivMsg) {
	vars := mux.Vars(r)
	target_type := vars["type"]
	target := vars["target"]
	action := vars["action"]

	if target_type == "channel" {
		target = "#" + target
	}

	if target_type != "channel" && target_type != "user" {
		l.Print("Invalid target type [" + target_type + "]")
		http.Error(w, "Invalid target type. " +
			"Valid types are [channel] and [user]", 400)
    		return
	}

	if action != "privmsg" && action != "action" {
		l.Print("Invalid target type [" + target_type + "]")
		http.Error(w, "Invalid action. " +
			"Valid actions are [privmsg] and [action]", 400)
    		return
	}

	message := r.FormValue("message")
    	if message == "" {
		l.Print("Unable to read message form value")
		http.Error(w, "Unable to read *message* from the form data",
			400)
    		return
    	}

	m := IRCPrivMsg {}
	m.Target = target
	m.Message = message
	m.Action = action

	// Send the message back through the channel
	ch <- m

	l.Print("Sent [" + message + "] to " + target_type + " [" + target +
		"] via HTTP")
	http.Error(w, "Message [" + message + "] sent to " + target_type +
		" [" + target + "]", 200)
	return
}

// Function that handles each of the incoming HTTP requests
func APIIndex(w http.ResponseWriter, r *http.Request, ch chan IRCPrivMsg) {
	http.Error(w, "POST to /{channel,user}/{$channel,$username}/{action,privmsg} with the form data {message} set", 200)
	return
}

func Execute(c string, args []string, channel string) (string, error) {
	var err error

	// Extract the last element of the path (filename) to make it slightly
	// more safe.
	c = path.Base(c)

	if c == "." {
		return "", errors.New("Empty command")
	}
	
	// Assemble the path to the script to execute
	// Add the channel as part of the path to enable different commands in different channels
	command := filepath.Join(cfg.CommandDirectory, "channel", channel, c)

        log.Print("Executing command: [", command, " ", strings.Join(args, " "), "]")
        cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
        return out.String(), err
}
