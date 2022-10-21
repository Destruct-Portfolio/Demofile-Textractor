package main

import (
	"fmt"
	"os"
	"strconv"
	"encoding/json"
	"io/ioutil"


	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
)

// Player messages
type PlayerLogs struct {
	SteamID64 string
	Messages []string
}

// Run like this: go run print_events.go -demo /path/to/demo.dem
func main() {

	// Result payload
	var data = map[string][]string {}

	f, err := os.Open(ex.DemoPathFromArgs())
	checkError(err)

	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()


	// Register handler for chat messages to print them
	p.RegisterEventHandler(func(e events.ChatMessage) {
		fmt.Printf("Chat - %s says: %s\n", formatPlayer(e.Sender), e.Text)
		data[strconv.FormatUint(e.Sender.SteamID64,10)] = append(data[strconv.FormatUint(e.Sender.SteamID64,10)], e.Text ) 
	})

	// Parse to end
	err = p.ParseToEnd()
	checkError(err)

	// Write result to JSON
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("temp.logs.json", file, 0644)

}

func formatPlayer(p *common.Player) string {
	if p == nil {
		return "?"
	}

	switch p.Team {
	case common.TeamTerrorists:
		return "[T]" + p.Name + "#" + strconv.FormatUint(p.SteamID64, 10)
	case common.TeamCounterTerrorists:
		return "[CT]" + p.Name+ "#" + strconv.FormatUint(p.SteamID64, 10)
	}

	return p.Name+ "#" + strconv.FormatUint(p.SteamID64, 10)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
