package fwatch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// Main -
func Main() {
	cfpath := ""

	(&cli.App{
		Name:  "fwatch",
		Usage: "watch file/directory modifying",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       GetConfigDir() + "/config.toml",
				Usage:       "config file path",
				Destination: &cfpath,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("config path is ", cfpath)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run fwatch",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Value:   false,
						Usage:   "show all logs",
					},
					&cli.BoolFlag{
						Name:    "daemon",
						Aliases: []string{"d"},
						Value:   false,
						Usage:   "launch daemon",
						Hidden:  true,
					},
				},
				Action: func(c *cli.Context) error {
					daemon := c.Bool("daemon")
					if daemon {
						pidFilePath := GetPidDir() + "/" + strconv.Itoa(os.Getpid())
						pidFile, err := os.Create(pidFilePath)
						if err != nil {
							log.Fatalf("cannot make pid file: %s", pidFilePath)
						}
						defer pidFile.Close()
						defer os.Remove(pidFilePath)

						Verbose = false
						IsService = true
						PidFile = pidFile
						pidFile.Write([]byte(cfpath + "\n"))

						log.SetOutput(pidFile)
					} else {
						Verbose = c.Bool("verbose")
						IsService = false
					}
					return Run(cfpath)
				},
			},
			{
				Name:  "start",
				Usage: "start fwatch service",
				Action: func(c *cli.Context) error {
					bin := os.Args[0]
					args := []string{"-c", cfpath, "run", "--daemon"}

					cmd := exec.Command(bin, args...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Start(); err != nil {
						log.Fatalf("fwatch: error: %v\n", err)
					}
					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "stop fwatch service",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					pids := GetPidList()
					if c.Bool("all") {
						for _, p := range pids {
							proc, err := os.FindProcess(p.Pid)
							if err == nil {
								proc.Signal(os.Interrupt)
								fmt.Printf("%d\t%s\n", p.Pid, p.Path)
							}
						}
					} else {
						pidstr := c.Args().First()
						pid, err := strconv.Atoi(pidstr)
						if err != nil {
							log.Fatalf("pid %d not found\n", pid)
						}
						for _, p := range pids {
							if p.Pid != pid {
								continue
							}
							proc, err := os.FindProcess(p.Pid)
							if err == nil {
								proc.Signal(os.Interrupt)
								fmt.Printf("%d\t%s\n", p.Pid, p.Path)
							}
						}
					}
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "show config list",
				Action: func(c *cli.Context) error {
					cf, err := Load(cfpath)
					if err != nil {
						return err
					}

					fmt.Println("No.\tpath\tscript\ttype")
					for i, t := range cf.Targets {
						fmt.Println(strconv.Itoa(i) + "\t" + t.Path + "\t\"" + t.Script + "\"\t" + strings.Join(t.Type, ","))
					}
					return nil
				},
			},
			{
				Name:  "register",
				Usage: "register new file watcher",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "file path that is watched",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "script",
						Aliases:  []string{"s"},
						Usage:    "execute script",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "file event type",
					},
				},
				Action: func(c *cli.Context) error {
					cf, err := Load(cfpath)
					if err != nil {
						return err
					}

					p := c.Path("path")
					s := c.String("script")
					t := c.StringSlice("type")

					if p == "" {
						log.Fatalf("path is required. you should set '-p path'\n")
					}
					if s == "" {
						log.Fatalf("script is required. you should set '-s \"shell script\"'\n")
					}

					if len(t) == 0 {
						t = append(t, "write")
					}

					err = cf.Register(p, s, t)

					print("%v\n", err)

					return err
				},
			},
			{
				Name:  "unregister",
				Usage: "unregister file watcher",
				Action: func(c *cli.Context) error {

					cf, err := Load(cfpath)
					if err != nil {
						return err
					}

					indexstr := c.Args().First()

					index, err := strconv.Atoi(indexstr)
					if err != nil {
						return err
					}

					err = cf.Unregister(index)

					return err
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "show fwatch service pid",
				Action: func(c *cli.Context) error {
					procs := GetPidList()
					fmt.Printf("pid\tpath\n")
					for _, proc := range procs {
						fmt.Printf("%d\t%s\n", proc.Pid, proc.Path)
					}
					return nil
				},
			},
		},
	}).Run(os.Args)
}
