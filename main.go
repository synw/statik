package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/synw/statik/conf"
	"github.com/synw/statik/watcher"
	"github.com/synw/statik/ws"
)

var v = flag.Bool("v", false, "Verbosity of the output: -v to enable")
var https = flag.Bool("https", false, "Start the server with https: -https to enable")
var root = flag.String("root", ".", "Define the static root folder")
var port = flag.Int("port", 8085, "Define the port to run on")
var nr = flag.Bool("nr", false, "Disable the autoreload")
var nw = flag.Bool("nw", false, "Disable the watchers")
var nc = flag.Bool("nc", false, "Do not use the config file")
var spa = flag.String("spa", "", "Path to run with single page app history mode")
var genConf = flag.Bool("conf", false, "generate a config file")
var genSSLcerts = flag.Bool("certs", false, "Help with generating self-signed ssl certificates to be able to run the server in https mode")

func main() {
	flag.Parse()

	if *genConf {
		conf.Create()
		fmt.Println("File statik.config.json created")
		return
	}

	conf.Init()
	// check flags to take precedence over config
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "https" {
			conf.HTTPS = *https
		} else if f.Name == "root" {
			conf.Root = *root
		} else if f.Name == "port" {
			conf.Port = *port
		} else if f.Name == "spa" {
			conf.SpaPath = *spa
		}
	})

	if *genSSLcerts {
		_genSSLCerts()
		return
	}

	if !*nw {
		if *v {
			fmt.Println("Running changes watcher")
		}
		go watcher.Watch(*v, !*nr)
		go ws.RunWs()
		if *v {
			for k, v := range conf.WatchBuilders {
				fmt.Println(k, v)
			}
		}
	}

	runServer()
}

func runServer() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", conf.Root)
	if conf.SpaPath != "" {
		if *v {
			fmt.Println("Running in spa history mode for path", conf.SpaPath)
		}
		baseSpaPath := "/" + conf.SpaPath
		e.Static(baseSpaPath+"/*", conf.Root+"/index.html")
		e.Static(baseSpaPath+"/", conf.Root+"/index.html")
		e.Static(baseSpaPath, conf.Root+"/index.html")
	}

	if conf.HTTPS {
		if *v {
			fmt.Println("Running https server")
		}
		err := e.StartTLS(":"+strconv.Itoa(conf.Port), "./cert.pem", "./key.pem")
		if err != nil {
			if strings.Contains(err.Error(), ".pem") {
				fmt.Println("Can not run server in https mode: certificates are missing")
				_genSSLCerts()
				return
			}
			panic(err.Error())
		}
	}

	if *v {
		fmt.Println("Running http server")
	}
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Port)))
}

func _genSSLCerts() {
	fmt.Println("Generate local certificates with a command like:")
	fmt.Println("openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout key.pem -out cert.pem")
}
