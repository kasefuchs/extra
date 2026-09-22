package link

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUnsupportedProtocol     = errors.New("unsupported protocol")
	ErrInvalidProxySettings    = errors.New("invalid proxy settings")
	ErrInvalidReceiverSettings = errors.New("invalid receiver settings")
	ErrInvalidSettingsType     = errors.New("invalid settings type")
	ErrSecuritySettings        = errors.New("failed to build security settings")
	ErrTransportSettings       = errors.New("failed to build transport settings")
)
