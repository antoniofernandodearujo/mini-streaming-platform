import { S3_BUCKET_NAME, REGION } from "@/config/aws/config";

export function getThumbnail(videoID: string) {
    return `https://${S3_BUCKET_NAME}.s3.${REGION}.amazonaws.com/thumbnails/${videoID}.jpg`;
};