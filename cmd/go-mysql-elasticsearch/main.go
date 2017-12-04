package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql-elasticsearch/river"
	log "github.com/sirupsen/logrus"
)
//解析启动参数
var configFile = flag.String("config", "./etc/river.toml", "go-mysql-elasticsearch config file")
var my_addr = flag.String("my_addr", "", "MySQL addr")
var my_user = flag.String("my_user", "", "MySQL user")
var my_pass = flag.String("my_pass", "", "MySQL password")
var es_addr = flag.String("es_addr", "", "Elasticsearch addr")
var data_dir = flag.String("data_dir", "", "path for go-mysql-elasticsearch to save data")
var server_id = flag.Int("server_id", 0, "MySQL server id, as a pseudo slave")
var flavor = flag.String("flavor", "", "flavor: mysql or mariadb")
var execution = flag.String("exec", "", "mysqldump execution path")
var logLevel = flag.String("log_level", "info", "log level")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	level, _ := log.ParseLevel(*logLevel)
	log.SetLevel(level)
	//监听进程信号
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	cfg, err := river.NewConfigWithFile(*configFile)
	if err != nil {
		println(errors.ErrorStack(err))
		return
	}
	//启动参数覆盖toml配置文件
	if len(*my_addr) > 0 {
		cfg.MyAddr = *my_addr
	}

	if len(*my_user) > 0 {
		cfg.MyUser = *my_user
	}

	if len(*my_pass) > 0 {
		cfg.MyPassword = *my_pass
	}

	if *server_id > 0 {
		cfg.ServerID = uint32(*server_id)
	}

	if len(*es_addr) > 0 {
		cfg.ESAddr = *es_addr
	}

	if len(*data_dir) > 0 {
		cfg.DataDir = *data_dir
	}

	if len(*flavor) > 0 {
		cfg.Flavor = *flavor
	}

	if len(*execution) > 0 {
		cfg.DumpExec = *execution
	}

	r, err := river.NewRiver(cfg)
	if err != nil {
		println(errors.ErrorStack(err))
		return
	}

	//同步过程说明：
	//1. 检查本地文件系统中是否有master.info文件
	//1.1 如果有，读取指定位置的binlog
	//1.2 如果没有，执行mysqldump，从语句中获取master.info信息
	// dump命令：mysqldump --host=127.0.0.1 --port=3306 --user=root --password=123456 --master-data --single-transaction --skip-lock-tables --compact --skip-opt --quick --no-create-info --skip-extended-insert openapi1 api_order_info --default-character-set=utf8
	done := make(chan struct{}, 1)
	go func() {
		r.Run()
		done <- struct{}{}
	}()

	select {
	case n := <-sc:
		log.Infof("receive signal %v, closing", n)
	case <-r.Ctx().Done():
		log.Infof("context is done with %v, closing", r.Ctx().Err())
	}

	r.Close()
	<-done
}
