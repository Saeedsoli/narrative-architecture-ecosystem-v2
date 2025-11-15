package entitlement

import "time"

type ResourceType string
const (
    ResourceTypeChapter ResourceType = "chapter"
    // ... سایر انواع
)

type Entitlement struct {
    ID           string
    UserID       string
    ResourceType ResourceType
    ResourceID   string // e.g., "16"
    Status       string // "active", "revoked"
    CreatedAt    time.Time
}