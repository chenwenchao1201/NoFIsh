package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	"log"
	"strings"
	"time"
)

var fishesMap = map[string]int{}

// 白名单list,后续可以通过页面增加
var whiteList = []string{"微信读书", "马士兵", "知识星球", "小报童", "xzgedu"}

// 黑名单list,后续可以通过页面增加
var blackList = []string{"google", "知乎", "即刻"}

// 判断是否摸鱼的等待时间
var waitTime = 5

// 设定每日开始时间
var beginTime = 9

// 设定每日结束时间
var endTime = 18

var lastLearnTime = time.Now()

// takeARest 20分钟提醒休息一下
func takeARest() {
	tick := time.Tick(20 * time.Minute)
	for {
		<-tick
		if isInWorkTime() {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "休息一下",
				Content: "看电脑20分钟了，休息一下比较好",
			})
		}
	}
}

// fishCheck 摸鱼检查
func fishCheck() {
	tick := time.Tick(20 * time.Second)
	for {
		if isInWorkTime() {
			fishCheckTask()
		}
		<-tick
	}
}

// fishCheckTask 具体摸鱼检查
func fishCheckTask() {

	title := strings.ToLower(robotgo.GetTitle())
	// 黑名单内直接判断
	if isInBlackList(title) {
		log.Println("当前窗口标题：", title, "，疑似在摸鱼,最近摸鱼时间：", lastLearnTime)
		// 记录摸鱼时间,如果超过5分钟就弹窗
		if time.Now().Sub(lastLearnTime) > time.Duration(waitTime)*time.Minute {
			fmt.Println("摸鱼时间超过5分钟")
			today := time.Now().Format("2006-01-02")
			// 获取今日日期
			fishesMap[today]++
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "摸鱼警告",
				Content: "你已经摸鱼5分钟了,今日共摸鱼" + fmt.Sprintf("%d", fishesMap[today]) + "次",
			})
			myApp.FishCount = fishesMap[today]
			lastLearnTime = time.Now()
		}
	} else {
		// 清空最近摸鱼时间
		lastLearnTime = time.Now()
		log.Println("当前窗口标题：", title, " 非摸鱼，重新开始计时：", lastLearnTime)
	}
}

// isInWorkTime 工作时间检测
func isInWorkTime() bool {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	if hour < beginTime || hour > endTime {
		return false
	}
	if hour == beginTime && minute < 30 {
		return false
	}
	if hour == endTime && minute > 0 {
		return false
	}
	return true
}

// isInWhiteList 白名单检测
func isInWhiteList(title string) bool {
	for _, white := range whiteList {
		if strings.Contains(title, white) {
			return true
		}
	}
	return false
}

// isInBlackList 黑名单检测
func isInBlackList(title string) bool {
	for _, black := range blackList {
		if strings.Contains(title, black) {
			return true
		}
	}
	return false

}
