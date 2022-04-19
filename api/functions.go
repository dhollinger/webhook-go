package api

import (
	"os/exec"
	"strings"

	"github.com/voxpupuli/webhook-go/config"
	"github.com/voxpupuli/webhook-go/lib/chatops"
	"github.com/voxpupuli/webhook-go/lib/orchestrators"
)

// Grab the ChatOps configuration from the config package and
// setup variable of type chatops.ChatOps and fill it with data
// from the configuration.
//
// This returns a reference to a chatops.ChatOps struct
func chatopsSetup() *chatops.ChatOps {
	conf := config.GetConfig().ChatOps
	c := chatops.ChatOps{
		Service:   conf.Service,
		Channel:   conf.Channel,
		User:      conf.User,
		AuthToken: conf.AuthToken,
		ServerURI: &conf.ServerUri,
	}

	return &c
}

// Determine if orchestration is enabled and either pass the cmd string slice to a
// the orchestrationExec function or localExec function.
//
// This returns an interface of the result of the execution and an error
func execute(cmd []string) (interface{}, error) {
	conf := config.GetConfig()
	var res interface{}
	var err error
	if conf.Orchestration.Enabled {
		res, err = orchestrationExec(cmd)
		if err != nil {
			return res, err
		}
	} else {
		res, err = localExec(cmd)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

// Executes r10k via orchestration tooling.
//
// Takes in an argument of is of type []string and passes the
// formatted command to the Deploy function in the ochestrators
// package.
//
// This function will return an interface{} containing the result
// and an error.
func orchestrationExec(cmd []string) (interface{}, error) {
	command := "\""
	for i := range cmd {
		command = command + cmd[i] + " "
	}
	command = strings.TrimSuffix(command, " ")
	command = command + "\""

	res, err := orchestrators.Deploy(command)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Executes r10k on the local system. This assumes that webhook-go, r10k,
// and the puppet server are all running on the same system.
//
// Takes in an argument, cmd, of type []string and returns a string result
// and error.
func localExec(cmd []string) (string, error) {
	args := cmd[1:]
	command := exec.Command(cmd[0], args...)

	res, err := command.CombinedOutput()
	if err != nil {
		return string(res), err
	}

	return string(res), nil
}
