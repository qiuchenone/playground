package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

const ConsulHost = "10.176.168.119" // consul ip
const ConsulPort = 8500             // consul端口

const HttpHost = "10.176.168.119" // http ip
const HttpPort = 8084             // 要健康检查的端口

func getConsulAddress() string {
	return fmt.Sprintf("%s:%d", ConsulHost, ConsulPort)
}

func getHttpAddressHealth() string {
	return fmt.Sprintf("http://%s:%d/health", HttpHost, HttpPort)
}

// Register 注册
func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = getConsulAddress()

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           getHttpAddressHealth(),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

// AllServices 获得所有注册的服务
func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = getConsulAddress()

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

// FilterSerivice 获得指定的服务
func FilterSerivice() {
	cfg := api.DefaultConfig()
	cfg.Address = getConsulAddress()

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func main() {
	_ = Register(HttpHost, HttpPort, "user-web", []string{"mxshop", "bobby"}, "user-web")
	//AllServices()
	//FilterSerivice()
	fmt.Println(fmt.Sprintf(`Service == "%s"`, "user-web"))
}
