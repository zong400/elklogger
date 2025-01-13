package main

import (
	elklogger "elklogger/udplogger"
	"log/slog"
)

func main() {
	elk := elklogger.NewUdpLogger("logstash-wbyb:5111", "test-elk-log")
	slog.SetDefault(elk.Logger)
	slog.Info("testing elk log for go.")
	slog.Warn("testing elk warnning for go.")
	slog.Error("testing elk error log for go.")
}
