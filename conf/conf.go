package conf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

// WatchBuilders : the watch/build commands
var WatchBuilders = make(map[string][]string)

// Port : the port to run on
var Port int = 8085

// HTTPS : run in https mode
var HTTPS bool = false

// Root : the root directory to serve
var Root string = "."

var _dir, _ = os.Getwd()

// WorkingDir : the working directory
var WorkingDir string = _dir

// Init : get the config
func Init() error {
	// get json conf
	viper.SetConfigName("statik.config")
	viper.AddConfigPath(".")
	viper.SetDefault("watch", map[string][]string{})
	viper.SetDefault("port", Port)
	viper.SetDefault("https", HTTPS)
	viper.SetDefault("root", Root)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return errors.New("Can not find config file")
		}
		log.Fatal(err)
	}
	Port = viper.GetInt("port")
	HTTPS = viper.GetBool("https")
	Root = viper.GetString("root")
	// watchers
	wb := viper.Get("watch").(map[string]interface{})
	for w := range wb {
		b := wb[w].([]interface{})
		c := []string{}
		for _, cmd := range b {
			c = append(c, cmd.(string))
		}
		WatchBuilders[w] = c
	}
	return nil
}

// Create : create a config file
func Create() {
	w := make(map[string][]string)
	w["."] = append(w["."], "echo something has changed")
	data := map[string]interface{}{
		"watch": w,
	}
	jsonString, _ := json.MarshalIndent(data, "", "    ")
	ioutil.WriteFile("statik.config.json", jsonString, os.ModePerm)
}

func generateRandomKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	key := hex.EncodeToString(bytes)
	return key
}
