package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/sambaiz/go-logging-sample/config"
	"github.com/sambaiz/go-logging-sample/log"
	"github.com/spf13/viper"

	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	done := make(chan bool, 1)

	gocron.Every(
		uint64(viper.GetInt64(config.IntervalSeconds)),
	).Seconds().Do(writeLog)
	cron := gocron.Start()

	go func() {
		beforeWait := func() {
			fmt.Println("stopping")
			cron <- true
		}
		afterWait := func() {
			if err := log.Shutdown(); err != nil {
				fmt.Println(err)
			}
		}
		done <- gracefulShutdown(
			viper.GetInt(config.ShutdownWaitSeconds),
			beforeWait,
			afterWait,
		)
	}()

	fmt.Println("started")
	<-cron
	<-done
}

func gracefulShutdown(waitSeconds int, beforeWait func(), afterWait func()) bool {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	beforeWait()

	wait, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(waitSeconds)*time.Second)
	<-wait.Done()
	cancel()

	afterWait()

	return true
}

func writeLog() {
	log.Logger.Info("value log",
		zap.Float64("value", rand.Float64()*100),
	)
}
