package error

import "fmt"

type ApiErrorKind int

const (
	NoTokenProvided ApiErrorKind = iota
	InvalidTokenPrefix
	InvalidApiKey
	MerchantNotFound
	NoTestSshKeyFound
	NoLiveSshKeyFound
	IpAddressNotWhitelisted
	CouldNotDetermineClientIpAddress
	PermissionDenied
	MissingEncryptionDataKeyOrIv
	DecryptionFailed
	Validation
	TooManyRequests
	ServiceUnavailable
	BusinessLogic
	Unknown
)

func (k ApiErrorKind) String() string {
	switch k {
	case NoTokenProvided:
		return "no token provided"
	case InvalidTokenPrefix:
		return "invalid token prefix"
	case InvalidApiKey:
		return "invalid API key"
	case MerchantNotFound:
		return "merchant not found"
	case NoTestSshKeyFound:
		return "no test SSH key found"
	case NoLiveSshKeyFound:
		return "no live SSH key found"
	case IpAddressNotWhitelisted:
		return "IP address not whitelisted"
	case CouldNotDetermineClientIpAddress:
		return "could not determine client IP address"
	case PermissionDenied:
		return "permission denied"
	case MissingEncryptionDataKeyOrIv:
		return "missing encryption data, key, or IV"
	case DecryptionFailed:
		return "decryption failed"
	case Validation:
		return "validation failed"
	case TooManyRequests:
		return "too many requests; try again later"
	case ServiceUnavailable:
		return "service is currently unavailable; try again later"
	case BusinessLogic:
		return "business logic error"
	case Unknown:
		return "unknown API error"
	default:
		return "unknown API error"
	}
}

type ApiError struct {
	Status   uint16
	Kind     ApiErrorKind
	Field    string
	Message  string
}

func (e *ApiError) Error() string {
	if e.Kind == Validation {
		return fmt.Sprintf("validation failed: %s: %s", e.Field, e.Message)
	}
	if e.Kind == BusinessLogic {
		return fmt.Sprintf("business logic error: %s", e.Message)
	}
	return e.Kind.String()
}

type ErrorType int

const (
	TypeConfiguration ErrorType = iota
	TypeNetwork
	TypeApi
	TypeParse
	TypeCrypto
	TypeInvalidEncryptedPayload
)

type Error struct {
	Kind    ErrorType
	Message string
	Status  uint16
	ApiKind ApiErrorKind
	Field   string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) IsRetryable() bool {
	if e.Kind == TypeNetwork {
		return true
	}
	if e.Kind == TypeApi {
		return e.Status >= 500 || e.ApiKind == TooManyRequests || e.ApiKind == ServiceUnavailable
	}
	return false
}

func NewConfigurationError(msg string) *Error {
	return &Error{
		Kind:    TypeConfiguration,
		Message: fmt.Sprintf("configuration error: %s", msg),
	}
}

func NewNetworkError(err error) *Error {
	return &Error{
		Kind:    TypeNetwork,
		Message: fmt.Sprintf("request failed: %v", err),
		Err:     err,
	}
}

func NewApiError(status uint16, kind ApiErrorKind, field, message string) *Error {
	msg := kind.String()
	if kind == Validation {
		msg = fmt.Sprintf("validation failed: %s: %s", field, message)
	} else if kind == BusinessLogic {
		msg = fmt.Sprintf("business logic error: %s", message)
	} else if kind == Unknown {
		msg = fmt.Sprintf("unknown API error: %s", message)
	}
	return &Error{
		Kind:    TypeApi,
		Status:  status,
		ApiKind: kind,
		Field:   field,
		Message: fmt.Sprintf("cNGN API error (%d): %s", status, msg),
	}
}

func NewParseError(err error) *Error {
	return &Error{
		Kind:    TypeParse,
		Message: fmt.Sprintf("invalid API response: %v", err),
		Err:     err,
	}
}

func NewCryptoError(msg string) *Error {
	return &Error{
		Kind:    TypeCrypto,
		Message: fmt.Sprintf("cryptography error: %s", msg),
	}
}

func NewInvalidEncryptedPayloadError(msg string) *Error {
	return &Error{
		Kind:    TypeInvalidEncryptedPayload,
		Message: fmt.Sprintf("invalid encrypted payload: %s", msg),
	}
}

func ClassifyApiError(status uint16, message string) (ApiErrorKind, string, string) {
	switch message {
	case "No token provided":
		return NoTokenProvided, "", ""
	case "Invalid token prefix":
		return InvalidTokenPrefix, "", ""
	case "Invalid api key":
		return InvalidApiKey, "", ""
	case "Merchant not found":
		return MerchantNotFound, "", ""
	case "No Test SSH Key found":
		return NoTestSshKeyFound, "", ""
	case "No Live SSH Key found":
		return NoLiveSshKeyFound, "", ""
	case "IP address not whitelisted":
		return IpAddressNotWhitelisted, "", ""
	case "Could not determine client IP address":
		return CouldNotDetermineClientIpAddress, "", ""
	case "Permission denied":
		return PermissionDenied, "", ""
	case "Missing encryption data, key, or IV":
		return MissingEncryptionDataKeyOrIv, "", ""
	case "Decryption failed":
		return DecryptionFailed, "", ""
	case "Too many requests. Please try again later.":
		return TooManyRequests, "", ""
	case "Service is currently unavailable. Please try again later.":
		return ServiceUnavailable, "", ""
	default:
		if status == 400 {
			return classifyBadRequest(message)
		}
		return Unknown, "", message
	}
}

func classifyBadRequest(message string) (ApiErrorKind, string, string) {
	for i := 0; i < len(message); i++ {
		if message[i] == ':' {
			field := message[:i]
			msg := message[i+1:]
			if len(field) > 0 && len(msg) > 0 {
				return Validation, field, msg
			}
		}
	}
	return BusinessLogic, "", message
}
