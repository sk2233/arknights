/*
@author: sk
@date: 2023/2/26
*/
package main

type Param struct {
	EventType int
	// 不固定 参数
	Player    *Player
	Enemy     *Enemy
	Atk, Def  *DataWarp
	Invalid   bool
	HurtValue int64
	RealHurt  bool // 是否为真实伤害
	Doctor    *Player
}
