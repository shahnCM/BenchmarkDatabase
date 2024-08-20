package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomString generates a random string of a specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomInt generates a random integer within a specified range
func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomDouble generates a random float64 within a specified range
func RandomDouble(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomBool generates a random boolean value
func RandomBool() bool {
	return rand.Intn(2) == 0
}

// RandomizeJSON randomizes the JSON structure
func RandomizeJSON() map[string]interface{} {
	jsonData := map[string]interface{}{
		"logged_at": time.Now().UTC().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"cpu": map[string]interface{}{
				"cpu_id":      RandomString(6),
				"temperature": fmt.Sprintf("%.2f", RandomDouble(30.0, 100.0)),
				"usage": map[string]interface{}{
					"user":   fmt.Sprintf("%.2f", RandomDouble(0.0, 100.0)),
					"system": fmt.Sprintf("%.2f", RandomDouble(0.0, 100.0)),
					"idle":   fmt.Sprintf("%.2f", RandomDouble(0.0, 100.0)),
				},
			},
			"user": map[string]interface{}{
				"user_id":  RandomString(6),
				"username": RandomString(8),
				"role":     RandomString(6),
				"preferences": map[string]interface{}{
					"theme":         []string{"dark", "light"}[rand.Intn(2)],
					"notifications": RandomBool(),
				},
			},
			"device": map[string]interface{}{
				"device_id": RandomString(8),
				"type":      RandomString(6),
				"location": map[string]interface{}{
					"latitude":  RandomDouble(-90.0, 90.0),
					"longitude": RandomDouble(-180.0, 180.0),
				},
			},
			"network": map[string]interface{}{
				"ip_address":      fmt.Sprintf("%d.%d.%d.%d", RandomInt(0, 255), RandomInt(0, 255), RandomInt(0, 255), RandomInt(0, 255)),
				"mac_address":     RandomString(12),
				"signal_strength": fmt.Sprintf("%.2f", RandomDouble(0.0, 100.0)),
			},
			"application": map[string]interface{}{
				"app_id":  RandomString(8),
				"version": fmt.Sprintf("v%d.%d", RandomInt(1, 10), RandomInt(0, 9)),
				"status":  []string{"active", "inactive"}[rand.Intn(2)],
			},
			"additional_info": map[string]interface{}{
				"field_1": RandomString(256),
				"field_2": RandomString(256),
				"field_3": RandomString(256),
				"field_4": RandomString(256),
			},
		},
	}
	return jsonData
}
