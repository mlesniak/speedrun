package seed

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Seed struct {
	Code string
	Seed int64
}

func New() Seed {
	return NewWords(3)
}

func NewWords(number int) Seed {
	words := randomWords(number)
	text := strings.Join(words, "-")
	seed := computeSeed(text)

	return Seed{
		Code: text,
		Seed: seed,
	}
}

func computeSeed(words string) int64 {
	sha256 := sha256.New()
	sha256.Write([]byte(words))
	hash := sha256.Sum(nil)

	// There are probably faster ways, but for now it suffices. We also ignore overflows, etc. since it's simply important
	// that the results are reproducible.
	k := fmt.Sprintf("%x", hash)
	numbers := k
	for _, c := range "abcdef" {
		numbers = strings.ReplaceAll(numbers, string(c), "")
	}
	numbers = numbers[:len(fmt.Sprintf("%d", math.MaxInt64))-1]

	seed, err := strconv.ParseInt(numbers, 10, 64)
	if err != nil {
		panic(err)
	}
	return seed
}

func randomWords(number int) []string {
	file, err := ioutil.ReadFile("words")
	if err != nil {
		panic(err)
	}

	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	words := strings.Split(string(file), "\n")
	result := make([]string, number)
	for i := 0; i < number; i++ {
		result[i] = words[rnd.Int()%len(words)]
		result[i] = strings.ToLower(result[i])
	}

	return result
}
