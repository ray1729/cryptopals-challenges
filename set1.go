package set1

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	//"math"
	"os"
	"sort"
	"unicode"

	"github.com/montanaflynn/stats"
)

func hex_bytes(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func hex_to_base64(s string) string {
	return base64.StdEncoding.EncodeToString(hex_bytes(s))
}

func fixed_xor(s, k string) string {
	s_bytes := hex_bytes(s)
	k_bytes := hex_bytes(k)
	n := len(s_bytes)
	if len(k_bytes) != n {
		log.Fatalf("Inputs to fixed_xor must be the same length")
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = s_bytes[i] ^ k_bytes[i]
	}
	return hex.EncodeToString(res)
}

func single_byte_xor(b byte, s_bytes []byte) []byte {
	n := len(s_bytes)
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = s_bytes[i] ^ b
	}
	return res
}

/* Letter frequencies from

https://en.wikipedia.org/wiki/Letter_frequency

I took a guess at 13% for space as wp just says "In English, the
space is slightly more frequent than the top letter (e)..."

It might be useful to consider n-grams, e.g. http://norvig.com/mayzner.html
or tokenize the string and score against a dictionary.
*/

var EnglishLetterFrequency = map[rune]float64{
	'a': 8.167,
	'b': 1.492,
	'c': 2.782,
	'd': 4.253,
	'e': 12.702,
	'f': 2.228,
	'g': 2.015,
	'h': 6.094,
	'i': 6.966,
	'j': 0.153,
	'k': 0.772,
	'l': 4.205,
	'm': 2.406,
	'n': 6.749,
	'o': 7.507,
	'p': 1.929,
	'q': 0.095,
	'r': 5.987,
	's': 6.327,
	't': 9.056,
	'u': 2.758,
	'v': 0.978,
	'w': 2.360,
	'x': 0.150,
	'y': 1.974,
	'z': 0.074,
	' ': 13.0,
}

func score(s string) float64 {
	frequencies := make(map[rune]int)
	n := 0
	for _, c := range s {
		n++
		frequencies[unicode.ToLower(c)]++
	}
	distance := 0.0
	// Score English letter characters (and space)
	for c, ef := range EnglishLetterFrequency {
		f := float64(frequencies[c]) * 100.0 / float64(n)
		d := ef - f
		distance += d * d
	}
	// Penalize non-ascii characters that appear in `s`
	// for c, f := range(frequencies) {
	//     if _, ok := EnglishLetterFrequency[c]; !ok {
	//         distance += 100.0*float64(f)
	//     }
	// }
	return 1.0 / (1.0 + distance)
}

func letter_freq_correlation(s string) float64 {
	frequencies := make(map[rune]float64)
	n := 0.0
	for _, c := range s {
		n++
		lc := unicode.ToLower(c)
		frequencies[lc]++
	}
	var xs, ys stats.Float64Data
	for char, freq := range EnglishLetterFrequency {
		xs = append(xs, freq)
		ys = append(ys, frequencies[char]/n)
	}
	for char, count := range frequencies {
		if _, ok := EnglishLetterFrequency[char]; ok {
			continue
		}
		xs = append(xs, 0.0)
		ys = append(ys, count/n)
	}
	correlation, err := xs.Correlation(ys)
	if err != nil {
		log.Fatalf("Error computing correlation: %v", err)
	}
	return correlation
}

func solve_single_byte_xor(s []byte) ([]byte, byte, float64) {
	var best []byte
	var key byte
	best_score := 0.0
	for b := byte(0); b < 255; b++ {
		t := single_byte_xor(b, s)
		s := letter_freq_correlation(string(t))
		if s > best_score {
			best = t
			key = b
			best_score = s
		}
	}
	return best, key, best_score
}

