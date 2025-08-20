package utils

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/goombaio/namegenerator"
)

func GenerateRandomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	randomNumber := rand.Int()
	return fmt.Sprintf("%s-%d", name, randomNumber)
}
