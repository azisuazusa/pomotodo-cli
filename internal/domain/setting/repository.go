package setting

import "context"

type SettingRepository interface {
	SetDropboxToken(ctx context.Context, token string) error
}
