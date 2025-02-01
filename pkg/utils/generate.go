package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GetClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	remoteAddr := r.RemoteAddr
	if strings.Contains(remoteAddr, ":") {
		return strings.Split(remoteAddr, ":")[0]
	}

	return remoteAddr
}

func GenerateRandomId(folder string) string {
	return fmt.Sprintf("%s-%012d", folder, rand.Intn(1000000000000))
}

func GenerateSku(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	cleanName := reg.ReplaceAllString(name, "")
	cleanName = strings.ToUpper(cleanName)

	namePrefix := cleanName
	if len(namePrefix) > 3 {
		namePrefix = namePrefix[:3]
	} else {
		for len(namePrefix) < 3 {
			namePrefix += "X"
		}
	}
	year := time.Now().Year()
	randomNum := rand.Intn(9000) + 1000
	sku := fmt.Sprintf("%s-%d-%04d", namePrefix, year, randomNum)

	return sku
}

func IsSkuValid(sku string) bool {
	pattern := regexp.MustCompile(`^[A-Z]{3}-\d{4}-\d{4}$`)
	return pattern.MatchString(sku)
}
