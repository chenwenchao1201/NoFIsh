package main

import (
	"NoFish/repository"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

// getToolBar 获取工具栏
func (app *Config) getToolBar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addTaskDialog()
		}),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addHoldingDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.refreshSum()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			app.setupDialog()
		}),
	)
	return toolbar
}

type AppTask struct {
	name     *widget.Entry
	desc     *widget.Entry
	deadline *widget.Entry
	score    *widget.Entry
	taskType int64
	priority int64
}

func (app *Config) addTaskDialog() dialog.Dialog {
	// 任务名
	taskNameEntry := widget.NewEntry()
	// 任务描述
	taskDescEntry := widget.NewMultiLineEntry()
	// 任务截止时间
	taskDeadlineEntry := widget.NewEntry()
	taskDeadlineEntry.PlaceHolder = "YYYY-MM-DD"
	taskDeadlineEntry.Validator = dateValidator
	// 任务积分
	taskScoreEntry := widget.NewEntry()
	taskScoreEntry.Validator = isIntValidator
	// 任务类型
	taskTypeEntry := widget.NewSelect([]string{"长期", "短期"}, func(s string) {
	})
	// 任务优先级
	taskPriorityEntry := widget.NewSelect([]string{"高", "中", "低"}, func(s string) {
	})

	// 新建一个对话框
	addForm := dialog.NewForm(
		"新增任务",
		"添加",
		"取消",
		[]*widget.FormItem{
			{Text: "任务名", Widget: taskNameEntry},
			{Text: "描述", Widget: taskDescEntry},
			{Text: "截止时间", Widget: taskDeadlineEntry},
			{Text: "积分", Widget: taskScoreEntry},
			{Text: "类型", Widget: taskTypeEntry},
			{Text: "优先级", Widget: taskPriorityEntry},
		},
		func(valid bool) {
			if valid {
				// strconv 处理字符串
				name := taskNameEntry.Text
				desc := taskDescEntry.Text
				deadline, _ := time.Parse("2006-01-02", taskDeadlineEntry.Text)
				score, _ := strconv.Atoi(taskScoreEntry.Text)
				taskType := taskTypeEntry.Selected
				taskTypeInt := 0
				if taskType == "长期" {
					taskTypeInt = 1
				} else {
					taskTypeInt = 2
				}
				priority := taskPriorityEntry.Selected
				priorityInt := 0
				if priority == "高" {
					priorityInt = 1
				} else if priority == "中" {
					priorityInt = 2
				} else {
					priorityInt = 3
				}
				// 保存到数据库
				_, err := app.DB.InsertTask(repository.Task{
					ID:          0,
					Name:        name,
					Description: desc,
					DueDate:     deadline,
					Completed:   false,
					Points:      score,
					IsLongTerm:  taskTypeInt,
					Priority:    priorityInt,
				})
				if err != nil {
					dialog.ShowError(err, app.MainWindow)
					app.ErrorLog.Println(err)
					return
				}

				// 刷新列表
				app.refreshTaskList()
			}
		},
		app.MainWindow)

	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}

func (app *Config) setupDialog() dialog.Dialog {
	return nil
}

func isIntValidator(text string) error {
	_, err := strconv.Atoi(text)
	if err != nil {
		return err
	}
	return nil
}

func isFloatValidator(text string) error {
	_, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return err
	}
	return nil
}

func dateValidator(text string) error {
	if _, err := time.Parse("2006-01-02", text); err != nil {
		return err
	}

	return nil
}