func detect_single_character_xor(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var best []byte
	best_score := 0.0
	for scanner.Scan() {
		cyphertext, _ := hex.DecodeString(scanner.Text())
		this, _, this_score := solve_single_byte_xor(cyphertext)
		if this_score > best_score {
			best = this
			best_score = this_score
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return string(best)
}

func repeating_key_xor(s_bytes, k_bytes []byte) []byte {
	s_len := len(s_bytes)
	k_len := len(k_bytes)
	res := make([]byte, s_len)
	for i := 0; i < s_len; i++ {
		res[i] = s_bytes[i] ^ k_bytes[i%k_len]
	}
	return res
}

// From "The Go Programming Language", p45
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func hamming_distance(s_bytes, t_bytes []byte) int {
	s_len := len(s_bytes)
	t_len := len(t_bytes)
	if s_len != t_len {
		log.Fatal("hamming_distance not implemented for slices of differing length")
	}
	var distance int
	for i := 0; i < s_len; i++ {
		distance += int(pc[s_bytes[i]^t_bytes[i]])
	}
	return distance
}

func read_base64_file(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := base64.NewDecoder(base64.StdEncoding, f)
	var res []byte
	buffer := make([]byte, 1024)
	for {
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		res = append(res, buffer...)
	}
	return res, nil
}

// func guess_keysize(cyphertext []byte, n_delta int) int {
//     best_score := math.Inf(0)
//     best_keysize := 0
//     for keysize := 2; keysize < 41; keysize++ {
//         var delta int
//         for d := 0; d < n_delta; d++ {
//             x := d*keysize
//             y := (d+1)*keysize
//             z := (d+2)*keysize
//             b1 := cyphertext[x:y]
//             b2 := cyphertext[y:z]
//             delta += hamming_distance(b1, b2)
//         }
//         score := float64(delta)/float64(keysize * n_delta)
//         if score < best_score {
//             best_score = score
//             best_keysize = keysize
//         }
//     }
//     return best_keysize
// }

type Pair struct {
	Key   int
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func rankByScore(scores map[int]float64) PairList {
	result := make(PairList, len(scores))
	i := 0
	for k, v := range scores {
		result[i] = Pair{k, v}
		i++
	}
	sort.Sort(result)
	return result
}

func score_key_sizes(min, max int, cyphertext []byte) PairList {
	scores := make(map[int]float64)
	for keysize := min; keysize < max+1; keysize++ {
		score := 0
		for n := 0; n < 30; n++ {
			if (n+2)*keysize > len(cyphertext) {
				break
			}
			b1 := cyphertext[n*keysize : (n+1)*keysize]
			b2 := cyphertext[(n+1)*keysize : (n+2)*keysize]
			score += hamming_distance(b1, b2)
		}
		scores[keysize] = float64(score) / float64(keysize)
	}
	return rankByScore(scores)
}

func guess_xor_key(cyphertext []byte, keysize int) []byte {
	key := make([]byte, keysize)
	for i := 0; i < keysize; i++ {
		var block []byte
		for j := i; j < len(cyphertext); j += keysize {
			block = append(block, cyphertext[j])
		}
		_, k, _ := solve_single_byte_xor(block)
		key[i] = k
	}
	return key
}

func decode_repeating_key_xor(cyphertext []byte) ([]byte, []byte) {
	var best_plaintext, best_key []byte
	best_score := 0.0
	for keysize := 2; keysize < 40; keysize++ {
		key := guess_xor_key(cyphertext, keysize)
		plaintext := repeating_key_xor(cyphertext, key)
		score := letter_freq_correlation(string(plaintext))
		if score > best_score {
			best_score = score
			best_plaintext = plaintext
			best_key = key
		}
	}
	return best_plaintext, best_key
}

// func find_key(cyphertext []byte) []byte {
//   scores := score_key_sizes(2, 40, cyphertext)
//   keysize := scores[0].Key
//   log.Printf("Trying key size %d\n", keysize)
//   key := guess_xor_key(cyphertext, keysize)
//   return key
// }
