package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"time"
)

// makeUI 创建UI
func (app *Config) makeUI() {
	fishCount, finishCount, prizeCount := app.getSum()

	// 创建一个容器
	summary := container.NewGridWithColumns(3, fishCount, finishCount, prizeCount)
	app.Summary = summary
	// 创建工具栏,绑定到主窗口上
	toolBar := app.getToolBar()
	app.ToolBar = toolBar
	tasksTabContent := app.tasksTab()
	holdingsTab := app.holdingsTab()

	// 创建标签页
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("当前任务", theme.HomeIcon(), tasksTabContent),
		container.NewTabItemWithIcon("任务设置", theme.InfoIcon(), canvas.NewText("任务设置和修改", nil)),
		container.NewTabItemWithIcon("奖品区域", theme.InfoIcon(), holdingsTab),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// add container to window

	finalContent := container.NewVBox(summary, toolBar, tabs)

	app.MainWindow.SetContent(finalContent)

	// 定时刷新总览 1分钟1次
	go func() {
		for range time.Tick(time.Second * 60) {
			app.refreshSum()
		}
	}()
}

// getSum 获取总览
func (app *Config) getSum() (*canvas.Text, *canvas.Text, *canvas.Text) {
	var fishCount, finishCount, prizeCount *canvas.Text

	fishCount = canvas.NewText(fmt.Sprintf("今日摸鱼次数: %d ", myApp.FishCount), nil)
	finishCount = canvas.NewText(fmt.Sprintf("今日完成数: %d ", myApp.FinishCount), nil)
	prizeCount = canvas.NewText(fmt.Sprintf("当前积分数: %d ", myApp.PrizeCount), nil)
	fishCount.TextSize, finishCount.TextSize, prizeCount.TextSize = 18, 18, 18

	fishCount.Alignment = fyne.TextAlignLeading
	finishCount.Alignment = fyne.TextAlignCenter
	prizeCount.Alignment = fyne.TextAlignTrailing
	return fishCount, finishCount, prizeCount
}

// refreshSum  刷新总览
func (app *Config) refreshSum() {
	app.InfoLog.Println("刷新总览")
	// 重新获取总览 并刷新
	fishCount, finishCount, prizeCount := app.getSum()
	app.Summary.Objects = []fyne.CanvasObject{fishCount, finishCount, prizeCount}
	app.Summary.Refresh()
}

// refreshTaskList 刷新任务列表
func (app *Config) refreshTaskList() {
	app.InfoLog.Println("刷新任务列表")
	app.Tasks = app.getTaskSlice()
	app.TasksTable.Refresh()
}

func (app *Config) refreshHoldingsTable() {
	app.InfoLog.Println("刷新奖品列表")
	app.Holdings = app.getHoldingSlice()
	app.HoldingsTable.Refresh()
}
