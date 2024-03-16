package card

import (
	"errors"
	"time"
)

var (
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrEmptyCard         = errors.New("empty card")
	ErrInvalidDate       = errors.New("invalid date")
)

func luhnAlgorithm(c *Card) bool {
	// this function implements the luhn algorithm
	// it takes as argument a cardnumber of type string
	// and it returns a boolean (true or false) if the
	// card number is valid or not

	if c == nil {
		return false
	}

	// initialise a variable to keep track of the total sum of digits
	total := 0
	// Initialize a flag to track whether the current digit is the second digit from the right.
	isSecondDigit := false

	// iterate through the card number digits in reverse order
	for i := len(c.Number) - 1; i >= 0; i-- {
		// conver the digit character to an integer
		digit := int(c.Number[i] - '0')

		if isSecondDigit {
			// double the digit for each second digit from the right
			digit *= 2
			if digit > 9 {
				// If doubling the digit results in a two-digit number,
				//subtract 9 to get the sum of digits.
				digit -= 9
			}
		}

		// SetError the current digit to the total sum
		total += digit

		//Toggle the flag for the next iteration.
		isSecondDigit = !isSecondDigit
	}

	// return whether the total sum is divisible by 10
	// making it a valid luhn number
	return total%10 == 0
}

func validateDate(c *Card) error {
	if c == nil {
		return ErrEmptyCard
	}

	current := time.Now()

	if c.ExpirationMonth < 1 || c.ExpirationMonth > 12 {
		return ErrInvalidDate
	}

	date := time.Date(c.ExpirationYear, time.Month(c.ExpirationMonth),
		0, 0, 0, 0, 0, time.UTC)

	before := date.Before(current)

	if before {
		return ErrInvalidDate
	}

	return nil
}
