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
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

func GetMainSprite(path string) model.ISprite {
	return config.SpriteFactory.CreateStaticSprite(R.SPRITE.MAIN, path)
}

func GetMainImg(path string) *ebiten.Image {
	return config.SpritesLoader.LoadStaticSprite(R.SPRITE.MAIN, path).Image
}

func ToIndex(pos complex128) (int, int) {
	return int(real(pos) / 64), int((imag(pos) - 32) / 64)
}

func ToPos(x, y int) complex128 {
	return utils.Int2Vector(x*64, y*64+32)
}

func GetDir(offset complex128) int {
	angle := utils.VectorAngle(offset) + math.Pi/4
	if angle < 0 {
		angle += math.Pi * 2
	}
	return int(angle / (math.Pi / 2))
}

func ParseHurt(param *Param) int64 {
	if param.RealHurt {
		return int64(param.Atk.Float())
	}
	return int64(math.Max(param.Atk.Float()-param.Def.Float(), param.Atk.Float()*0.05))
	//return int(math.Max(Atk.Float()*100/(def.Float()+100), Atk.Float()*0.05))
}

func ToFrame(second float64) int64 {
	return int64(second * 60)
}

func ToMoveSpeed(speed float64) float64 {
	return 64 * speed / 60
}

func ToAttackRange(attackRange float64) float64 {
	return 64 * attackRange
}

func GetMinPlayer(pos complex128, maxLen float64) *Player {
	l2 := maxLen * maxLen
	players := PlayerManager.GetPlayers(func(player *Player) bool {
		return utils.VectorLen2(player.Grid.Pos+GridSize/2-pos) < l2
	})
	if len(players) <= 0 {
		return nil
	}
	player := players[0]
	minL2 := utils.VectorLen2(player.Grid.Pos + GridSize/2 - pos)
	for i := 1; i < len(players); i++ {
		if l2 = utils.VectorLen2(players[i].Grid.Pos + GridSize/2 - pos); l2 < minL2 {
			minL2 = l2
			player = players[i]
		}
	}
	return player
}

func DrawRange(screen *ebiten.Image, pos complex128, dir int, offsets [][]int, clr *model.Color) {
	for i := 0; i < len(offsets); i++ {
		x, y := Transforms[dir](offsets[i][0], offsets[i][1])
		utils.FillRect(screen, pos+utils.Int2Vector(x*64, y*64), GridSize, clr)
	}
}

var (
	rectMask = model.NewRect(GridSize - 2 - 2i) // 使用偏移 解决沾边碰撞问题
)

func CollisionEnemy(pos complex128, dir int, offsets [][]int) *Enemy {
	for i := 0; i < len(offsets); i++ {
		x, y := Transforms[dir](offsets[i][0], offsets[i][1])
		rectMask.Pos = pos + utils.Int2Vector(x*64, y*64) + 1 + 1i
		if temp := utils.CollisionRect(R.LAYER.ENEMY, EnemyTag, rectMask); temp != nil {
			return temp.(*Enemy)
		}
	}
	return nil
}

func CollisionEnemies(pos complex128, dir int, offsets [][]int) []*Enemy {
	res := make([]*Enemy, 0)
	for i := 0; i < len(offsets); i++ {
		x, y := Transforms[dir](offsets[i][0], offsets[i][1])
		rectMask.Pos = pos + utils.Int2Vector(x*64, y*64) + 1 + 1i
		temp := utils.CollisionRectList(R.LAYER.ENEMY, EnemyTag, rectMask)
		for j := 0; j < len(temp); j++ {
			res = append(res, temp[j].(*Enemy))
		}
	}
	return res
}

func GetMinLastEnemy(enemies []*Enemy) *Enemy {
	if len(enemies) <= 0 {
		return nil
	}
	enemy := enemies[0]
	minLast := enemy.GetLast()
	for i := 1; i < len(enemies); i++ {
		last := enemies[i].GetLast()
		if last < minLast {
			minLast = last
			enemy = enemies[i]
		}
	}
	return enemy
}

func CollisionPlayers(pos complex128, dir int, offsets [][]int) []*Player {
	res := make([]*Player, 0)
	for i := 0; i < len(offsets); i++ {
		x, y := Transforms[dir](offsets[i][0], offsets[i][1])
		rectMask.Pos = pos + utils.Int2Vector(x*64, y*64) + 1 + 1i
		player := PlayerManager.GetPlayer(func(player *Player) bool {
			return utils.PointCollision(rectMask, player.Grid.Pos+GridSize/2)
		})
		if player != nil {
			res = append(res, player)
		}
	}
	return res
}

func SortEnemyByLast(enemies []*Enemy) {
	sort.Slice(enemies, func(i, j int) bool {
		return enemies[i].GetLast() < enemies[i].GetLast()
	})
}

func SortPlayerByHurt(players []*Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Data.Hp-players[i].Hp > players[j].Data.Hp-players[j].Hp
	})
}

func InvokePosDraw(src any, pos complex128, screen *ebiten.Image) {
	if tar, ok := src.(IPosDraw); ok {
		tar.PosDraw(pos, screen)
	}
}
