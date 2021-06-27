package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"time"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var ███████████████████████████████

	// player
	player          = playerblok{}
	playermovepause bool
	// map
	actualmapbottomx, actualmapbottomy  int
	textureon, worlddetailon, outlineon bool
	drawblok, drawbloknext              int
	selectblok, selectblokworld         int
	tilew                               = 128
	tileh                               = 64
	draww                               = 10
	drawh                               = 10
	drawa                               = draww * drawh
	drawactualh                         = tileh * drawh
	drawactualw                         = tilew * draww
	worldw                              = 1000
	worldh                              = 1000
	worlda                              = worldw * worldh
	world                               = make([]blok, worlda)
	worldunderground                    = make([]blok, worlda)

	//img
	flame     = rl.NewRectangle(47, 115, 18, 16)
	grass1    = rl.NewRectangle(0, 0, 128, 90)
	playerimg = rl.NewRectangle(0, 93, 16, 16)
	tree1     = rl.NewRectangle(0, 352, 32, 32)

	trees = make([]rl.Rectangle, 15)
	// core
	centerlineson, options, paused, scanlines, pixelnoise, ghosting bool
	mouseblok                                                       int
	mousepos                                                        rl.Vector2
	gridon, debugon, fadeblinkon                                    bool
	monw, monh                                                      int
	fps                                                             = 30
	framecount                                                      int
	imgs                                                            rl.Texture2D
	camera, cameraend                                               rl.Camera2D
	fadeblink                                                       = float32(0.2)
	onoff2, onoff3, onoff6, onoff10, onoff15, onoff30, onoff60      bool
)

// MARK: structs
type blok struct {
	activ, solid           bool
	img, imgsolid          rl.Rectangle
	color, color2          rl.Color
	name, use1, use2, use3 string
}
type playerblok struct {
	move                         bool
	x, y, bloknumber, drawnumber int
	img                          rl.Rectangle
}

