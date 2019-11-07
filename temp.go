package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/rehtt/gogpio"
    "github.com/kardianos/service"
    "log"
    "os"
    "strconv"
    "time"
)

var (
    max, min    int
    showTemp, h bool
)

func main() {

    flag.BoolVar(&h, "h", false, "help")
    flag.IntVar(&max, "m", 55, "开启风扇温度")
    flag.IntVar(&min, "n", 50, "关闭风扇温度")
    flag.BoolVar(&showTemp, "s", false, "查看当前温度")
    flag.Usage = help
    flag.Parse()
    if h {
        flag.Usage()
    } else {
        if !showTemp {
            //服务
            cfg := &service.Config{
                Name:        "temp",
                DisplayName: "RPI CPU Temp",
                Description: "RPI CPU Temp",
            }
            prg := &program{}
            s, err := service.New(prg, cfg)
            if err != nil {
                log.Fatal(err)
            }
            // logger 用于记录系统日志
            logger, err := s.Logger(nil)
            if err != nil {
                log.Fatal(err)
            }
            if len(os.Args) == 2 { //如果有命令则执行
                err = service.Control(s, os.Args[1])
                if err != nil {
                    log.Fatal(err)
                }
            } else { //否则说明是方法启动了
                err = s.Run()
                if err != nil {
                    logger.Error(err)
                }
            }
            if err != nil {
                logger.Error(err)
            }
        } else {
            fmt.Println("当前温度：", t(), "°")
        }
    }

}

func t() int {
    txt, _ := os.Open("/sys/class/thermal/thermal_zone0/temp")
    defer txt.Close()
    buf := bufio.NewReader(txt)
    t, _, _ := buf.ReadLine()
    temp, _ := strconv.Atoi(string(t))
    temp = temp / 1000
    return temp
}

func help() {
    fmt.Fprintln(os.Stderr, `使用：temp [选项] [服务操作]
选项：`)
    flag.PrintDefaults()
    fmt.Fprintln(os.Stderr,`
服务操作：
    install        安装服务（开启服务必做）
    start        运行服务
    stop        停止服务
    uninstall    卸载服务

示例：
    temp -s            显示温度
    temp install && temp start        安装服务并启动
    temp stop        停止服务
    temp -m 50    start    带参数启动服务
启动服务需要root权限`)
}

//服务
type program struct{}

func (p *program) Start(s service.Service) error {
    log.Println("开始服务")
    go p.run()
    return nil
}
func (p *program) Stop(s service.Service) error {
    log.Println("停止服务")
    return nil
}
func (p *program) run() {
    // 这里放置程序要执行的代码……
    pin, _ := gogpio.Open(14, gogpio.OUT)
    for ; ; {
        if t() >= max {
            pin.High()
        } else if t() <= min {
            pin.Low()
        }
        time.Sleep(30 * time.Second)
    }
}
