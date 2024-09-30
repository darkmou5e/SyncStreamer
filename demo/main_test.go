package main

import (
	"reflect"
	"testing"
)

func TestMatrixBuilder(t *testing.T) {
	textSample := "??"
	textMatrixSample := []int{
		//             X
		1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0,
	}

	matrix, width := textToLedMatrix(textSample)

	if !reflect.DeepEqual(matrix, textMatrixSample) {
		t.Errorf("\n%v\n!=\n%v\n", matrix, textMatrixSample)
	}

	if width != 11 {
		t.Errorf("\n%v\n!=\n%v\n", width, 11)
	}
}

func TestMatrixRepeat(t *testing.T) {
	srcMatrixWidth := 11
	srcMatrixOffset := 3
	srcMatrix := []int{
		//             X
		1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1,
		1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1,
		0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0,
	}

	destMatrixWidth := 22
	destMatrix := []int{
		//                               V
		1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0,
		1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1,
		0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1,
	}

	resultMatrix := repeatMatrix(srcMatrix, srcMatrixWidth, srcMatrixOffset, destMatrixWidth)

	if !reflect.DeepEqual(resultMatrix, destMatrix) {
		t.Errorf("\n%v\n!=\n%v\n", resultMatrix, destMatrix)
	}
}
