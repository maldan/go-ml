package mfixed

import (
	"fmt"
	"strconv"
	"strings"
)

type FixedPoint struct {
	Integer    int
	Fractional int
}

func New(i int, f int) FixedPoint {
	return FixedPoint{Integer: i, Fractional: f}
}

func NewF(f float32) FixedPoint {
	integerPart := int(f)
	fractionalPart := int((f - float32(integerPart)) * 10000.0)
	return FixedPoint{Integer: integerPart, Fractional: fractionalPart}
}

func NewFromString(s string) FixedPoint {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return FixedPoint{}
	}

	integerPart, _ := strconv.Atoi(parts[0])
	fractionalPart, _ := strconv.Atoi(parts[1])

	// Умножаем дробную часть на multiplier
	fractionalPart /= 10000

	return FixedPoint{Integer: integerPart, Fractional: fractionalPart}
}

func (f FixedPoint) Add(v FixedPoint) FixedPoint {
	f.Integer = f.Integer + v.Integer
	f.Fractional = f.Fractional + v.Fractional

	if f.Fractional >= 10000 {
		f.Integer++
		f.Fractional -= 10000
	}

	return f
}

// Subtract выполняет вычитание двух чисел с фиксированной точностью
func (f FixedPoint) Sub(v FixedPoint) FixedPoint {
	result := FixedPoint{
		Integer:    f.Integer - v.Integer,
		Fractional: f.Fractional - v.Fractional,
	}

	if result.Fractional < 0 {
		result.Integer--
		result.Fractional += 10000
	}

	return result
}

func (f FixedPoint) Mul(v FixedPoint) FixedPoint {
	result := FixedPoint{
		Integer:    f.Integer*v.Integer + f.Fractional*v.Integer/10000 + f.Integer*v.Fractional/10000,
		Fractional: f.Fractional * v.Fractional / 10000,
	}

	result.Integer += result.Fractional / 10000
	result.Fractional %= 10000

	return result
}

func (f FixedPoint) Floor() FixedPoint {
	f.Fractional = 0
	return f
}

// ToString возвращает строковое представление числа
func (f FixedPoint) ToString() string {
	return fmt.Sprintf("%d.%04d", f.Integer, f.Fractional)
}
