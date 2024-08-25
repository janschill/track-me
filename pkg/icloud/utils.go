package icloud

import "fmt"

func base62ToInt(char byte) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	} else if char >= 'A' && char <= 'Z' {
		return int(char - 'A' + 10)
	} else if char >= 'a' && char <= 'z' {
		return int(char - 'a' + 36)
	}
	return -1
}

func getPartitionFromToken(token string) string {
	var partition int
	if token[0] == 'A' {
		partition = base62ToInt(token[1])
	} else {
		partition = base62ToInt(token[1])*62 + base62ToInt(token[2])
	}

	if partition < 10 {
		return fmt.Sprintf("0%d", partition)
	}
	return fmt.Sprintf("%d", partition)
}
