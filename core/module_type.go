package core

type ModuleType string

const (
	ModuleTypeInput      ModuleType = "input"
	ModuleTypeOutput     ModuleType = "output"
	ModuleTypeCompress   ModuleType = "compress"
	ModuleTypeDecompress ModuleType = "decompress"
	ModuleTypeEncrypt    ModuleType = "encrypt"
	ModuleTypeDecrypt    ModuleType = "decrypt"
)
