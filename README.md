# Logger
实现了slog的logger方法，把日志通过udp发送到接收端，自动获取Hostname加入日志，实例化时候传入appname代表一组相同应用

# 接收端
需要upd接口，已测试filebeat udp，通过filebeat接入elk日志系统

# 用法

```
import elklogger "github.com/zong400/elklogger/udplogger"

func main() {
    // appname是一组app的名字
    elk := elklogger.NewUdpLogger("logstash-wbyb:5111", "appname", slog.LevelInfo)
    slog.SetDefault(elk.Logger)
    slog.Info("test upd logger")
}
```