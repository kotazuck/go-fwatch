package fwatch

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func print(format string, a ...interface{}) {
	if Verbose {
		if len(a) > 0 {
			fmt.Printf("fwatch: "+format+"\n", a...)
		} else {
			fmt.Printf("fwatch: " + format + "\n")
		}
	}
}

// GetConfigDir - fwatch config directory path
func GetConfigDir() string {
	// get home directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("home directory not found\n")
	}

	// get config directory path
	// $XDG_CONFIG_HOME (~/.config)
	configDir := homeDir + "/.config"
	if _, err = os.Stat(configDir); err != nil {
		if os.IsNotExist(err) && os.Mkdir(homeDir+"/.config", 0755) != nil {
			log.Fatalf("config directory ($XDG_CONFIG_HOME ~/.config) make failed\n")
		}
	}

	// get fwatch config directory path
	// $XDG_CONFIG_HOME/fwatch (~/.config/fwatch)
	appConfigDir := configDir + "/fwatch"
	if _, err := os.Stat(appConfigDir); err != nil {
		if os.IsNotExist(err) && os.Mkdir(appConfigDir, 0755) != nil {
			log.Fatalf("fwatch config directory ($XDG_CONFIG_HOME/fwatch ~/.config/fwatch) make failed\n")
		}
	}

	// get fwatch pid files directory path
	// $XDG_CONFIG_HOME/fwatch/pids (~/.config/fwatch/pids)
	pidsDir := appConfigDir + "/pids"
	if _, err := os.Stat(pidsDir); err != nil {
		if os.IsNotExist(err) && os.Mkdir(pidsDir, 0755) != nil {
			log.Fatalf("fwatch config directory ($XDG_CONFIG_HOME/fwatch ~/.config/fwatch) make failed\n")
		}
	}

	return appConfigDir
}

// GetPidDir - get fwatch temp dir
func GetPidDir() string {
	dir := GetConfigDir() + "/pids"
	if err := os.Mkdir(dir, 0755); err != nil {
		if os.IsNotExist(err) && os.MkdirAll(dir, 0755) != nil {
			log.Fatalf("pid directory make failed\n")
		} else if os.IsExist(err) {
			return dir
		}
	}
	return dir
}

// GetPidList -
func GetPidList() []struct {
	Path string
	Pid  int
} {
	dir := GetPidDir()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("pid directory not found\n")
	}

	pids := []struct {
		Path string
		Pid  int
	}{}

	for _, f := range files {
		pid, _ := strconv.Atoi(f.Name())
		if _, err := os.FindProcess(pid); err != nil || pid == 0 {
			continue
		}
		fp, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			continue
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		if scanner.Scan() {
			path := scanner.Text()
			pids = append(pids, struct {
				Path string
				Pid  int
			}{
				path,
				pid,
			})
		}
	}
	return pids
}
