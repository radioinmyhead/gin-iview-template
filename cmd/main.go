//go:generate go-bindata -o static.go -fs -prefix ../ui/dist/ ../ui/dist/...
package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/BurntSushi/toml"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"os"
)

// Conf the config info
type Conf struct {
	Main struct {
		DB string
	}
	Server struct {
		Port      string
		AppPrefix string
	}
}

// NewConf get default conf
func NewConf() *Conf {
	conf := &Conf{}
	conf.Main.DB = "10.10.10.10"
	conf.Server.Port = "0.0.0.0:8080"
	conf.Server.AppPrefix = "/ui/"
	return conf
}

var conf *Conf

func main() {
	conf = NewConf()
	defer CloseDB()
	app := cli.NewApp()
	app.Name = "ctl"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "debug"},
		&cli.StringFlag{Name: "conf"},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		conffile := c.String("conf")
		if conffile != "" {
			log.Debug(conffile)
			_, err := toml.DecodeFile(conffile, conf)
			if err != nil {
				return err
			}
		}
		if err := InitDB(conf.Main.DB); err != nil {
			return err
		}
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:  "server",
			Usage: "run as server mode",
			Action: func(c *cli.Context) error {
				r := gin.Default()
				r.StaticFS(conf.Server.AppPrefix, AssetFile())
				r.NoRoute(func(c *gin.Context) {
					c.Redirect(http.StatusMovedPermanently, conf.Server.AppPrefix)
				})

				r.POST("/api/auth/login", func(c *gin.Context) {
					var login map[string]string
					if err := c.ShouldBindJSON(&login); err != nil {
						c.JSON(400, gin.H{"err": err.Error()})
						return
					}
					log.Info(login["password"], login["username"])
					c.JSON(200, gin.H{"token": "super_admin"})
					return
				})
				r.GET("/api/auth/userinfo", func(c *gin.Context) {
					token := c.Query("token")
					log.Debugf("%#v\n", token)
					c.JSON(200, gin.H{
						"name":    "username",
						"user_id": "4291d7da9005377ec9aec4a71ea837f",
						"access":  []string{"super_admin", "admin"},
						"token":   token,
						"avatar":  "https://gw.alipayobjects.com/zos/rmsportal/jZUIxmJycoymBprLOUbT.png",
					})
					return
				})
				r.Run(conf.Server.Port)
				return nil
			},
		}, {
			Name:  "worker",
			Usage: "run as worker mode",
			Action: func(c *cli.Context) error {
				fmt.Println("run in worker")
				return nil
			},
		},
		{
			Name: "mainctl",
			Subcommands: []*cli.Command{
				{
					Name: "subctl",
					Action: func(c *cli.Context) (err error) {
						fmt.Println("here i am")
						return nil
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// InitDB open db with dbaddress
func InitDB(dbaddress string) error {
	return nil
}

// CloseDB close db if db is opened
func CloseDB() {
	check_db_close := true
	if check_db_close {
		return
	}
	return
}
