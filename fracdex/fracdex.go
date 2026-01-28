// Licensed under CC0-1.0 Universial by https://github.com/rocicorp/fracdex
package fracdex

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	base62Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	smallestInt  = "A00000000000000000000000000"
	zero         = "a0"
)

// KeyBetween returns a key that sorts lexicographically between a and b.
// Either a or b can be empty strings. If a is empty it indicates smallest key,
// If b is empty it indicates largest key.
// b must be empty string or > a.
func KeyBetween(a, b string) (string, error) {
	if a != "" {
		err := validateOrderKey(a)
		if err != nil {
			return "", err
		}
	}
	if b != "" {
		err := validateOrderKey(b)
		if err != nil {
			return "", err
		}
	}
	if a != "" && b != "" && a >= b {
		return "", fmt.Errorf("%s >= %s", a, b)
	}
	if a == "" {
		if b == "" {
			return zero, nil
		}

		ib, err := getIntPart(b)
		if err != nil {
			return "", err
		}
		fb := b[len(ib):]
		if ib == smallestInt {
			return ib + midpoint("", fb), nil
		}
		if ib < b {
			return ib, nil
		}
		res, err := decrementInt(ib)
		if err != nil {
			return "", err
		}
		if res == "" {
			return "", errors.New("range underflow")
		}
		return res, nil
	}

	if b == "" {
		ia, err := getIntPart(a)
		if err != nil {
			return "", err
		}
		fa := a[len(ia):]
		i, err := incrementInt(ia)
		if err != nil {
			return "", err
		}
		if i == "" {
			return ia + midpoint(fa, ""), nil
		}
		return i, nil
	}

	ia, err := getIntPart(a)
	if err != nil {
		return "", err
	}
	fa := a[len(ia):]
	ib, err := getIntPart(b)
	if err != nil {
		return "", err
	}
	fb := b[len(ib):]
	if ia == ib {
		return ia + midpoint(fa, fb), nil
	}
	i, err := incrementInt(ia)
	if err != nil {
		return "", err
	}
	if i == "" {
		return "", errors.New("range overflow")
	}
	if i < b {
		return i, nil
	}
	return ia + midpoint(fa, ""), nil
}

// `a < b` lexicographically if `b` is non-empty.
// a == "" means first possible string.
// b == "" means last possible string.
func midpoint(a, b string) string {
	if b != "" {
		// remove longest common prefix.  pad `a` with 0s as we
		// go.  note that we don't need to pad `b`, because it can't
		// end before `a` while traversing the common prefix.
		i := 0
		for ; i < len(b); i++ {
			c := byte('0')
			if len(a) > i {
				c = a[i]
			}
			if c != b[i] {
				break
			}
		}
		if i > 0 {
			if i > len(a) {
				return b[0:i] + midpoint("", b[i:])
			}
			return b[0:i] + midpoint(a[i:], b[i:])
		}
	}

	// first digits (or lack of digit) are different
	digitA := 0
	if a != "" {
		digitA = strings.Index(base62Digits, string(a[0]))
	}
	digitB := len(base62Digits)
	if b != "" {
		digitB = strings.Index(base62Digits, string(b[0]))
	}
	if digitB-digitA > 1 {
		midDigit := int(math.Round(0.5 * float64(digitA+digitB)))
		return string(base62Digits[midDigit])
	}

	// first digits are consecutive
	if len(b) > 1 {
		return b[0:1]
	}

	// `b` is empty or has length 1 (a single digit).
	// the first digit of `a` is the previous digit to `b`,
	// or 9 if `b` is null.
	// given, for example, midpoint('49', '5'), return
	// '4' + midpoint('9', null), which will become
	// '4' + '9' + midpoint('', null), which is '495'
	sa := ""
	if len(a) > 0 {
		sa = a[1:]
	}
	return string(base62Digits[digitA]) + midpoint(sa, "")
}

func validateInt(i string) error {
	exp, err := getIntLen(i[0])
	if err != nil {
		return err
	}
	if len(i) != exp {
		return fmt.Errorf("invalid integer part of order key: %s", i)
	}
	return nil
}

func getIntLen(head byte) (int, error) {
	if head >= 'a' && head <= 'z' {
		return int(head - 'a' + 2), nil
	} else if head >= 'A' && head <= 'Z' {
		return int('Z' - head + 2), nil
	} else {
		return 0, fmt.Errorf("invalid order key head: %s", string(head))
	}
}

func getIntPart(key string) (string, error) {
	intPartLen, err := getIntLen(key[0])
	if err != nil {
		return "", err
	}
	if intPartLen > len(key) {
		return "", fmt.Errorf("invalid order key: %s", key)
	}
	return key[0:intPartLen], nil
}

func validateOrderKey(key string) error {
	if key == smallestInt {
		return fmt.Errorf("invalid order key: %s", key)
	}
	// getIntPart will return error if the first character is bad,
	// or the key is too short.  we'd call it to check these things
	// even if we didn't need the result
	i, err := getIntPart(key)
	if err != nil {
		return err
	}
	f := key[len(i):]
	if strings.HasSuffix(f, "0") {
		return fmt.Errorf("invalid order key: %s", key)
	}
	return nil
}