func raylib() { // MARK: raylib
	rl.InitWindow(monw, monh, "GAME TITLE")
	//rl.ToggleFullscreen()
	rl.SetExitKey(rl.KeyEnd) // key to end the game and close window
	// MARK: load images
	imgs = rl.LoadTexture("imgs.png") // load images
	createimgs()
	createmaps()
	paused = true

	rl.SetTargetFPS(fps)
	//rl.HideCursor()
	//	rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		framecount++
		mousepos = rl.GetMousePosition()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//	drawnocameraback()
		rl.BeginMode2D(camera)

		drawlayers()

		if gridon {
			drawgrid()
		}

		rl.EndMode2D()
		drawnocamera()

		if debugon {
			drawdebug()
		}
		update()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func update() { // MARK: update
	input()
	timers()

	updateplayer()
	updatedrawmap()
}
func updatedrawmap() { // MARK: updateplayer

	if player.bloknumber >= 5005 {
		drawbloknext = player.bloknumber - ((worldw + 1) * (drawh / 2))
	} else {
		drawbloknext = 0
	}

}

func updateplayer() { // MARK: updateplayer
	if selectblokworld != player.bloknumber && selectblok != 0 {
		player.move = true
	} else {
		player.move = false
	}

	if onoff3 {
		if !playermovepause && selectblok != 0 {
			if selectblokworld != player.bloknumber {
				if selectblok > player.drawnumber { // up left movements
					if (selectblok-player.drawnumber)%10 == 0 {
						if !world[player.bloknumber+worldw].solid {
							player.bloknumber += worldw
						}
					} else if (selectblok-player.drawnumber)%9 == 0 {
						if !world[player.bloknumber+(worldw-1)].solid {
							player.bloknumber += worldw - 1
						}
					} else if (selectblok-player.drawnumber)%11 == 0 {
						if !world[player.bloknumber+(worldw+1)].solid {
							player.bloknumber += worldw + 1
						}
					} else if (selectblok - player.drawnumber) < 10 {
						if !world[player.bloknumber+1].solid {
							player.bloknumber += 1
						}
					} else {
						switch rInt(1, 8) {
						case 1:
							if !world[player.bloknumber-(worldw-1)].solid {
								player.bloknumber -= worldw - 1
							}
						case 2:
							if !world[player.bloknumber+1].solid {
								player.bloknumber += 1
							}
						case 3:
							if !world[player.bloknumber+(worldw+1)].solid {
								player.bloknumber += worldw + 1
							}
						case 4:
							if !world[player.bloknumber+worldw].solid {
								player.bloknumber += worldw
							}
						case 5:
							if !world[player.bloknumber+(worldw-1)].solid {
								player.bloknumber += worldw - 1
							}
						case 6:
							if !world[player.bloknumber-1].solid {
								player.bloknumber -= 1
							}
						case 7:
							if !world[player.bloknumber-(worldw+1)].solid {
								player.bloknumber -= worldw + 1
							}
						}
					}
				} else if selectblok < player.drawnumber { // right down movements
					if (player.drawnumber-selectblok)%10 == 0 {
						if !world[player.bloknumber-worldw].solid {
							player.bloknumber -= worldw
						}
					} else if (player.drawnumber-selectblok)%9 == 0 {
						if !world[player.bloknumber-(worldw-1)].solid {
							player.bloknumber -= worldw - 1
						}
					} else if (player.drawnumber-selectblok)%11 == 0 {
						if !world[player.bloknumber-(worldw+1)].solid {
							player.bloknumber -= worldw + 1
						}
					} else if (player.drawnumber - selectblok) < 10 {
						if !world[player.bloknumber-1].solid {
							player.bloknumber -= 1
						}
					} else {
						switch rInt(1, 8) {
						case 1:
							if !world[player.bloknumber+(worldw+1)].solid {
								player.bloknumber += worldw + 1
							}
						case 2:
							if !world[player.bloknumber+1].solid {
								player.bloknumber += 1
							}
						case 3:
							if !world[player.bloknumber-(worldw-1)].solid {
								player.bloknumber -= worldw - 1
							}
						case 4:
							if !world[player.bloknumber-worldw].solid {
								player.bloknumber -= worldw
							}
						case 5:
							if !world[player.bloknumber-(worldw+1)].solid {
								player.bloknumber -= worldw + 1
							}
						case 6:
							if !world[player.bloknumber-1].solid {
								player.bloknumber -= 1
							}
						case 7:
							if !world[player.bloknumber+(worldw-1)].solid {
								player.bloknumber += worldw - 1
							}
						}
					}

				}
			}
		}
		playermovepause = true
	} else {
		playermovepause = false
		//	updatedrawmap()
	}
}
func drawnocamera() { // MARK: drawnocamera

	if scanlines {
		for a := 0; a < monh; a++ {
			rl.DrawLine(0, a, monw, a, rl.Fade(rl.Black, 0.2))
			a += 2
		}
	}
	if pixelnoise {
		for a := 0; a < 100; a++ {
			width := rFloat32(1, 3)
			rec := rl.NewRectangle(rFloat32(0, monw), rFloat32(0, monh), width, width)
			rl.DrawRectangleRec(rec, rl.Fade(rl.Black, rF32(0.4, 1.1)))
		}

	}

}
func drawlayers() { // MARK: drawlayers

	// layer 1
	v1 := rl.NewVector2(float32(monw/2), float32(monh-(monh-drawactualh)/2))
	actualmapbottomy = int(v1.Y)
	actualmapbottomx = int(v1.X)
	v1orig := v1
	count := 0
	drawblok = drawbloknext
	for a := 0; a < drawa; a++ {

		v2 := rl.NewVector2(v1.X-float32(tilew/2), v1.Y-float32(tileh/2))
		v3 := rl.NewVector2(v1.X, v1.Y-float32(tileh))
		v4 := rl.NewVector2(v1.X+float32(tilew/2), v1.Y-float32(tileh/2))

		if world[drawblok].activ && world[drawblok].name != "water" {

			rl.DrawTriangle(v2, v1, v3, rl.Fade(world[drawblok].color, 0.2))
			rl.DrawTriangle(v3, v1, v4, rl.Fade(world[drawblok].color, 0.2))

			if ghosting {
				change := rF32(-0.5, 1.0)

				v1new := rl.NewVector2(v1.X+change, v1.Y+change)
				v2new := rl.NewVector2(v2.X+change, v2.Y+change)
				v3new := rl.NewVector2(v3.X+change, v3.Y+change)
				v4new := rl.NewVector2(v4.X+change, v4.Y+change)

				rl.DrawTriangle(v2new, v1new, v3new, rl.Fade(world[drawblok].color, 0.04))
				rl.DrawTriangle(v3new, v1new, v4new, rl.Fade(world[drawblok].color, 0.04))

			}

			rl.DrawLineV(v1, v2, rl.Fade(rl.Black, 0.1))
			rl.DrawLineV(v2, v3, rl.Fade(rl.Black, 0.1))
			rl.DrawLineV(v3, v4, rl.Fade(rl.Black, 0.1))
			rl.DrawLineV(v1, v4, rl.Fade(rl.Black, 0.1))

			if worlddetailon {
				if world[drawblok].solid && world[drawblok].name == "tree" {
					destrec := rl.NewRectangle(v1.X, v1.Y-60, 72, 72)
					origin := rl.NewVector2(36, 36)
					rl.DrawTexturePro(imgs, world[drawblok].imgsolid, destrec, origin, 0, rl.White)
					if ghosting {
						changex := rFloat32(-3, 4)
						changey := rFloat32(-3, 4)
						destrec.X += changex
						destrec.Y += changey

						rl.DrawTexturePro(imgs, world[drawblok].imgsolid, destrec, origin, 0, rl.Fade(rl.White, 0.2))
					}

				}
			}

			// tile img
			//	imgv2 := rl.NewVector2(v2.X, v3.Y-24)
			//	rl.DrawTextureRec(imgs, grass1, imgv2, rl.White)

		} else if world[drawblok].activ && world[drawblok].name == "water" {

			smltileh := tileh / 4
			smltilew := tilew / 4
			v1sml := v1
			v1sml.Y += 3.0
			yorig := v1sml.Y
			xorig := v1sml.X
			count2 := 0
			for a := 0; a < 16; a++ {
				v2sml := rl.NewVector2(v1sml.X-float32(smltilew/2), v1sml.Y-float32(smltileh/2))
				v3sml := rl.NewVector2(v1sml.X, v1sml.Y-float32(smltileh))
				v4sml := rl.NewVector2(v1sml.X+float32(smltilew/2), v1sml.Y-float32(smltileh/2))

				blue := randombluedark()
				if flipcoin() {
					if flipcoin() {
						rl.DrawTriangle(v2sml, v1sml, v3sml, rl.Fade(blue, 0.2))
						rl.DrawTriangle(v3sml, v1sml, v4sml, rl.Fade(blue, 0.2))
					} else {
						change := rFloat32(2, 6)
						v1new := rl.NewVector2(v1sml.X, v1sml.Y-change)
						v2new := rl.NewVector2(v2sml.X, v2sml.Y-change)
						v3new := rl.NewVector2(v3sml.X, v3sml.Y-change)
						v4new := rl.NewVector2(v4sml.X, v4sml.Y-change)
						rl.DrawTriangle(v2new, v1new, v3new, rl.Fade(blue, 0.2))
						rl.DrawTriangle(v3new, v1new, v4new, rl.Fade(blue, 0.2))
					}

				} else {
					if flipcoin() {
						rl.DrawTriangle(v2sml, v1sml, v3sml, rl.Fade(blue, 0.4))
						rl.DrawTriangle(v3sml, v1sml, v4sml, rl.Fade(blue, 0.4))
					} else {

						change := rFloat32(2, 6)
						v1new := rl.NewVector2(v1sml.X, v1sml.Y-change)
						v2new := rl.NewVector2(v2sml.X, v2sml.Y-change)
						v3new := rl.NewVector2(v3sml.X, v3sml.Y-change)
						v4new := rl.NewVector2(v4sml.X, v4sml.Y-change)
						rl.DrawTriangle(v2new, v1new, v3new, rl.Fade(blue, 0.4))
						rl.DrawTriangle(v3new, v1new, v4new, rl.Fade(blue, 0.4))

					}
				}

				v1sml.Y -= float32(smltileh) / 2
				v1sml.X -= float32(smltilew) / 2

				count2++
				if count2 == 4 {
					count2 = 0
					v1sml.X = xorig
					v1sml.X += float32(smltilew / 2)
					xorig = v1sml.X
					v1sml.Y = yorig
					v1sml.Y -= float32(smltileh / 2)
					yorig = v1sml.Y
				}
			}

		}
		if rl.CheckCollisionPointTriangle(mousepos, v2, v1, v3) || rl.CheckCollisionPointTriangle(mousepos, v3, v1, v4) {
			rl.DrawTriangle(v2, v1, v3, rl.Fade(rl.Magenta, fadeblink))
			rl.DrawTriangle(v3, v1, v4, rl.Fade(rl.Magenta, fadeblink))

			//blok text
			textx := actualmapbottomx + 80
			texty := actualmapbottomy - 20
			textlen := rl.MeasureText(world[drawblok].name, 20)
			rl.DrawText(world[drawblok].name, textx, texty, 20, rl.Magenta)
			rl.DrawText(world[drawblok].use1, textx+textlen+8, texty-8, 10, rl.Magenta)
			rl.DrawText(world[drawblok].use2, textx+textlen+8, texty+4, 10, rl.Magenta)
			rl.DrawText(world[drawblok].use3, textx+textlen+8, texty+16, 10, rl.Magenta)

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				selectblok = a
				selectblokworld = drawblok
			}
		}

		if outlineon {
			rl.DrawLineV(v1, v2, rl.Black)
			rl.DrawLineV(v2, v3, rl.Black)
			rl.DrawLineV(v3, v4, rl.Black)
			rl.DrawLineV(v1, v4, rl.Black)
		}

		// draw iso grid
		if gridon {
			rl.DrawLineV(v1, v2, rl.Magenta)
			rl.DrawLineV(v2, v3, rl.Magenta)
			rl.DrawLineV(v3, v4, rl.Magenta)
			rl.DrawLineV(v1, v4, rl.Magenta)

			number := strconv.Itoa(a)
			number2 := strconv.Itoa(drawblok)
			rl.DrawText(number, int(v1.X-10), int(v3.Y+10), 20, rl.White)
			rl.DrawText(number2, int(v1.X-20), int(v2.Y), 10, rl.White)
		}

		v1.X -= float32(tilew / 2)
		v1.Y -= float32(tileh / 2)

		count++
		drawblok++
		if count == draww {
			count = 0
			drawblok -= draww
			drawblok += worldw
			v1 = v1orig
			v1.Y -= float32(tileh / 2)
			v1.X += float32(tilew / 2)
			v1orig = v1

		}
	}

	// layer 2
	v1 = rl.NewVector2(float32(monw/2), float32(monh-(monh-drawactualh)/2))
	v1orig = v1
	count = 0
	drawblok = drawbloknext
	for a := 0; a < drawa; a++ {

		//v2 := rl.NewVector2(v1.X-float32(tilew/2), v1.Y-float32(tileh/2))
		//v3 := rl.NewVector2(v1.X, v1.Y-float32(tileh))
		//v4 := rl.NewVector2(v1.X+float32(tilew/2), v1.Y-float32(tileh/2))

		// draw player
		if player.bloknumber == drawblok {

			player.x = int(v1.X)
			player.y = int(v1.Y - float32(tileh/2))
			player.drawnumber = a
			if onoff6 {
				destrec := rl.NewRectangle(v1.X, v1.Y-64, 64, 64)
				origin := rl.NewVector2(32, 32)

				rl.DrawTexturePro(imgs, player.img, destrec, origin, 0, rl.Magenta)

				//	rl.DrawCircle(int(v1.X), int(v1.Y-float32(tileh/2)), 15, rl.Fade(rl.Magenta, fadeblink))
			} else {

				destrec := rl.NewRectangle(v1.X, v1.Y-66, 64, 64)
				origin := rl.NewVector2(32, 32)

				rl.DrawTexturePro(imgs, player.img, destrec, origin, 0, rl.Magenta)

				//	v2 := rl.NewVector2(v1.X, v1.Y)
				//	rl.DrawTextureRec(imgs, player.img, v2, rl.White)
				//	rl.DrawCircle(int(v1.X), int(v1.Y-float32(tileh/2)-10), 15, rl.Fade(rl.Magenta, fadeblink))

			}
		}

		v1.X -= float32(tilew / 2)
		v1.Y -= float32(tileh / 2)

		count++
		drawblok++
		if count == draww {
			count = 0
			drawblok -= draww
			drawblok += worldw
			v1 = v1orig
			v1.Y -= float32(tileh / 2)
			v1.X += float32(tilew / 2)
			v1orig = v1

		}
	}

}

