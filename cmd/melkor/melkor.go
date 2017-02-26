package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/alde/melkor"
	"github.com/alde/melkor/config"
	"github.com/alde/melkor/crawlers"
	"github.com/alde/melkor/server"
	"github.com/alde/melkor/version"

	"github.com/sirupsen/logrus"
	"github.com/braintree/manners"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf(version.Version)
		os.Exit(0)
	}

	go catchInterrupt()

	cfg := config.Initialize()
	setupLogging(cfg)

	crawlers := initializeCrawlers(cfg)
	go doCrawl(crawlers, cfg)

	bind := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	logrus.WithFields(logrus.Fields{
		"version": version.Version,
		"address": cfg.Address,
		"port":    cfg.Port,
	}).Info("Launching Melkor")
	router := server.NewRouter(cfg, crawlers)
	if err := manners.ListenAndServe(bind, router); err != nil {
		logrus.WithError(err).Fatal("Unrecoverable error!")
	}
}

func setupLogging(cfg *config.Config) {
	if cfg.LogFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	if s != os.Interrupt && s != os.Kill {
		return
	}
	logrus.Info("Shutting down.")
	os.Exit(0)
}

func doCrawl(crawlers melkor.Crawlers, cfg *config.Config) {
	for {
		for _, crawler := range crawlers {
			if err := crawler.DoCrawl(); err != nil {
				logrus.WithError(err).
					WithField("resource", crawler.Resource()).
					Error("Error while crawling")
			}
		}
		time.Sleep(time.Duration(cfg.CrawlInterval) * time.Second)
	}
}

func initializeCrawlers(c *config.Config) melkor.Crawlers {
	ic := crawlers.NewInstancesCrawler(c)
	return melkor.Crawlers{
		ic.Resource(): ic,
	}
}
