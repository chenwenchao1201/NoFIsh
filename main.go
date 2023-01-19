package main

import (
	"NoFish/repository"
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/goki/freetype/truetype"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type App struct {
	output *widget.Label
}

type Config struct {
	App        fyne.App
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	MainWindow fyne.Window
	// 存放概况
	Summary *fyne.Container
	ToolBar *widget.Toolbar
	// 存放httpClient的字段
	HttpClient *http.Client
	// 存放摸鱼次数
	FishCount   int
	FinishCount int
	PrizeCount  int
	// 数据库
	DB repository.Repository
	// 任务相关
	Tasks         [][]interface{}
	TasksTable    *widget.Table
	Holdings      [][]interface{}
	HoldingsTable *widget.Table

	// 添加任务临时存放
	appTask *AppTask
}

var myApp Config

//func main() {
//
//	// 创建应用
//	app := app.NewWithID("com.earl")
//	myApp.App = app
//	// 窗口初始化
//	initApp(app)
//	// 检查是否摸鱼
//	go fishCheck()
//	// 提醒休息一下，不管是不是在工作
//	go takeARest()
//	// 启动
//	myApp.MainWindow.ShowAndRun()
//}

func initApp(app fyne.App) {
	// 赋予一个初始值
	myApp.HttpClient = &http.Client{}
	// 创建logger
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	myApp.MainWindow = app.NewWindow("摸鱼观察者")
	myApp.MainWindow.Resize(fyne.NewSize(800, 600))
	myApp.MainWindow.SetFixedSize(true) // 设置成自适应大小
	myApp.MainWindow.SetMaster()        // 设置成主窗口

	// 数据库初始化
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}
	myApp.setupDB(sqlDB)
	// ui初始化
	myApp.makeUI()
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

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("db in:", path)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) addHoldingDialog() dialog.Dialog {
	addAmountEntry := widget.NewEntry()
	purChaseDateEntry := widget.NewEntry()
	purchasePriceEntry := widget.NewEntry()
	//app.AddHoldingsPurChaseAmountEntry = addAmountEntry
	//app.AddHoldingsPurChaseDateEntry = purChaseDateEntry
	//app.AddHoldingsPurChasePriceEntry = purchasePriceEntry

	dateValidator := func(text string) error {
		if _, err := time.Parse("2006-01-02", text); err != nil {
			return err
		}

		return nil
	}
	purChaseDateEntry.Validator = dateValidator

	isIntValidator := func(text string) error {
		_, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		return nil
	}

	addAmountEntry.Validator = isIntValidator

	isFloatValidator := func(text string) error {
		_, err := strconv.ParseFloat(text, 32)
		if err != nil {
			return err
		}
		return nil
	}

	purchasePriceEntry.Validator = isFloatValidator

	purChaseDateEntry.PlaceHolder = "YYYY-MM-DD"

	addForm := dialog.NewForm(
		"Add gold",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Amount om toz", Widget: addAmountEntry},
			{Text: "PurChase price", Widget: purchasePriceEntry},
			{Text: "PurChase date", Widget: purChaseDateEntry},
		},
		func(valid bool) {
			if valid {
				amount, _ := strconv.Atoi(addAmountEntry.Text)
				purchasePrice, _ := strconv.ParseFloat(purchasePriceEntry.Text, 32)
				purchaseDate, _ := time.Parse("2006-01-02", purChaseDateEntry.Text)
				purchasePrice = purchasePrice * 100.0

				_, err := app.DB.InsertHolding(repository.Holdings{
					Amount:        amount,
					PurchaseDate:  purchaseDate,
					PurchasePrice: int(purchasePrice),
				})

				if err != nil {
					app.ErrorLog.Println(err)
					return
				}

				app.refreshHoldingsTable()

			}
		},
		app.MainWindow)

	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}
