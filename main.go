package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/go-vgo/robotgo"
	"github.com/goki/freetype/truetype"
	"os"
	"strings"
	"time"
)

type App struct {
	output *widget.Label
}

var myApp App

var fishesMap = map[string]int{}

// 白名单list,后续可以通过页面增加
var whiteList = []string{"微信读书", "马士兵", "知识星球", "小报童"}

// 黑名单list,后续可以通过页面增加
var blackList = []string{"google", "知乎", "即刻"}

// 等待时间
var waitTime = 5

// 设定每日开始时间
var beginTime = 9

// 设定每日结束时间
var endTime = 18

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	// 每日开始结束时间初始化
	output, entry, button := myApp.makeUI()
	// Vbox就是从上到下堆叠组件
	w.SetContent(container.NewVBox(output, entry, button))
	// 设置大小，否则就会按照最小匹配来展示
	w.Resize(fyne.Size{Height: 500, Width: 400})
	// 检查是否摸鱼
	go fishCheck()
	// 提醒休息一下，不管是不是在工作
	go takeARest()
	w.ShowAndRun() // 等于w show + a run
	// 程序就会卡在这里，下面的代码不会执行,除非关闭窗口
	println("run exit")
}

// takeARest 20分钟提醒休息一下
func takeARest() {
	tick := time.Tick(20 * time.Minute)
	for {
		if isInWorkTime() {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "休息一下",
				Content: "看电脑20分钟了，休息一下比较好",
			})
		}
		<-tick
	}
}

func fishCheck() {
	tick := time.Tick(time.Second)
	for {
		if isInWorkTime() {
			fishCheckTask()
		}
		<-tick
	}
}

func fishCheckTask() {
	lastLearnTime := time.Now()
	title := strings.ToLower(robotgo.GetTitle())
	// 白名单内直接 ，然后判断是否在黑名单内
	if isInWhiteList(title) {
		return
	}
	// 黑名单内直接判断
	if isInBlackList(title) {
		// 记录摸鱼时间,如果超过5分钟就弹窗
		// todo 5分钟以后可以手动配置
		if time.Now().Sub(lastLearnTime) > time.Duration(waitTime)*time.Minute {
			fmt.Println("摸鱼时间超过5分钟")
			// 获取今日日期
			today := time.Now().Format("2006-01-02")
			fishesMap[today] = fishesMap[today] + 1
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "摸鱼警告",
				Content: "你已经摸鱼5分钟了,今日共摸鱼" + fmt.Sprintf("%d", fishesMap[today]) + "次",
			})
			lastLearnTime = time.Now()
		}
	} else {
		// 清空最近摸鱼时间
		lastLearnTime = time.Now()
	}
}

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

func isInWhiteList(title string) bool {
	for _, white := range whiteList {
		if strings.Contains(title, white) {
			return true
		}
	}
	return false
}

func isInBlackList(title string) bool {
	for _, black := range blackList {
		if strings.Contains(title, black) {
			return true
		}
	}
	return false

}

// 挂载在App类型上的方法
func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Button) {
	// 定义3个组件
	output := widget.NewLabel("Hello world!")
	entry := widget.NewEntry()
	// 添加按钮事件，把entry组件里的内容放到app的output中
	btn := widget.NewButton("确认", func() {
		app.output.SetText(entry.Text)
	})
	// 设置按钮优先级，不同优先级有不同配色，还会跟着主题变化颜色
	btn.Importance = widget.HighImportance
	// 关联2个output
	app.output = output
	return output, entry, btn
}

// 初始化中文字体文件
func init() {
	fontPath, err := findfont.Find("ShangShouJianSongXianXiTi-2.ttf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found 'arial.ttf' in '%s'\n", fontPath)

	// load the font with the freetype library
	// 原作者使用的ioutil.ReadFile已经弃用
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	os.Setenv("FYNE_FONT", fontPath)

}
