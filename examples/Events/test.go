/*
 * @Author: YanHui Li
 * @Date: 2022-02-10 11:28:59
 * @LastEditTime: 2022-02-24 16:32:39
 * @LastEditors: YanHui Li
 * @Description:
 * @FilePath: \go-onvif\examples\Events\test.go
 *
 */
package main

import (
	"log"
	"time"

	"github.com/liyanhui1998/go-onvif"
	event "github.com/liyanhui1998/go-onvif/types/events"
)

func main() {
	/* 连接指定设备 */
	dev, err := onvif.NewDevice(onvif.DeviceParams{Ipddr: "10.1.1.200", Username: "admin", Password: "123qweasdZXC"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	/* 调用CreatePullPointSubscription方法获取SubscriptionReference地址 */
	subres := event.CreatePullPointSubscriptionResponse{}
	err = dev.CallMethodInterface(event.CreatePullPointSubscription{}, &subres, "")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Print(subres)
	/* 测试读取事件消息 */
	for i := 0; i < 10; i++ {
		/* 调用PullMessages方法获取事件消息 */
		pull := event.PullMessagesResponse{}
		err := dev.CallMethodInterface(event.PullMessages{Timeout: "PT10S", MessageLimit: 2}, &pull, string(subres.SubscriptionReference.Address))
		if err != nil {
			log.Fatalf(err.Error())
		}
		if pull.NotificationMessage.Topic.TopicKinds != "" {
			log.Printf("%v", pull.NotificationMessage)
		}
		time.Sleep(100 * time.Millisecond)
	}
	/* 调用Unsubscribe方法取消事件订阅 */
	err = dev.CallMethodInterface(event.Unsubscribe{}, event.UnsubscribeResponse{}, string(subres.SubscriptionReference.Address))
	if err != nil {
		log.Fatalf(err.Error())
	}
}
