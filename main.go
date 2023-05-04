package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/user"
	"path"
)

type Config struct {
	Ip   string
	Port string
}

func getConfigName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(currentUser.HomeDir, ".open-proxy"), nil
}

func load() (map[string]*Config, error) {
	configFileName, err := getConfigName()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*Config, 0), nil
		}
		return nil, err
	}

	var configs = make(map[string]*Config, 0)
	if err := json.Unmarshal(file, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func save(configs map[string]*Config) error {
	configFileName, err := getConfigName()
	if err != nil {
		return err
	}

	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(configs)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

func printOutput(config *Config) {
	ip := config.Ip
	port := config.Port
	fmt.Printf("export https_proxy=http://%s:%s; export http_proxy=http://%s:%s; export all_proxy=socks5://%s:%s\n", ip, port, ip, port, ip, port)
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "load",
				Aliases: []string{"l"},
				Usage:   "Load a config",
				Action: func(cCtx *cli.Context) error {
					// If there is no args, configName will be blank
					// use blank string as the default config name
					configName := cCtx.Args().First()
					configs, err := load()
					if err != nil {
						return err
					}
					pConfig, ok := configs[configName]
					if !ok {
						return fmt.Errorf("config name do not exist")
					}
					printOutput(pConfig)
					return nil
				},
			},
			{
				Name:      "new",
				Usage:     "Add a config",
				UsageText: "new [command options] <ip> <port> (<ip> and <port> are optional)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "The name of the config",
					},
				},
				Action: func(cCtx *cli.Context) error {
					configs, err := load()
					if err != nil {
						return err
					}
					key := cCtx.String("name")

					ip := cCtx.Args().First()
					if ip == "" {
						ip = "127.0.0.1"
					}

					port := cCtx.Args().Get(1)

					if port == "" {
						port = "7890"
					}

					rlt := &Config{Ip: ip, Port: port}
					configs[key] = rlt
					err = save(configs)
					if err != nil {
						return err
					}

					printOutput(rlt)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
