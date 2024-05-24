/*
@author: sk
@date: 2023/2/25
*/
package main

import "image/color"

type CardData struct { // 伤害 = 攻击 * 100/(100+防御)
	Name                 string
	Img                  string
	Place, ChangePlace   int   // 放置后 更改 后的区域
	Atk, Def             int64 // 攻击/防御
	BlockNum, CostNum    int
	Hp                   int64
	AttackTime, CoolTime int64    // 多少帧   攻击一次/冷却完毕
	Skills               []string // 挂载的技能
	Range                [][]int  // 普通攻击的检测区域
	Attack               string
	Career               int
}

type EnemyData struct { // 先做一些白板   详细的状态机制之后再做
	Name        string
	Img         string
	Hp          int64
	Atk, Def    int // 攻击/防御
	MoveSpeed   float64
	AttackTime  int64
	AttackRange float64 // 近战该值为0  无意义  通过 Attack 攻击方式使用
	CanBlock    bool
	HurtLife    int    //   进入基地后 减多少
	Attack      string // 攻击方式
}

type Info struct {
	Text  string
	Color color.Color
}
