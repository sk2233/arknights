/*
@author: sk
@date: 2023/2/26
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Range struct {
	offsets [][]int
}

func (r *Range) DirDraw(screen *ebiten.Image, pos complex128, dir int) {
	DrawRange(screen, pos, dir, r.offsets, ColorRange)
}

func (r *Range) CollisionEnemy(pos complex128, dir int) *Enemy {
	return CollisionEnemy(pos, dir, r.offsets)
}

func (r *Range) CollisionEnemies(pos complex128, dir int) []*Enemy {
	return CollisionEnemies(pos, dir, r.offsets)
}

func (r *Range) AddRange(offsets [][]int) {
	r.offsets = append(r.offsets, offsets...)
}

func (r *Range) CollisionPlayers(pos complex128, dir int) []*Player {
	return CollisionPlayers(pos, dir, r.offsets)
}

func (r *Range) SetRange(offsets [][]int) {
	r.offsets = offsets
}

func NewRange(offsets [][]int) *Range {
	return &Range{offsets: offsets}
}
