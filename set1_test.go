package set1

import (
    "testing"
)

func TestChallenge1(t *testing.T) {
    input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
    expected := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
    observed := hex_to_base64(input)
    if observed != expected {
        t.Errorf("hex_to_base64(%s) != %s", input, expected)
    }
}

func TestChallenge2(t *testing.T) {
    s := "1c0111001f010100061a024b53535009181c"
    k := "686974207468652062756c6c277320657965"
    expected := "746865206b696420646f6e277420706c6179"
    observed := fixed_xor(s, k)
    if observed != expected {
        t.Errorf("fixed_xor(%s,%s) != %s", s, k, expected)
    }
}

func TestChallenge3(t *testing.T) {
    s := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
    expected := "Cooking MC's like a pound of bacon"
    observed, _ := solve_single_byte_xor(s)
    if observed != expected {
        t.Errorf("solve_single_byte_xor(%s) != %s", s, expected)
    }
}

func TestChallenge4(t *testing.T) {
    filename := "data/4.txt"
    expected := "Now that the party is jumping\n"
    observed := detect_single_character_xor(filename)
    if observed != expected {
        t.Errorf("detect_single_character_xor(%s) != %s", filename, expected)
    }
}

func TestChallenge5(t *testing.T) {
    input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
    expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272" +
        "a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
    observed := repeating_key_xor(input, "ICE")
    if observed != expected {
        t.Errorf("repeating_key_xor() failed")
    }
}
