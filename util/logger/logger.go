package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var prefix string;


func CreateLog(prefix string) *CustomLog {
	const LOG_LEVEL = "DEBUG"

	logger := log.New()

	customFormatter := &log.TextFormatter{
		TimestampFormat : "2006-01-02 15:04:05",
		FullTimestamp:true,
		//DisableLevelTruncation:true,
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetFormatter(customFormatter)
	logger.SetOutput(os.Stdout)

	//Only log the warning severity or above.
	//logLevel := viper.GetString("log.level")
	if LOG_LEVEL != ""{
		switch LOG_LEVEL {
		case "PANIC" : logger.SetLevel(log.PanicLevel)
		case "FATAL" : logger.SetLevel(log.FatalLevel)
		case "ERROR" : logger.SetLevel(log.ErrorLevel)
		case "WARN" : logger.SetLevel(log.WarnLevel)
		case "INFO" : logger.SetLevel(log.InfoLevel)
		case "DEBUG" : logger.SetLevel(log.DebugLevel)
		case "TRACE" : logger.SetLevel(log.TraceLevel)
		}
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	customLog := CustomLog{
		Prefix: prefix,
		Log: logger,
		Level: LOG_LEVEL,
	}

	return &customLog
}

type CustomLog struct {
	Prefix 			string 				`json:"prefix"`
	Log 			*log.Logger			`json:"prefix"`
	Level 			string 				`json:"level"`
}

func (mylog CustomLog) Debug(args ...interface{}) {
	mylog.Log.Debug("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Info(args ...interface{}) {
	mylog.Log.Info("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Warn(args ...interface{}) {
	mylog.Log.Warn("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Error(args ...interface{}) {
	mylog.Log.Error("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Fatal(args ...interface{}) {
	mylog.Log.Fatal("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Println(args ...interface{}) {
	mylog.Log.Println("["+mylog.Prefix +"]", args)
}

func (mylog CustomLog) Print(args ...interface{}) {
	mylog.Log.Print("["+mylog.Prefix +"]", args)
}


