package reward

import (
	"crypto/rand"
	"fmt"
)

const promoCodeRandomLength = 10
const promoCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generatePromoCode(prefix string) (string, error) {
	randomPart := make([]byte, promoCodeRandomLength)
	buffer := make([]byte, promoCodeRandomLength)

	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate promo code: %w", err)
	}

	for i := range randomPart {
		randomPart[i] = promoCodeAlphabet[int(buffer[i])%len(promoCodeAlphabet)]
	}

	return prefix + "-" + string(randomPart), nil
}
