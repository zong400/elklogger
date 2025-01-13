package udplogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

type UdpHandler struct {
	slog.Handler
	appname string
	podname string
}

type UdpLogger struct {
	*slog.Logger
}

func NewUdpHandler(w io.Writer, opts *slog.HandlerOptions, appname, podname string) *UdpHandler {
	return &UdpHandler{
		Handler: slog.NewJSONHandler(w, opts),
		appname: appname,
		podname: podname,
	}
}

// Enabled 当前日志级别是否开启
// func (h *UdpHandler) Enabled(ctx context.Context, level slog.Level) bool {
// 	return h.Handler.Enabled(ctx, level)
// }

// Handle 处理日志记录，仅在 Enabled() 返回 true 时才会被调用
func (h *UdpHandler) Handle(ctx context.Context, record slog.Record) error {
	record.Add("appname", h.appname)
	record.Add("podname", h.podname)
	// record.AddAttrs(slog.Int64("timestamp", record.Time.UnixMicro()))
	return h.Handler.Handle(ctx, record)
}

// WithAttrs 从现有的 handler 创建一个新的 handler，并将新增属性附加到新的 handler
// func (h *UdpHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
// 	return h.Handler.WithAttrs(attrs)
// }

// WithGroup 从现有的 handler 创建一个新的 handler，并将指定分组附加到新的 handler
// func (h *UdpHandler) WithGroup(name string) slog.Handler {
// 	return h.Handler.WithGroup(name)
// }

// 初始化udp发送日志
func NewUdpLogger(updaddr, appname string) *UdpLogger {
	podname := os.Getenv("HOSTNAME")
	var writer io.Writer
	udpWriter, err := NewUDPWriter(updaddr)
	if err != nil {
		slog.Error("udp连接出错，只输出到console " + err.Error())
		writer = os.Stdout
	} else {
		writer = io.MultiWriter(os.Stdout, udpWriter)
	}

	l := slog.New(NewUdpHandler(
		writer,
		&slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.MessageKey {
					return slog.Attr{Key: "message", Value: a.Value}
				} else if a.Key == slog.TimeKey {
					// return slog.Attr{Key: "@timestamp", Value: slog.Int64Value(a.Value.Time().UnixMicro())}
					return slog.Attr{Key: "@timestamp", Value: a.Value}
				}
				return a
			},
		},
		appname,
		podname,
	))

	return &UdpLogger{l}
}

func NewUDPWriter(addr string) (io.Writer, error) {
	// 解析 UDP 地址
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %v", err)
	}

	// 返回 UDPConn，它实现了 io.Writer 接口
	return conn, nil
}
