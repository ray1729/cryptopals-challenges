package set1

import (
    "bufio"
    "encoding/base64"
    "encoding/hex"
    "log"
    "os"
    "unicode"
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
        res[i] = s_bytes[i]^k_bytes[i]
    }
    return hex.EncodeToString(res)
}

func single_byte_xor(b byte, s string) string {
    s_bytes := hex_bytes(s)
    n := len(s_bytes)
    res := make([]byte, n)
    for i := 0; i < n; i++ {
        res[i] = s_bytes[i]^b
    }
    return string(res)
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
    for c, ef := range(EnglishLetterFrequency) {
        f := float64(frequencies[c])*100.0/float64(n)
        d := ef - f
        distance += d*d
    }
    // Penalize non-ascii characters that appear in `s`
    // for c, f := range(frequencies) {
    //     if _, ok := EnglishLetterFrequency[c]; !ok {
    //         distance += 100.0*float64(f)
    //     }
    // }
    return 1.0/(1.0+distance)
}

func solve_single_byte_xor(s string) (string, float64) {
    best := ""
    best_score := 0.0
    for b := byte(0); b < 255; b++ {
        t := single_byte_xor(b, s)
        s := score(t)
        if s > best_score {
            best = t
            best_score = s
        }
    }
    return best, best_score
}

func detect_single_character_xor(filename string) string {
    f, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    best := ""
    best_score := 0.0
    for scanner.Scan() {
        this, this_score := solve_single_byte_xor(scanner.Text())
        if this_score > best_score {
            best = this
            best_score = this_score
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return best
}

func repeating_key_xor(s, k string) string {
    s_bytes := []byte(s)
    k_bytes := []byte(k)
    s_len := len(s_bytes)
    k_len := len(k_bytes)
    res := make([]byte, s_len)
    for i := 0; i < s_len; i++ {
        res[i] = s_bytes[i]^k_bytes[i % k_len]
    }
    return hex.EncodeToString(res)
}


// From "The Go Programming Language", p45
var pc [256]byte

func init() {
    for i := range pc {
        pc[i] = pc[i/2] + byte(i&1)
    }
}

func hamming_distance(s, t string) int {
    s_bytes := []byte(s)
    s_len := len(s_bytes)
    t_bytes := []byte(t)
    t_len := len(t_bytes)
    if s_len != t_len {
        log.Fatal("hamming_distance not implemented for strings of differing length")
    }
    var distance int
    for i := 0; i < s_len; i++ {
        distance += int(pc[s_bytes[i]^t_bytes[i]])
    }
    return distance
}
