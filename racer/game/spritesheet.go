package game

import (
	"github.com/tanema/amore/gfx"
)

var (
	background   *gfx.Image // our background image (loaded below)
	sprites      *gfx.Image // our spritesheet (loaded below)
	backgrounds  = map[string]*gfx.Quad{}
	spriteSheet  = map[string]*gfx.Quad{}
	billboards   = []*gfx.Quad{}
	plants       = []*gfx.Quad{}
	carQuads     = []*gfx.Quad{}
	sprite_scale float32
)

func initSpriteSheets() {
	background, _ = gfx.NewImage("images/background.png")
	sprites, _ = gfx.NewImage("images/sprites.png")

	backgrounds["hills"] = gfx.NewQuad(5, 5, 1280, 480, background.Width, background.Height)
	backgrounds["sky"] = gfx.NewQuad(5, 495, 1280, 480, background.Width, background.Height)
	backgrounds["trees"] = gfx.NewQuad(5, 985, 1280, 480, background.Width, background.Height)

	spriteSheet["palm_tree"] = gfx.NewQuad(5, 5, 215, 540, sprites.Width, sprites.Height)
	spriteSheet["billboard08"] = gfx.NewQuad(230, 5, 385, 265, sprites.Width, sprites.Height)
	spriteSheet["tree1"] = gfx.NewQuad(625, 5, 360, 360, sprites.Width, sprites.Height)
	spriteSheet["dead_tree1"] = gfx.NewQuad(5, 555, 135, 332, sprites.Width, sprites.Height)
	spriteSheet["billboard09"] = gfx.NewQuad(150, 555, 328, 282, sprites.Width, sprites.Height)
	spriteSheet["boulder3"] = gfx.NewQuad(230, 280, 320, 220, sprites.Width, sprites.Height)
	spriteSheet["column"] = gfx.NewQuad(995, 5, 200, 315, sprites.Width, sprites.Height)
	spriteSheet["billboard01"] = gfx.NewQuad(625, 375, 300, 170, sprites.Width, sprites.Height)
	spriteSheet["billboard06"] = gfx.NewQuad(488, 555, 298, 190, sprites.Width, sprites.Height)
	spriteSheet["billboard05"] = gfx.NewQuad(5, 897, 298, 190, sprites.Width, sprites.Height)
	spriteSheet["billboard07"] = gfx.NewQuad(313, 897, 298, 190, sprites.Width, sprites.Height)
	spriteSheet["boulder2"] = gfx.NewQuad(621, 897, 298, 140, sprites.Width, sprites.Height)
	spriteSheet["tree2"] = gfx.NewQuad(1205, 5, 282, 295, sprites.Width, sprites.Height)
	spriteSheet["billboard04"] = gfx.NewQuad(1205, 310, 268, 170, sprites.Width, sprites.Height)
	spriteSheet["dead_tree2"] = gfx.NewQuad(1205, 490, 150, 260, sprites.Width, sprites.Height)
	spriteSheet["boulder1"] = gfx.NewQuad(1205, 760, 168, 248, sprites.Width, sprites.Height)
	spriteSheet["bush1"] = gfx.NewQuad(5, 1097, 240, 155, sprites.Width, sprites.Height)
	spriteSheet["cactus"] = gfx.NewQuad(929, 897, 235, 118, sprites.Width, sprites.Height)
	spriteSheet["bush2"] = gfx.NewQuad(255, 1097, 232, 152, sprites.Width, sprites.Height)
	spriteSheet["billboard03"] = gfx.NewQuad(5, 1262, 230, 220, sprites.Width, sprites.Height)
	spriteSheet["billboard02"] = gfx.NewQuad(245, 1262, 215, 220, sprites.Width, sprites.Height)
	spriteSheet["stump"] = gfx.NewQuad(995, 330, 195, 140, sprites.Width, sprites.Height)
	spriteSheet["semi"] = gfx.NewQuad(1365, 490, 122, 144, sprites.Width, sprites.Height)
	spriteSheet["truck"] = gfx.NewQuad(1365, 644, 100, 78, sprites.Width, sprites.Height)
	spriteSheet["car03"] = gfx.NewQuad(1383, 760, 88, 55, sprites.Width, sprites.Height)
	spriteSheet["car02"] = gfx.NewQuad(1383, 825, 80, 59, sprites.Width, sprites.Height)
	spriteSheet["car04"] = gfx.NewQuad(1383, 894, 80, 57, sprites.Width, sprites.Height)
	spriteSheet["car01"] = gfx.NewQuad(1205, 1018, 80, 56, sprites.Width, sprites.Height)
	spriteSheet["player_uphill_left"] = gfx.NewQuad(1383, 961, 80, 45, sprites.Width, sprites.Height)
	spriteSheet["player_uphill_straight"] = gfx.NewQuad(1295, 1018, 80, 45, sprites.Width, sprites.Height)
	spriteSheet["player_uphill_right"] = gfx.NewQuad(1385, 1018, 80, 45, sprites.Width, sprites.Height)
	spriteSheet["player_left"] = gfx.NewQuad(995, 480, 80, 41, sprites.Width, sprites.Height)
	spriteSheet["player_straight"] = gfx.NewQuad(1085, 480, 80, 41, sprites.Width, sprites.Height)
	spriteSheet["player_right"] = gfx.NewQuad(995, 531, 80, 41, sprites.Width, sprites.Height)

	sprite_scale = 0.3 * (1 / spriteSheet["player_straight"].GetWidth()) // the reference sprite width should be 1/3rd the (half-)roadWidth

	billboards = []*gfx.Quad{spriteSheet["billboard01"], spriteSheet["billboard02"], spriteSheet["billboard03"], spriteSheet["billboard04"], spriteSheet["billboard05"], spriteSheet["billboard06"], spriteSheet["billboard07"], spriteSheet["billboard08"], spriteSheet["billboard09"]}
	plants = []*gfx.Quad{spriteSheet["tree1"], spriteSheet["tree2"], spriteSheet["dead_tree1"], spriteSheet["dead_tree2"], spriteSheet["palm_tree"], spriteSheet["bush1"], spriteSheet["bush2"], spriteSheet["cactus"], spriteSheet["stump"], spriteSheet["boulder1"], spriteSheet["boulder2"], spriteSheet["boulder3"]}
	carQuads = []*gfx.Quad{spriteSheet["car01"], spriteSheet["car02"], spriteSheet["car03"], spriteSheet["car04"], spriteSheet["semi"], spriteSheet["truck"]}
}
