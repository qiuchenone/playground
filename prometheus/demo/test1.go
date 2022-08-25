package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	//初始化日志服务
	logger := log.New(os.Stdout, "[Memory]", log.Lshortfile|log.Ldate|log.Ltime)

	//初始一个http handler
	http.Handle("/metrics", promhttp.Handler())

	//初始化一个容器
	var (
		diskPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "memeory_percent",
			Help: "memeory use percent",
		},
			[]string{"percent"},
		)

		cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		})

		cpusTemprature = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "CPUs_Temperature",
				Help: "the temperature of CPUs.",
			},
			[]string{
				"cpuName", // Which cpu temperature?
			},
		)

		hdFailures = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "`hd_errors_total`",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		)
	)

	prometheus.MustRegister(diskPercent)
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(cpusTemprature)
	prometheus.MustRegister(hdFailures)

	// 启动web服务，监听1010端口
	go func() {
		logger.Println("ListenAndServe at:localhost:8898")
		//err := http.ListenAndServe("localhost:8898", nil)
		err := http.ListenAndServe(":8898", nil)
		if err != nil {
			logger.Fatal("ListenAndServe: ", err)
		}
	}()

	//收集内存使用的百分比
	for {
		// 内存
		logger.Println("start collect memory used percent!")
		v, err := mem.VirtualMemory()
		if err != nil {
			logger.Println("get memeory use percent error:%s", err)
		}
		usedPercent := v.UsedPercent
		logger.Println("get memeory use percent:", usedPercent)
		diskPercent.WithLabelValues("usedMemory").Set(usedPercent)

		// CPU
		cpuTemp.Set(65.3)

		// 一次性统计多个cpu的温度
		cpusTemprature.WithLabelValues("cpu1").Set(rand.Float64())
		cpusTemprature.WithLabelValues("cpu2").Set(rand.Float64())
		cpusTemprature.WithLabelValues("cpu3").Set(rand.Float64())

		// 硬盘
		hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

		time.Sleep(time.Second * 2)
	}
}
