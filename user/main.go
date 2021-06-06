package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"user/domain/repository"
	service2 "user/domain/service"
	"user/handler"
	pb "user/proto/user"
)

func main() {
	src := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("lastest"),
	)
	src.Init()
	db, err := gorm.Open("mysql", "root:123456@/micro?charset=utf8&parseTime=True&loc=Local")
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
