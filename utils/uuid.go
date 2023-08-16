package utils

import (
	"github.com/google/uuid"
	"strings"
)

func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
