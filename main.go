package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bamchoh/pollydent"

	"gopkg.in/yaml.v2"

	youtube "google.golang.org/api/youtube/v3"
)

type AppConfig struct {
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	SofTalkDir string `yaml:"softalk_dir"`
}

func Load(filename string) (*AppConfig, error) {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

func setLog(prefix string) *log.Logger {
	basedir := filepath.Dir(os.Args[0])
	logPath := filepath.Join(basedir, prefix+".log")
	f, _ := os.Create(logPath)
	return log.New(f, prefix+":", 0)
}

func main() {
	logger := setLog("sofys")

	ac, err := Load("sofys-polly.yml")
	if err != nil {
		logger.Println("Load error")
		logger.Println(err)
	}

	var p TextToSpeacher
	if ac.SofTalkDir != "" {
		p = NewSofTalk(ac.SofTalkDir)
	} else {
		p = pollydent.NewPollydent(
			ac.AccessKey,
			ac.SecretKey,
			nil,
		)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		// fmt.Println(txt)
		dec := json.NewDecoder(strings.NewReader(txt))
		msg := new(youtube.LiveChatMessage)
		err := dec.Decode(msg)
		if err != nil {
			logger.Println(err)
			continue
		}
		reg := regexp.MustCompile("https?://.*")
		readContent := reg.ReplaceAllString(msg.Snippet.DisplayMessage, "")
		p.ReadAloud(readContent)
	}

	if err := scanner.Err(); err != nil {
		logger.Println("reading standard input:", err)
	}
}
