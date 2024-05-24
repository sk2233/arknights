package main

import (
	"GameBase2/app"
	"GameBase2/config"
	"GameBase2/room"
	R "arknights/res"
	"embed"
)

var (
	//go:embed res
	files         embed.FS
	StackRoom     *room.StackUIRoom
	UIInfo        *uiInfo
	CardShow      *cardShow
	Cursor        *cursor
	Tip           *tip
	GameManager   *gameManager
	ClickManager  *clickManager
	GridManager   *gridManager
	PlayerManager *playerManager
	EnemyManager  *enemyManager
)

func main() {
	config.ViewSize = GameSize
	config.Debug = true
	config.ShowFps = true
	config.Files = &files // 先使用内部资源 ，不存在  再寻找外部资源文件
	InitFont()
	app.Run(NewMainApp(), 1280, 720)
}

type MainApp struct {
	*app.App
}

// Init 必须先传入实例  初始化使用该方法
func NewMainApp() *MainApp {
	res := &MainApp{}
	res.App = app.NewApp()
	GameManager = NewGameManager()
	ClickManager = NewClickManager()
	GridManager = NewGridManager()
	PlayerManager = NewPlayerManager()
	EnemyManager = NewEnemyManager()
	StackRoom = room.NewStackUIRoom(config.RoomFactory.LoadAndCreate(R.MAP.MAIN).(*room.Room))
	StackRoom.AddManager(GameManager)
	StackRoom.AddManager(ClickManager)
	StackRoom.AddManager(GridManager)
	StackRoom.AddManager(PlayerManager)
	StackRoom.AddManager(EnemyManager)
	res.PushRoom(StackRoom)
	return res
}
