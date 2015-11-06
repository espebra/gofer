package config

import (
	"encoding/json"
	"os"
)

type HTTP struct {
	Host		string
	Port		int
}

type Channel struct {
        Ch		chan string
        Name            string
        Key             string
}

type Configuration struct {
        Nickname	string  
        Username	string  
        Server          string  
        Port            int     
        TLS             bool    
        TLSSkipVerify   bool    
        Debug		bool    
        Verbose		bool    
        Channels        []Channel
	HTTP		HTTP
}

func (c *Configuration) Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	}

	return nil
}
