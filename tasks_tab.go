package main

import (
	"NoFish/repository"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func (app *Config) tasksTab() *fyne.Container {
	app.Tasks = app.getTaskSlice()

	app.TasksTable = app.getTasksTable()

	tasksContainer := container.NewBorder(nil, nil, nil, nil, container.NewAdaptiveGrid(1, app.TasksTable))
	return tasksContainer

}

func (app *Config) getTasksTable() *widget.Table {

	table := widget.NewTable(
		func() (int, int) {
			// 设置长度和开始位置
			return len(app.Tasks), len(app.Tasks[0])
		},
		func() fyne.CanvasObject {
			// 设置表头
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			// 设置表格内容
			if i.Col == (len(app.Tasks[0])-1) && i.Row != 0 {
				// last cell - put in buttton
				w := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("删除任务", "确定删除当前任务？", func(delete bool) {
						if delete {
							id, _ := strconv.Atoi(app.Tasks[i.Row][0].(string))
							err := app.DB.DeleteHolding(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						app.refreshTaskList()
					}, app.MainWindow)
				})
				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{widget.NewLabel(app.Tasks[i.Row][i.Col].(string))}
			}
		})

	colWidths := []float32{30, 100, 300, 100, 100, 100, 30}
	for i := 0; i < len(colWidths); i++ {
		table.SetColumnWidth(i, colWidths[i])
	}
	return table
}

func (app *Config) getTaskSlice() [][]interface{} {
	var slice [][]interface{}

	tasks, err := app.currentTasks()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	slice = append(slice, []interface{}{"ID", "名字", "描述", "截止日期", "优先级", "积分", "删除"})

	for _, x := range tasks {
		var currentRow []interface{}
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, x.Name)
		currentRow = append(currentRow, x.Description)
		currentRow = append(currentRow, x.DueDate.Format("2006-01-02"))
		currentRow = append(currentRow, strconv.FormatInt(int64(x.Priority), 10))
		currentRow = append(currentRow, strconv.FormatInt(int64(x.Points), 10))
		currentRow = append(currentRow, widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {}))
		slice = append(slice, currentRow)
	}
	return slice
}

func (app *Config) currentTasks() ([]repository.Task, error) {
	tasks, err := app.DB.AllTasks()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return tasks, nil
}
