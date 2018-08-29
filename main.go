package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ggiamarchi/http-check/logger"
	"github.com/gin-gonic/gin"
	cli "github.com/jawher/mow.cli"
	yaml "gopkg.in/yaml.v2"
)

type check struct {
	Name    string  `json:"name" yaml:"name"`
	Command command `json:"command" yaml:"command"`
	Status  status  `json:"status" yaml:"status"`
}

func (c *check) String() string {
	return fmt.Sprintf("%+v", *c)
}

type command struct {
	Executable string        `json:"executable" yaml:"executable"`
	Args       []interface{} `json:"args" yaml:"args"`
}

func (c *command) String() string {
	return fmt.Sprintf("%+v", *c)
}

type status struct {
	Failure int `json:"failure" yaml:"failure"`
	Success int `json:"success" yaml:"success"`
}

func (s *status) String() string {
	return fmt.Sprintf("%+v", *s)
}

type appConfig struct {
	Checks []check
	Server struct {
		Port int
	}
}

func (c *appConfig) String() string {
	return fmt.Sprintf("%+v", *c)
}

func main() {

	app := cli.App("http-check", "HTTP Check exposes system commands as a single HTTP endpoint")

	app.Command("server", "Run HTTP Check server", func(cmd *cli.Cmd) {

		var configFile = cmd.StringOpt("c config", "/etc/http-check/http-check.yml", "HTTP Check YAML configuration file")

		cmd.Action = func() {
			logger.Init(false)
			logger.Info("Starting HTTP Check server...")

			appConfig := loadAppConfig(*configFile)

			s := &http.Server{
				Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
				Handler:      api(appConfig),
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
			}
			s.ListenAndServe()
		}
	})

	app.Run(os.Args)
}

func api(appConfig *appConfig) *gin.Engine {
	api := gin.New()
	api.Use(logger.APILogger(), gin.Recovery())

	v1 := api.Group("/v1")

	checks := make(map[string]check)
	for _, check := range appConfig.Checks {
		checks[check.Name] = check
	}

	v1.GET("/check/:name", func(c *gin.Context) {

		check := checks[c.Param("name")]

		stdout, stderr, err := execCommand(check.Command.Executable, check.Command.Args...)

		responseCode := check.Status.Success
		errorMsg := ""

		if err != nil {
			logger.Info("error  :: %s", err)
			responseCode = check.Status.Failure
			errorMsg = err.Error()
		}

		c.JSON(responseCode, gin.H{
			"stdout": stdout,
			"stderr": stderr,
			"error":  errorMsg,
		})
	})

	return api
}

func loadAppConfig(file string) *appConfig {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	c := appConfig{}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		panic(err)
	}

	return &c
}

func execCommand(command string, args ...interface{}) (string, string, error) {

	fmtCommand := fmt.Sprintf(command, args...)

	splitCommand := strings.Split(fmtCommand, " ")

	logger.Info("Executing command :: %s :: with args :: %v => %s", command, args, fmtCommand)

	cmdName := splitCommand[0]
	cmdArgs := splitCommand[1:len(splitCommand)]

	cmd := exec.Command(cmdName, cmdArgs...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}
