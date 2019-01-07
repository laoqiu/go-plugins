# go-plugins

A repository for [Micro Plugins](https://github.com/micro/go-plugins)

## Nats Streaming Usage

```
import (
    ...
    stan "github.com/laoqiu/go-plugins/broker/nats-streaming"
)
...
func main() {
	sbroker := stan.NewBroker()
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.example"),
		micro.Version("latest"),
		micro.Broker(sbroker),
	)

	// Initialise service
	service.Init()

	if err := sbroker.Init(); err != nil {
		log.Fatal(err)
	}
	if err := sbroker.Connect(); err != nil {
		log.Fatal(err)
	}
	sub, err := sbroker.Subscribe("foo", func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return p.Ack()
	}, stan.SetManualAckMode(), stan.DurableName("i-will-remember"), broker.Queue("bar"))
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	...
}
```

## Q&A
1. 使用命令行的--broker=stan，默认broker会被替换，但是在main中引用broker，并非stan broker，即使broker调用到了stan broker类的Connect方法
2. 不先调用sbroker.Connect()，执行sbroker.Subscribe将conn为nil的引用错误，调试发现原因是此时broker并未正常连接上，延时处理发现正常工作，我觉得应该是broker的设计问题
3. 暂时不知道如何在main中使用micro.RegisterSubscriber调用插件的正确方法
