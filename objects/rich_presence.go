package dota2

import "fmt"

type RichPresence struct{
	SteamID			uint64
	RichPresenceKV	string
}

func (rp *RichPresence) String() string {
	return fmt.Sprintf("Player %d: %s", rp.SteamID, rp.RichPresenceKV)
}