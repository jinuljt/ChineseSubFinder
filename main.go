package main

import (
	"github.com/allanpk716/ChineseSubFinder/common"
	"github.com/allanpk716/ChineseSubFinder/model"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	var err error
	log = model.GetLogger()
	configViper, err = model.InitConfigure()
	if err != nil {
		log.Errorln("InitConfigure", err)
		return 
	}
	config, err = model.ReadConfig(configViper)
	if err != nil {
		log.Errorln("ReadConfig", err)
		return 
	}
}

func main() {
	if log == nil {
		panic("log init error")
	}
	if configViper == nil {
		panic("init viper error")
	}
	if config == nil {
		panic("read config error")
	}
	httpProxy := config.HttpProxy
	if config.UseProxy == false {
		httpProxy = ""
	}
	// 判断文件夹是否存在
	if model.IsDir(config.MovieFolder) == false {
		log.Errorln("MovieFolder not found")
		return
	}
	log.Infoln("MovieFolder:", config.MovieFolder)

	// 下载实例
	downloader := NewDownloader(common.ReqParam{
		HttpProxy: httpProxy,
		DebugMode: config.DebugMode,
		SaveMultiSub: config.SaveMultiSub,
		Threads: config.Threads,
	})
	//任务还没执行完，下一次执行时间到来，下一次执行就跳过不执行
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	// 定时器
	entryID, err := c.AddFunc("@every " + config.EveryTime, func() {
		defer func() {
			log.Infoln("Download One End...")
		}()

		log.Infoln("Download One Started...")
		log.Infoln("Download Movie Sub Started...")
		// 开始下载
		err := downloader.DownloadSub4Movie(config.MovieFolder)
		if err != nil {
			log.Errorln("DownloadSub4Movie", err)
			return
		}
		log.Infoln("Download Movie Sub End...")

		log.Infoln("Download Series Sub Started...")
		err = downloader.DownloadSub4Series(config.SeriesFolder)
		if err != nil {
			log.Errorln("DownloadSub4Series", err)
			return
		}
		log.Infoln("Download Series Sub End...")
	})
	if err != nil {
		log.Errorln("cron entryID:", entryID, "Error:", err)
		return
	}

	log.Infoln("First Time Download Start")
	// 立即触发第一次的更新
	// 开始下载
	log.Infoln("Download Movie Sub Started...")
	err = downloader.DownloadSub4Movie(config.MovieFolder)
	if err != nil {
		log.Errorln("DownloadSub4Movie", err)
		return
	}
	log.Infoln("Download Movie Sub End...")

	log.Infoln("Download Series Sub Started...")
	err = downloader.DownloadSub4Series(config.SeriesFolder)
	if err != nil {
		log.Errorln("DownloadSub4Series", err)
		return
	}
	log.Infoln("Download Series Sub End...")

	log.Infoln("First Time Download End")

	c.Start()

	// 阻塞
	select {}
}

var(
	log         *logrus.Logger
	configViper *viper.Viper
	config      *common.Config
)
