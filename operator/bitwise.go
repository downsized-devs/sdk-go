package operator

/*
 1. Left shift given number 1 by position-1 to create a number that has only set bit as position-th bit.
    temp = 1 << (position-1)
 2. If bitwise AND of number and temp is non-zero, then result is SET else result is NOT SET.
*/
func CheckBitOnPosition(number, position int) bool {
	return number&(1<<(position-1)) != 0
}
