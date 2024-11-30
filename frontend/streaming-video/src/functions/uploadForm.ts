import {
    ACCEPTED_TYPES,
    MAX_FILE_SIZE
} from "@/config/aws/config";

// Validações de tipo e tamanho
export function validateFile(file: File) {
  if (!ACCEPTED_TYPES.includes(file.type)) {
    throw new Error("Tipo de arquivo não suportado. Envie apenas vídeos.");
  }

  if (file.size > MAX_FILE_SIZE) {
    throw new Error("Arquivo excede o tamanho máximo permitido de 500 MB.");
  }
}

// Divide o arquivo em chunks
export function createChunks(file: File, chunkSize: number = 5 * 1024 * 1024) {
  const chunks = [];
  let offset = 0;

  while (offset < file.size) {
      chunks.push(file.slice(offset, offset + chunkSize));
      offset += chunkSize;
  }

  return chunks;
}

  