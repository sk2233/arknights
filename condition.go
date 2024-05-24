/*
@author: sk
@date: 2023/3/5
*/
package main

func IsCareer(career int) func(player *Player) bool {
	return func(player *Player) bool {
		return player.Data.Career == career
	}
}

func IsAny(*Player) bool {
	return true
}
