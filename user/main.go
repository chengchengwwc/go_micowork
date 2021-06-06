package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"user/common"
	"user/domain/repository"
	service2 "user/domain/service"
	"user/handler"
	pb "user/proto/user"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
		return
	}
	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	src := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("lastest"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(consulRegistry),
	)
	src.Init()
	mysqlInfo := common.GetMySqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	db.SingularTable(true)
	userDataService := service2.NewUserDataService(repository.NewUserRepository(db))
	err = pb.RegisterUserHandler(src.Server(), &handler.User{UserDataService: userDataService})
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := src.Run(); err != nil {
		fmt.Println(err)
	}

}
