package main

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 3 {
		configFilePath := os.Args[1]
		chains := os.Args[2]
		height, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		cfg, err := conf.ReadConfig(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		global.Config = cfg
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   constant.DefaultTimeFormat,
			DisableHTMLEscape: true,
		})
		if level, err := logrus.ParseLevel(cfg.Log.LogLevel); err == nil {
			logrus.SetLevel(level)
		}
		repository.InitMgo(cfg.Mongo, context.Background())
		cache.InitRedisClient(cfg.Redis)
		task.LoadTaskConf(cfg.Task)
		addTransferDataTask := new(task.AddTransferDataTask)
		addTransferDataTask.RunWithParam(chains, height)
	} else {
		fmt.Println("./exec  [config_filepath] [chains] [height]")
	}

}
