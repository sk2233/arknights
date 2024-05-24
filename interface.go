/*
@author: sk
@date: 2023/2/25
*/
package main

import (
	"GameBase2/model"

	"github.com/hajimehoshi/ebiten/v2"
)

type IUpdateClick interface {
	UpdateClick(pos complex128)
}

type IPosDraw interface {
	PosDraw(pos complex128, screen *ebiten.Image)
}

type ISkill interface { // 基本技能  被动 之类的
	CanHandle(param *Param) bool
	Action(param *Param)
} // 可以实现  IName 主动触发的    IUpdate  需要冷却的

type IInitiativeSkill interface {
	ISkill
	IProgress
	model.IName
	model.IUpdate
}

type IRecover interface {
	Recover(recover int64)
}

type IProgress interface { // 需要冷却的技能的  冷却状态
	GetProgress() float64 //  0 ～ 1
}

type IAttack interface {
	Attack()
}
