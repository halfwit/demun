package service 

import (
	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
}

func NewService() (*Service, error) {
	srv, err := daemon.New("demun", "dmenu daemon", daemon.SystemDaemon)
	if err != nil {
		return nil, err
	}

	return &Service{srv}, nil
}

func (service *Service) Manage(command string) (string, error) {
	switch command {
	case "install":
		return service.Install()
	case "remove":
		return service.Remove()
	case "start":
		return service.Start()
	case "stop":
		return service.Stop()
	case "status":
		return service.Status()
	default:
		return "Unknown command for daemon", nil
	}
}
