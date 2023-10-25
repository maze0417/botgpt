package botgpt

import (
	"botgpt/internal/config"
	"botgpt/internal/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func Run() {

	//migration.Migrate()

	env := config.GetConfig()

	r := NewRouter() //router
	port := fmt.Sprintf("0.0.0.0:%v", env.GetString("http.port"))

	fmt.Printf("Listening and serving HTTP on %s\n", port)
	setupLog()

	err := r.Run(port)
	if err != nil {
		panic(err)
	}

}

func setupLog() {
	env := config.GetConfig()
	isJson := env.GetBool("log.json")
	level := env.GetString("log.level")
	logrus.AddHook(&utils.ConsoleLogHook{})
	logrus.SetOutput(os.Stdout)
	logrus.AddHook(&utils.CustomFieldHook{})
	if isJson {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	switch {
	case level == "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case level == "info":
		logrus.SetLevel(logrus.InfoLevel)
	case level == "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case level == "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case level == "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case level == "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetReportCaller(true)
}
