package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/serdyanuk/microtask/config"
	"github.com/streadway/amqp"
)

const connectionWaitTimeout = time.Second * 60

func connect(cfg *config.Rabbitmq) (*amqp.Connection, error) {
	// waiting for rabbitmq service start
	if !wait.New(wait.WithProto("tcp"), wait.WithDeadline(connectionWaitTimeout)).Do([]string{cfg.Host + cfg.Addr}) {
		return nil, fmt.Errorf("rabbitmq connection timeout")
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s%s/", cfg.User, cfg.Password, cfg.Host, cfg.Addr))
	if err != nil {
		return nil, err
	}
	log.Println("Service connected to rabbitmq")

	return conn, nil
}
