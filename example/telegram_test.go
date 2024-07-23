package example

import (
	"WarnNotify/lib"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func LoadConfig() lib.Config {
	var conf lib.Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}

func Test_TelegramNewNotify(t *testing.T) {
	lib.NewNotify(LoadConfig())
	time.Sleep(time.Second * 10)
}

func Test_TelegramWarnMessage(t *testing.T) {
	if err := lib.NewNotify(LoadConfig()).WarnMessage("hello"); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
}

func Test_TelegramWatch(t *testing.T) {
	sg := sync.WaitGroup{}
	i := 0
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		sg.Add(1)
		defer sg.Done()
		for range time.Tick(time.Second * 3) {
			i += rand.Intn(10)
		}
		log.Println("go routine 1 end")
	}(ctx)
	cancel2, err := lib.NewNotify(LoadConfig()).Watch(func() string {
		return fmt.Sprintf("got i %d", i)
	}, time.Second*5)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
		cancel2()
	}()
	sg.Wait()
	time.Sleep(time.Second * 10)
}

func Test1(t *testing.T) {
	sg := sync.WaitGroup{}
	i := 0
	ctx, cancel := context.WithCancel(context.Background())
	sg.Add(1)
	go func(ctx context.Context) {
		defer sg.Done()
		for range time.Tick(time.Second) {
			log.Printf("echo %d", i)
			i++
		}
		log.Println("go routine 1 end")
	}(ctx)
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	sg.Wait()
}
