package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

// Task ...
type Task struct {
	Content  string
	Complete bool
}

func main() {

	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "add, list, and complete tasks"
	app.Commands = []*cli.Command{
		{
			Name:  "add",
			Usage: "add a task",
			Action: func(c *cli.Context) error {
				task := Task{Content: c.Args().First(), Complete: false}
				AddTask(task, "tasks.json")
				return nil
			},
		},
		{
			Name:  "complete",
			Usage: "complete a task",
			Action: func(c *cli.Context) error {
				idx, err := strconv.Atoi(c.Args().First())
				if err != nil {
					panic(err)
				}
				CompleteTask(idx)
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "print all uncompleted tasks in list",
			Action: func(c *cli.Context) error {
				showAllList, err := strconv.ParseBool(c.Args().First())
				if err != nil {
					ListTasks(false)
				}
				ListTasks(showAllList)
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// AddTask ...
func AddTask(task Task, filename string) {
	j, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	j = append(j, "\n"...)
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if _, err = f.Write(j); err != nil {
		panic(err)
	}
}

// ListTasks ...
func ListTasks(showAll bool) {
	file := openTaskFile()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		j := scanner.Text()
		t := Task{}
		err := json.Unmarshal([]byte(j), &t)
		if err != nil {
			panic(err)
		}
		fmt.Println("№\ttodo\t\tfinished")
		if showAll {
			fmt.Printf("[%d]\t%s\t%v\n", i, t.Content, t.Complete)
		} else if !t.Complete {
			fmt.Printf("[%d]\t%s\t\n", i, t.Content)
		}
		i++
	}
	if i == 1 {
		fmt.Println("Everything is done!")
	}
}

func openTaskFile() *os.File {
	if _, err := os.Stat("tasks.json"); os.IsNotExist(err) {
		log.Fatal("tasks file does not exist")
		return nil
	}
	file, err := os.Open("tasks.json")
	if err != nil {
		panic(err)
	}
	return file
}

// CompleteTask ...
func CompleteTask(idx int) {
	file := openTaskFile()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		j := scanner.Text()
		t := Task{}
		err := json.Unmarshal([]byte(j), &t)
		if err != nil {
			panic(err)
		}
		if !t.Complete {
			if idx == i {
				t.Complete = true
			}
			i++
		}
		AddTask(t, ".tempfile")
	}
	os.Rename(".tempfile", "tasks.json")
	os.Remove(".tempfile")
}