// returns error if x is invalid, or if range is exceeded.
func incrementInt(x string) (string, error) {
	err := validateInt(x)
	if err != nil {
		return "", err
	}
	digs := strings.Split(x, "")
	head := digs[0]
	digs = digs[1:]
	carry := true
	for i := len(digs) - 1; carry && i >= 0; i-- {
		d := strings.Index(base62Digits, digs[i]) + 1
		if d == len(base62Digits) {
			digs[i] = "0"
		} else {
			digs[i] = string(base62Digits[d])
			carry = false
		}
	}
	if carry {
		if head == "Z" {
			return "a0", nil
		}
		if head == "z" {
			return "", nil
		}
		h := string(head[0] + 1)
		if h > "a" {
			digs = append(digs, "0")
		} else {
			digs = digs[1:]
		}
		return h + strings.Join(digs, ""), nil
	}
	return head + strings.Join(digs, ""), nil
}

func decrementInt(x string) (string, error) {
	err := validateInt(x)
	if err != nil {
		return "", err
	}
	digs := strings.Split(x, "")
	head := digs[0]
	digs = digs[1:]
	borrow := true
	for i := len(digs) - 1; borrow && i >= 0; i-- {
		d := strings.Index(base62Digits, digs[i]) - 1
		if d == -1 {
			digs[i] = string(base62Digits[len(base62Digits)-1])
		} else {
			digs[i] = string(base62Digits[d])
			borrow = false
		}
	}

	if borrow {
		if head == "a" {
			return "Z" + string(base62Digits[len(base62Digits)-1]), nil
		}
		if head == "A" {
			return "", nil
		}
		h := head[0] - 1
		if h < 'Z' {
			digs = append(digs, string(base62Digits[len(base62Digits)-1]))
		} else {
			digs = digs[1:]
		}
		return string(h) + strings.Join(digs, ""), nil
	}

	return head + strings.Join(digs, ""), nil
}

// Float64Approx converts a key as generated by KeyBetween() to a float64.
// Because the range of keys is far larger than float64 can represent
// accurately, this is necessarily approximate. But for many use cases it should
// be, as they say, close enough for jazz.
func Float64Approx(key string) (float64, error) {
	if key == "" {
		return 0.0, errors.New("invalid order key")
	}

	err := validateOrderKey(key)
	if err != nil {
		return 0.0, err
	}

	ip, err := getIntPart(key)
	if err != nil {
		return 0.0, err
	}

	digs := strings.Split(ip, "")
	head := digs[0]
	digs = digs[1:]
	rv := float64(0)
	for i := 0; i < len(digs); i++ {
		d := digs[len(digs)-i-1]
		p := strings.Index(base62Digits, d)
		if p == -1 {
			return 0.0, fmt.Errorf("invalid order key: %s", key)
		}
		rv += math.Pow(float64(len(base62Digits)), float64(i)) * float64(p)
	}

	fp := key[len(ip):]
	for i, d := range fp {
		p := strings.Index(base62Digits, string(d))
		if p == -1 {
			return 0.0, fmt.Errorf("invalid key: %s", key)
		}
		rv += (float64(p) / math.Pow(float64(len(base62Digits)), float64(i+1)))
	}

	if head < "a" {
		rv *= -1
	}

	return rv, nil
}

// NKeysBetween returns n keys between a and b that sorts lexicographically.
// Either a or b can be empty strings. If a is empty it indicates smallest key,
// If b is empty it indicates largest key.
// b must be empty string or > a.
func NKeysBetween(a, b string, n uint) ([]string, error) {
	if n == 0 {
		return []string{}, nil
	}
	if n == 1 {
		c, err := KeyBetween(a, b)
		if err != nil {
			return nil, err
		}
		return []string{c}, nil
	}
	if b == "" {
		c, err := KeyBetween(a, b)
		if err != nil {
			return nil, err
		}
		result := make([]string, 0, n)
		result = append(result, c)
		for i := 0; i < int(n)-1; i++ {
			c, err = KeyBetween(c, b)
			if err != nil {
				return nil, err
			}
			result = append(result, c)
		}
		return result, nil
	}
	if a == "" {
		c, err := KeyBetween(a, b)
		if err != nil {
			return nil, err
		}
		result := make([]string, 0, n)
		result = append(result, c)
		for i := 0; i < int(n)-1; i++ {
			c, err = KeyBetween(a, c)
			if err != nil {
				return nil, err
			}
			result = append(result, c)
		}
		reverse(result)
		return result, nil
	}
	mid := n / 2
	c, err := KeyBetween(a, b)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, n)
	{
		r, err := NKeysBetween(a, c, mid)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	result = append(result, c)
	{
		r, err := NKeysBetween(c, b, n-mid-1)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	return result, nil
}

func reverse(values []string) {
	for i := 0; i < len(values)/2; i++ {
		j := len(values) - i - 1
		values[i], values[j] = values[j], values[i]
	}
}
