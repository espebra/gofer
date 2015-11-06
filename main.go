package main

import (
	"crypto/tls"
	"log"
	"strconv"
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

	i.AddCallback("PRIVMSG", func(e *irc.Event) {
		var message = e.Message()
		var nick = e.Nick
		var sender = e.Arguments[0]
		if sender == i.GetNick() {
			sender = "private message"
		}
		log.Print("Received [" + message + "] on [" +
			sender + "] from [" + nick + "] on IRC")
	})

	router := mux.NewRouter()
	http.Handle("/", httpInterceptor(router))
	router.HandleFunc("/", reqHandler(APIIndex)).Methods("HEAD", "GET")
	router.HandleFunc("/{type}/{target}/{action}", reqHandler(APIHandler)).Methods("POST")

	// Make sure we reconnect if disconnected. Not sure if this needs to
	// be a goroutine.
	go i.Loop()

	err = http.ListenAndServe(cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port), nil)
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

	//target := r.FormValue("target")
    	//if target == "" {
	//	l.Print("Unable to read target form value")
	//	http.Error(w, "Unable to read *target* from the form data", 400)
    	//	return
    	//}

	m := IRCPrivMsg {}
	m.Target = target
	m.Message = message
	m.Action = action

	// Send the message back through the channel
	ch <- m

	l.Print("Sent [" + message + "] to [" + target + "] via HTTP")
	http.Error(w, "Message [" + message + "] sent to [" + target + "]", 200)
	return
}

// Function that handles each of the incoming HTTP requests
func APIIndex(w http.ResponseWriter, r *http.Request, ch chan IRCPrivMsg) {
	http.Error(w, "POST to /{channel,user}/{$channel,$username}/{action,privmsg} with the form data {message} set", 200)
	return
}
