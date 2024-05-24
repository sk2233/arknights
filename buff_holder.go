/*
@author: sk
@date: 2023/3/4
*/
package main

type BuffHolder struct {
	datas map[int]*DataWarp
	buffs map[string]*Buff
	state int // 添加 各种状态  现阶段 仅有  限制 移动  限制攻击  暂时仅对敌人 有效
	dirty bool
}

func NewBuffHolder() *BuffHolder {
	return &BuffHolder{datas: make(map[int]*DataWarp), buffs: make(map[string]*Buff), dirty: false}
}

func (b *BuffHolder) IsState(state int) bool {
	return b.state&state == 0
}

func (b *BuffHolder) Refresh() {
	for _, data := range b.datas {
		data.Reset()
	}
	b.state = 0
	for _, buff := range b.buffs {
		if data, ok := b.datas[buff.Field]; ok {
			data.Add += buff.Add
			data.Mul += buff.Mul
		} else {
			b.state |= buff.State
		}
	}
}

func (b *BuffHolder) GetData(field int) *DataWarp {
	return b.datas[field]
}

func (b *BuffHolder) CloneData(field int) *DataWarp { // 事件中的 必须使用Clone的
	return b.GetData(field).Clone()
}

func (b *BuffHolder) Update() {
	for name, buff := range b.buffs {
		buff.Timer--
		if buff.Timer == 0 {
			b.RemoveBuff(name)
		}
	}
	if b.dirty {
		b.dirty = false
		b.Refresh()
	}
}

func (b *BuffHolder) AddBuff(buff *Buff) {
	b.buffs[buff.Name] = buff // 重复的 直接 覆盖
	b.dirty = true
}

func (b *BuffHolder) RemoveBuff(name string) {
	delete(b.buffs, name)
	b.dirty = true
}

func (b *BuffHolder) InitData(field int, data float64) {
	b.datas[field] = NewDataWarp(data)
}

type DataWarp struct {
	Data float64
	Add  float64
	Mul  float64
}

func NewDataWarp(data float64) *DataWarp {
	return &DataWarp{Data: data, Mul: 1}
}

func (d *DataWarp) Float() float64 {
	return d.Data*d.Mul + d.Add
}

func (d *DataWarp) Int() int64 {
	return int64(d.Float())
}

func (d *DataWarp) Reset() {
	d.Add = 0
	d.Mul = 1
}

func (d *DataWarp) Clone() *DataWarp {
	return &DataWarp{Data: d.Data, Add: d.Add, Mul: d.Mul}
}

type Buff struct {
	Name  string
	Timer int64
	Field int
	State int
	Add   float64
	Mul   float64
}
