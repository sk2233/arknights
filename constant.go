/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/model"
	"GameBase2/utils"
	R "arknights/res"

	"golang.org/x/image/font"
)

const (
	GridSize = 64 + 64i
	GameSize = 1280 + 720i
)

var (
	GameSpeeds = []float64{1, 2}
	Transforms []func(x, y int) (int, int)
)

func init() {
	Transforms = make([]func(x, y int) (int, int), 4)
	Transforms[0] = func(x, y int) (int, int) {
		return x, y
	}
	Transforms[1] = func(x, y int) (int, int) {
		return -y, x
	}
	Transforms[2] = func(x, y int) (int, int) {
		return -x, -y
	}
	Transforms[3] = func(x, y int) (int, int) {
		return y, -x
	}
}

const (
	PlaceLand = 1 << iota
	PlaceHighland
	PlaceStart // 仅是用于特殊展示
	PlaceEnd
	PlaceNone
)

var (
	ColorTextMask      = &model.Color{A: 0.5} // p point
	ColorPlaceHighLand = utils.RGB(214, 205, 183)
	ColorPlaceLand     = utils.RGB(156, 147, 122)
	ColorButton        = utils.RGB(45, 45, 45)
	ColorHint          = &model.Color{G: 1, A: 0.5}
	ColorRange         = &model.Color{R: 1, G: 1, A: 0.5}
	ColorPlayerHp      = utils.RGB(85, 176, 228)
	ColorEnemyHp       = utils.RGB(255, 128, 0)
	ColorCardMask      = &model.Color{R: 1, A: 0.5} // p point
	ColorRecover       = &model.Color{G: 1, A: 0.2}
	ColorGaiZhiHua     = &model.Color{B: 1, A: 0.2}
)

var (
	Font18 font.Face
	Font32 font.Face
	Font48 font.Face
	Font64 font.Face
)

func InitFont() {
	Font18 = config.FontFactory.CreateFont(R.RAW.IPIX, 64, 18)
	Font32 = config.FontFactory.CreateFont(R.RAW.IPIX, 64, 32)
	Font48 = config.FontFactory.CreateFont(R.RAW.IPIX, 64, 48)
	Font64 = config.FontFactory.CreateFont(R.RAW.IPIX, 64, 64)
}

const (
	EventTypePlayerInit = iota
	EventTypePlayerHurt
	EventTypePlayerRetreat
	EventTypePlayerAttack
	EventTypePlayerRecover
	EventTypePlayerSkill
	EventTypePlayerInitiative // 主动的技能必须有名称   冷却进度
	EventTypeEnemyAttack
	EventTypeEnemyDeath
)

const (
	FieldState = iota // 修改状态 专用 正好当默认值
	FieldAttackSpeed
	FieldMoveSpeed
	FieldAtk
	FieldDef
	FieldMaxHp
)

const (
	StateMove = 1 << iota
	StateAttack
)

const (
	CareerReinstall = iota
	CareerGuards
	CareerSniper
	CareerMedical
	CareerMagian
	CareerSummon // 召唤物  退出战斗 不会返回卡组
)

const (
	EnemyTag = 1 << iota
)
