"use client";

import { useState, useRef } from "react";
import { ButtonNavigate, SubmitButton } from "@/components/Button";
import { CreateMultipartUploadCommand, UploadPartCommand, CompleteMultipartUploadCommand } from "@aws-sdk/client-s3";
import { S3_BUCKET_NAME, s3Client } from "@/config/aws/config";
import { createChunks, validateFile } from "@/functions/uploadForm";

export default function VideoUploadForm() {
  const [message, setMessage] = useState<string>("");
  const [isUploading, setIsUploading] = useState<boolean>(false);
  const [progress, setProgress] = useState<number>(0); // Estado para o progresso
  const formRef = useRef<HTMLFormElement>(null);

  // Upload do vídeo em chunks para o S3
  async function uploadVideoToS3(file: File) {
    const fileName = `videos/${file.name}`;
    const chunks = createChunks(file);

    // Inicia o upload multipart
    const createCommand = new CreateMultipartUploadCommand({
      Bucket: S3_BUCKET_NAME,
      Key: fileName,
    });

    const { UploadId } = await s3Client.send(createCommand);

    if (!UploadId) {
      throw new Error("Falha ao iniciar o upload multipart.");
    }

    const parts: { ETag: string; PartNumber: number }[] = [];

    for (let i = 0; i < chunks.length; i++) {
      setMessage(`Enviando chunk ${i + 1} de ${chunks.length}...`);

      const uploadPartCommand = new UploadPartCommand({
        Bucket: S3_BUCKET_NAME,
        Key: fileName,
        PartNumber: i + 1,
        UploadId,
        Body: chunks[i],
      });

      const { ETag } = await s3Client.send(uploadPartCommand);

      if (!ETag) {
        throw new Error(`Falha ao enviar o chunk ${i + 1}.`);
      }

      parts.push({ ETag, PartNumber: i + 1 });

      // Atualiza o progresso
      const progressValue = Math.round(((i + 1) / chunks.length) * 100);
      setProgress(progressValue);
    }

    // Finaliza o upload multipart
    const completeCommand = new CompleteMultipartUploadCommand({
      Bucket: S3_BUCKET_NAME,
      Key: fileName,
      UploadId,
      MultipartUpload: { Parts: parts },
    });

    await s3Client.send(completeCommand);

    return "Upload concluído com sucesso!";
  }

  // Função para tratar o envio do vídeo
  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const fileInput = event.currentTarget.querySelector<HTMLInputElement>("#video-upload");
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
      setMessage("Por favor, selecione um arquivo para enviar.");
      return;
    }

    const file = fileInput.files[0];

    try {
      validateFile(file); // Validações de tipo e tamanho

      setIsUploading(true);
      setProgress(0); // Reseta o progresso antes de começar
      const result = await uploadVideoToS3(file);
      setMessage(result);
    } catch (error: any) {
      setMessage(error.message || "Erro no upload.");
    } finally {
      setIsUploading(false);
      formRef.current?.reset();
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="w-full max-w-md p-8 space-y-8 bg-gray-900 rounded-xl shadow-2xl">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold text-white">Envie Seu Vídeo</h2>
          <p className="mt-2 text-sm text-gray-400">Compartilhe seu conteúdo incrível com o mundo!</p>
        </div>
        <form ref={formRef} onSubmit={handleSubmit} className="mt-8 space-y-6">
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="video-upload" className="sr-only">
                Escolha o arquivo de vídeo
              </label>
              <input
                id="video-upload"
                name="video"
                type="file"
                accept="video/*"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-700 placeholder-gray-500 text-gray-300 rounded-t-md focus:outline-none focus:ring-gray-500 focus:border-gray-500 focus:z-10 sm:text-sm bg-gray-800"
              />
            </div>
          </div>

          <div>
            <SubmitButton isUploading={isUploading} />
          </div>

            <div className="text-center text-sm font-bold hover:text-white transition-all duration-150 text-gray-400">
              <ButtonNavigate isUploading={isUploading} to="/videolist">Ir para a lista de vídeos</ButtonNavigate>
            </div>

        </form>
        {isUploading && (
          <div className="mt-4">
            <div className="relative pt-1">
              <div className="flex mb-2 items-center justify-between">
                <div>
                  <span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-white bg-gray-600">
                    Progresso do Upload
                  </span>
                </div>
                <div className="text-right">
                  <span className="text-xs font-semibold inline-block text-gray-500">
                    {progress}%
                  </span>
                </div>
              </div>
              <div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-gray-700">
                <div
                  style={{ width: `${progress}%` }}
                  className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-green-500"
                ></div>
              </div>
            </div>
          </div>
        )}
        {message && (
          <p className={`mt-2 text-sm ${message.includes("sucesso") ? "text-green-400" : "text-red-400"}`}>
            {message}
          </p>
        )}
      </div>
    </div>
  );
}
