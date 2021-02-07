package main

import (
	"fmt"
	"github.com/bbawn/boredgames/internal/games/set"
	"os"
)

// Utility to rename files from old numeric scheme to Card string format:
// <color><count><shape><shadingL
func main() {
	i := 1
	for _, shading := range []set.Shading{set.Filled, set.Stripe, set.Outline} {
		for _, shape := range []set.Shape{set.Squiggle, set.Diamond, set.Oval} {
			for _, color := range []set.Color{set.Red, set.Purple, set.Green} {
				for count := 1; count < 4; count++ {
					c := set.Card{Color: color, Count: byte(count), Shading: shading, Shape: shape}
					cs := c.String()
					oldFile := fmt.Sprintf("ui/img/%d.gif", i)
					newFile := fmt.Sprintf("ui/img/%s.gif", cs)
					fmt.Printf("renaming %s to %s\n", oldFile, newFile)
					err := os.Rename(oldFile, newFile)
					if err != nil {
						fmt.Println("Rename failed:", err)
					}
					i += 1
				}
			}
		}
	}
}
