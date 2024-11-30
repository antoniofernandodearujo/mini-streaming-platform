import { S3Client } from "@aws-sdk/client-s3";

export const S3_BUCKET_NAME = process.env.NEXT_PUBLIC_S3_BUCKET_NAME || "video-upload-demo";
export const REGION = process.env.NEXT_PUBLIC_AWS_REGION; // Altere para a sua região
export const MAX_FILE_SIZE = 500 * 1024 * 1024; // 500 MB
export const ACCEPTED_TYPES = ["video/mp4", "video/avi", "video/mov"];

// Configure o cliente do S3
const accessKeyId = process.env.NEXT_PUBLIC_AWS_ACCESS_KEY_ID;
const secretAccessKey = process.env.NEXT_PUBLIC_AWS_SECRET_ACCESS_KEY;

// Verifica se as credenciais estão definidas
if (!accessKeyId || !secretAccessKey) {
    throw new Error("AWS credentials are not defined");
}

export const s3Client = new S3Client({
    region: REGION,
    credentials: {
      accessKeyId,
      secretAccessKey
    },
});