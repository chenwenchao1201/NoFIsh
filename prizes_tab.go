package main

import (
	"NoFish/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) prizesTab() *fyne.Container {
	app.Prizes = app.getPrizeSlice()

	app.PrizesTable = app.getPrizesTable()

	prizesContainer := container.NewBorder(nil, nil, nil, nil, container.NewAdaptiveGrid(1, app.PrizesTable))
	return prizesContainer

}

func (app *Config) getPrizesTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(app.Prizes), len(app.Prizes[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == (len(app.Prizes[0])-1) && i.Row != 0 {
				// last cell - put in buttton
				w := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {

					dialog.ShowConfirm("删除奖品", "确认删除？", func(delete bool) {
						if delete {
							id, _ := strconv.Atoi(app.Prizes[i.Row][0].(string))
							err := app.DB.DeletePrize(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						app.refreshPrizesTable()
					}, app.MainWindow)
				})
				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else if i.Col == (len(app.Prizes[0])-2) && i.Row != 0 {
				w := widget.NewButtonWithIcon("兑换", theme.ContentAddIcon(), func() {
					app.exchangePrize(i.Row)
				})
				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{widget.NewLabel(app.Prizes[i.Row][i.Col].(string))}
			}
		})

	colWidths := []float32{50, 200, 200, 200, 110}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}
	return t
}

// exchangePrize  兑换奖品
func (app *Config) exchangePrize(row int) {

}

// getPrizeSlice 从数据库中获取奖品信息
func (app *Config) getPrizeSlice() [][]interface{} {

	var slice [][]interface{}

	prizes, err := app.currentPrizes()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	slice = append(slice, []interface{}{"ID", "描述", "积分", "是否重复", "删除?"})

	for _, x := range prizes {

		var currentRow []interface{}
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, x.Description)
		currentRow = append(currentRow, strconv.Itoa(x.Points))
		switch int64(x.IsRepeat) {
		case 1:
			currentRow = append(currentRow, "是")
		case 0:
			currentRow = append(currentRow, "否")
		}
		currentRow = append(currentRow, widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {}))

		slice = append(slice, currentRow)

	}

	return slice

}

func (app *Config) currentPrizes() ([]repository.Prize, error) {
	prizes, err := app.DB.AllPrizes()

	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}

	return prizes, nil
}
