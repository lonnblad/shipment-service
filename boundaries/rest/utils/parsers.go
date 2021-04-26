package utils

import (
	"fmt"
	"strconv"
)

// ParseLimit will take a limit as a string, a default and a max limit.
// It will return the limit as an integer or an error.
//
// If the provided limit is an empty string, ParseLimit will return the
// defaultLimit.
// If the provided limit can't be converted to an integer, is negative
// or above the provided max, ParseLimit will return an error.
func ParseLimit(limitStr string, defaultLimit, maxLimit int) (_ int, err error) {
	if limitStr == "" {
		return defaultLimit, nil
	}

	var limit int

	if limit, err = strconv.Atoi(limitStr); err != nil {
		err = fmt.Errorf("limit [%s] is not an integer: %w", limitStr, err)
		return
	}

	if limit < 0 {
		err = fmt.Errorf("limit must be a non-negative integer")
		return
	}

	if limit > maxLimit {
		err = fmt.Errorf("limit: %d is greater than max: %d", limit, maxLimit)
		return
	}

	return limit, nil
}

// ParseOffset will take an offset as a string and a default offset.
// It will return the offset as an integer or an error.
//
// If the provided offset is an empty string, ParseOffset will return the
// defaultOffset.
// If the provided offset can't be converted to an integer or is negative,
// ParseOffset will return an error.
func ParseOffset(offsetStr string, defaultOffset int) (_ int, err error) {
	if offsetStr == "" {
		return defaultOffset, nil
	}

	var offset int

	if offset, err = strconv.Atoi(offsetStr); err != nil {
		err = fmt.Errorf("offset [%s] is not a integer: %w", offsetStr, err)
		return
	}

	if offset < 0 {
		err = fmt.Errorf("offset must be a non-negative integer")
		return
	}

	return offset, nil
}