func createmaps() { // MARK: createmaps

	for a := 0; a < len(world); a++ {
		world[a].activ = true
		world[a].color = randomgreen()

		//create trees
		if rolldice() == 6 {
			world[a].solid = true
			world[a].name = "tree"
			world[a].use1 = "wood"
			world[a].use2 = "birds"
			world[a].use3 = "leaves"
			world[a].imgsolid = trees[rInt(0, len(trees))]
		}
	}

	for a := 0; a < len(world); a++ {

		// create water
		if a > 100000 && a < 900000 {

			if rolldice()+rolldice()+rolldice() == 18 {

				length := rInt(4, 17)
				origlen := length
				waterblok := a
				origblok := waterblok

				for {
					for a := 0; a < length; a++ {
						world[waterblok].color = randombluedark()
						world[waterblok].solid = true
						world[waterblok].name = "water"
						world[a].use1 = "fire"
						world[a].use2 = "ice"
						world[a].use3 = "fish"
						waterblok++
					}
					waterblok += worldw
					waterblok -= length
					change := rInt(1, 4)
					waterblok += change / 2
					length -= change
					if length <= 1 {
						break
					}
				}
				length = origlen
				waterblok = origblok
				for {
					for a := 0; a < length; a++ {
						world[waterblok].color = randombluedark()
						world[waterblok].solid = true
						world[waterblok].name = "water"
						world[a].use1 = "fire"
						world[a].use2 = "ice"
						world[a].use3 = "fish"
						waterblok--
					}
					waterblok -= worldw
					waterblok += length
					change := rInt(1, 4)
					waterblok += change / 2
					length -= change
					if length <= 1 {
						break
					}
				}

			}
		}

	}

	//water boundary left bottom
	for a := 0; a < 3000; a++ {
		world[a].color = randombluedark()
		world[a].solid = true
		world[a].name = "water"
		world[a].use1 = "fire"
		world[a].use2 = "ice"
		world[a].use3 = "fish"
		if rolldice() == 6 {
			world[a+worldw].color = randombluedark()
			world[a+worldw].solid = true
			world[a+worldw].name = "water"
			world[a].use1 = "fire"
			world[a].use2 = "ice"
			world[a].use3 = "fish"
			if flipcoin() {
				world[a+worldw*2].color = randombluedark()
				world[a+worldw*2].solid = true
				world[a+worldw*2].name = "water"
			}

		}

	}

}
func createimgs() { // MARK: createimgs

	count := 0
	for a := 0; a < len(trees); a++ {
		trees[a] = rl.NewRectangle(tree1.X, tree1.Y, 32, 32)
		tree1.X += 32
		count++
		if count == 4 {
			count = 0
			tree1.X = 0
			tree1.Y += 32
		}
	}

}

