package obfuscator

import "math/rand"

func randomRange(lower, upper int) int {
	return lower + rand.Intn(upper-lower+1)
}

func strShuffle(str string) string {
	inRune := []rune(str)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}
