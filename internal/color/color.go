package color

func RGBToDiscord(r, g, b int) int {
	return (r << 16) + (g << 8) + b
}