// MARK: core	core	core	core	core	core	core	core	core	core	core
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLogLevel(rl.LogError) // hides info window
	rl.InitWindow(monw, monh, "setres")
	setres(0, 0)
	rl.CloseWindow()
	setinitialvalues()
	raylib()

}
func input() { // MARK: input
	if rl.IsKeyPressed(rl.KeyF2) {
		if textureon {
			textureon = false
		} else {
			textureon = true
		}
	}
	if rl.IsKeyPressed(rl.KeyF1) {
		if worlddetailon {
			worlddetailon = false
		} else {
			worlddetailon = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpAdd) {
		camera.Zoom += 0.2
		camera.Offset.X = float32(player.x)
		camera.Offset.Y = float32(player.y)
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {

		camera.Zoom -= 0.2
		camera.Target.X = float32(monw - player.x)
		camera.Target.Y = float32(monh - player.y)
	}
	if rl.IsKeyPressed(rl.KeyPause) {
		if paused {
			paused = false
		} else {
			paused = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpDivide) {
		if centerlineson {
			centerlineson = false
		} else {
			centerlineson = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
		if gridon {
			gridon = false
		} else {
			gridon = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKp0) {
		if gridon {
			gridon = false
		} else {
			gridon = true
		}
	}
}
func drawdebug() { // MARK: DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG DEBUG

	//centerlines
	if centerlineson {
		rl.DrawLine(monw/2, 0, monw/2, monh, rl.Magenta)
		rl.DrawLine(0, monh/2, monw, monh/2, rl.Magenta)
	}
	rl.DrawRectangle(monw-300, 0, 300, monh, rl.Fade(rl.Black, 0.5))
	textx := monw - 290
	textx2 := monw - 145
	texty := 10

	drawactualhtext := strconv.Itoa(drawactualh)
	drawactualwtext := strconv.Itoa(drawactualw)
	playerdrawnumbertext := strconv.Itoa(player.drawnumber)
	playerbloknumbertexdt := strconv.Itoa(player.bloknumber)
	selectbloktext := strconv.Itoa(selectblok)
	selectblokworldtext := strconv.Itoa(selectblokworld)

	camerazoomtext := fmt.Sprintf("%g", camera.Zoom)
	mousex := fmt.Sprintf("%g", mousepos.X)
	mousey := fmt.Sprintf("%g", mousepos.Y)

	rl.DrawText("camerazoomtext", textx, texty, 10, rl.White)
	rl.DrawText(camerazoomtext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("mousex", textx, texty, 10, rl.White)
	rl.DrawText(mousex, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("mousey", textx, texty, 10, rl.White)
	rl.DrawText(mousey, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("drawactualhtext", textx, texty, 10, rl.White)
	rl.DrawText(drawactualhtext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("drawactualwtext", textx, texty, 10, rl.White)
	rl.DrawText(drawactualwtext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("playerdrawnumbertext", textx, texty, 10, rl.White)
	rl.DrawText(playerdrawnumbertext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("playerbloknumbertexdt", textx, texty, 10, rl.White)
	rl.DrawText(playerbloknumbertexdt, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("selectbloktext", textx, texty, 10, rl.White)
	rl.DrawText(selectbloktext, textx2, texty, 10, rl.White)
	texty += 12
	rl.DrawText("selectblokworldtext", textx, texty, 10, rl.White)
	rl.DrawText(selectblokworldtext, textx2, texty, 10, rl.White)
	texty += 12

	// fps
	rl.DrawRectangle(monw-110, monh-110, 100, 40, rl.Black)
	rl.DrawFPS(monw-100, monh-100)

}
func timers() { // MARK: timers
	if player.move {
		if onoff2 {
			player.img.X += 16
			if player.img.X > 70 {
				player.img.X = 0
			}

		}
	}
	if framecount%2 == 0 {
		if onoff2 {
			onoff2 = false
		} else {
			onoff2 = true
		}
	}
	if framecount%3 == 0 {
		if onoff3 {
			onoff3 = false
		} else {
			onoff3 = true
		}
	}
	if framecount%6 == 0 {
		if onoff6 {
			onoff6 = false
		} else {
			onoff6 = true
		}
	}
	if framecount%10 == 0 {
		if onoff10 {
			onoff10 = false
		} else {
			onoff10 = true
		}
	}
	if framecount%15 == 0 {
		if onoff15 {
			onoff15 = false
		} else {
			onoff15 = true
		}
	}
	if framecount%30 == 0 {
		if onoff30 {
			onoff30 = false
		} else {
			onoff30 = true
		}
	}
	if framecount%60 == 0 {
		if onoff60 {
			onoff60 = false
		} else {
			onoff60 = true
		}
	}
	if fadeblinkon {
		if fadeblink > 0.2 {
			fadeblink -= 0.05
		} else {
			fadeblinkon = false
		}
	} else {
		if fadeblink < 0.6 {
			fadeblink += 0.05
		} else {
			fadeblinkon = true
		}
	}
}

func setres(w, h int) { // MARK: setres

	if w == 0 {
		monw = rl.GetMonitorWidth(0)
		monh = rl.GetMonitorHeight(0)

		camera.Target.X = 0
		camera.Target.Y = 0
		camera.Offset.X = 0
		camera.Offset.Y = 0
		camera.Zoom = 1.0

	} else {
		monw = w
		monh = h
	}

}

func setinitialvalues() { // MARK: setinitialvalues

	pixelnoise = true
	scanlines = true
	ghosting = true
	worlddetailon = true
	player.bloknumber = worldw * (worldh / 2)
	player.bloknumber += worldw / 2
	updateplayer()
	player.img = playerimg

}

// MARK:  █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █ █
func drawgrid() { // MARK: drawgrid

	x := 16
	for {
		rl.DrawLine(x, 0, x, monh, rl.Fade(rl.Magenta, 0.1))
		x += 16
		if x > monw {
			break
		}
	}
	y := 16
	for {
		rl.DrawLine(0, y, monw, y, rl.Fade(rl.Magenta, 0.1))
		y += 16
		if y > monh {
			break
		}
	}

}

// MARK: colors
// https://www.rapidtables.com/web/color/RGB_Color.html
func darkred() rl.Color {
	color := rl.NewColor(55, 0, 0, 255)
	return color
}
func semidarkred() rl.Color {
	color := rl.NewColor(70, 0, 0, 255)
	return color
}
func brightred() rl.Color {
	color := rl.NewColor(230, 0, 0, 255)
	return color
}
func randomgrey() rl.Color {
	color := rl.NewColor(uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(160, 193)), uint8(rInt(0, 255)))
	return color
}
func randombluelight() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 180)), uint8(rInt(120, 256)), uint8(rInt(120, 256)), 255)
	return color
}
func randombluedark() rl.Color {
	color := rl.NewColor(0, 0, uint8(rInt(120, 250)), 255)
	return color
}
func randomyellow() rl.Color {
	color := rl.NewColor(255, uint8(rInt(150, 256)), 0, 255)
	return color
}
func randomorange() rl.Color {
	color := rl.NewColor(uint8(rInt(250, 256)), uint8(rInt(60, 210)), 0, 255)
	return color
}
func randomred() rl.Color {
	color := rl.NewColor(uint8(rInt(128, 256)), uint8(rInt(0, 129)), uint8(rInt(0, 129)), 255)
	return color
}
func randomgreen() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 170)), uint8(rInt(100, 256)), uint8(rInt(0, 50)), 255)
	return color
}
func randomcolor() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 256)), uint8(rInt(0, 256)), uint8(rInt(0, 256)), 255)
	return color
}
func brightyellow() rl.Color {
	color := rl.NewColor(uint8(255), uint8(255), uint8(0), 255)
	return color
}
func brightbrown() rl.Color {
	color := rl.NewColor(uint8(218), uint8(165), uint8(32), 255)
	return color
}
func brightgrey() rl.Color {
	color := rl.NewColor(uint8(212), uint8(212), uint8(213), 255)
	return color
}

// random numbers
func rF32(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
