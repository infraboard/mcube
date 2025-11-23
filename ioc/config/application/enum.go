package application

// 加密算法
type EncryptAlgorithm string

const (
	// AES-GCM
	ENCRYPT_ALGORITHM_AES_GCM EncryptAlgorithm = "AES-GCM"
	// AES-CBC
	ENCRYPT_ALGORITHM_AES_CBC EncryptAlgorithm = "AES-CBC"
)
