package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var (
	filePath string
	job      string
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	flag.StringVar(&filePath, "file", "/etc/crontab", "crontab file path")
	flag.StringVar(&job, "job", "", "cron job: * * * * * ls")
}

func main() {
	flag.Parse()

	c := myCron{
		Cron: cron.New(),
	}
	if job != "" {
		c.AddJob(job)
	} else {
		c.AddFromFile(filePath)
	}

	c.Run()
}

var lineExp = regexp.MustCompile(`^(\S+\s+\S+\s+\S+\s+\S+\s+\S+)\s+(.+)$`)

type myCron struct {
	*cron.Cron
}

func (c *myCron) AddJob(job string) {
	m := lineExp.FindStringSubmatch(job)
	c.AddFunc("0 "+m[1], func() {
		out, err := exec.Command("bash", "-c", m[2]).CombinedOutput()
		if err != nil {
			log.Printf("command: %s out: \r\n%s with error: %s", m[2], out, err.Error())
		} else {
			log.Printf("command: %s out: \r\n%s", m[2], out)
		}
	})
}

func (c *myCron) AddFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("open file %s error: %s", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		if lineExp.MatchString(line) == true {
			c.AddJob(line)
		}
	}
}

func (c *myCron) Run() {
	if len(c.Entries()) == 0 {
		log.Fatal("exited: no job found")
	}

	c.Cron.Run()
}
