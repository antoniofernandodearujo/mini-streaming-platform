package utils

import (
	"log"
	"os"
)

// EnsureDirectoryExists verifica e cria o diretório, se necessário
func EnsureDirectoryExists(path string) {
	log.Printf("Verificando/criando diretório: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Erro ao criar diretório '%s': %v", path, err)
		}
		log.Printf("Diretório criado: %s", path)
	} else {
		log.Printf("O diretório '%s' já existe", path)
	}
}
