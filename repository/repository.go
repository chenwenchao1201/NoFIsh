package repository

import (
	"errors"
	"time"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

// Repository is the interface which must be satisfied in order to
// connect to a database
type Repository interface {
	Migrate() error
	// Holdings
	InsertHolding(h Holdings) (*Holdings, error)
	AllHoldings() ([]Holdings, error)
	GetHoldingByID(id int) (*Holdings, error)
	UpdateHolding(id int64, updated Holdings) error
	DeleteHolding(id int64) error
	// tasks
	InsertTask(t Task) (*Task, error)
	AllTasks() ([]Task, error)
	//GetTaskByID(id int) (*Task, error)
	//UpdateTask(id int64, updated Task) error
	//DeleteTask(id int64) error
	//// prizes
	//InsertPrize(p Prize) (*Prize, error)
	//AllPrizes() ([]Prize, error)
	//GetPrizeByID(id int) (*Prize, error)
	//UpdatePrize(id int64, updated Prize) error
	//DeletePrize(id int64) error
}

// Holdings is the type for the user's gold holdings
type Holdings struct {
	ID            int64     `json:"id"`
	Amount        int       `json:"amount"`
	PurchaseDate  time.Time `json:"purchase_date"`
	PurchasePrice int       `json:"purchase_price"`
}

type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	// 对应积分
	Points int `json:"points"`
	// 短期 or 长期
	IsLongTerm int `json:"is_long_term"`
	// 优先级
	Priority int `json:"priority"`
}

type Prize struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	// 对应积分
	Points int `json:"points"`
	// 是否重复兑换
	IsRepeat int `json:"is_repeat"`
	// 兑换时间
	ExchangeTime time.Time `json:"exchange_time"`
}

type summary struct {
	ID          int64  `json:"id"`
	FishCount   int64  `json:"fish_count"`
	FinishCount int64  `json:"finish_count"`
	PrizeCount  int64  `json:"prize_count"`
	day         string `json:"day"`
}
