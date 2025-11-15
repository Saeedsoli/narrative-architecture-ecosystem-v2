package entitlement

import "context"

type Repository interface {
    HasAccess(ctx context.Context, userID string, resourceType ResourceType, resourceID string) (bool, error)
}