// apps/platform/lib/server/s3.ts
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';

const endpoint = process.env.S3_ENDPOINT;
const region = process.env.S3_REGION || 'auto';
const accessKeyId = process.env.S3_ACCESS_KEY_ID!;
const secretAccessKey = process.env.S3_SECRET_ACCESS_KEY!;
const bucket = process.env.S3_BUCKET!;
const publicBase = process.env.S3_PUBLIC_BASE_URL; // اختیاری

if (!bucket) throw new Error('S3_BUCKET is not set');
if (!accessKeyId || !secretAccessKey) {
  throw new Error('S3 credentials are not set');
}

export const s3 = new S3Client({
  region,
  endpoint,
  forcePathStyle: !!endpoint, // برای MinIO/R2
  credentials: {
    accessKeyId,
    secretAccessKey,
  },
});

export async function putObject(key: string, body: Buffer, contentType: string) {
  await s3.send(new PutObjectCommand({
    Bucket: bucket,
    Key: key,
    Body: body,
    ContentType: contentType,
    ACL: 'public-read', // اگر R2/AWS ACL را می‌پذیرند؛ در صورت استفاده از signed URL می‌توانید حذف کنید
  }));
}

export function buildPublicUrl(key: string): string {
  if (publicBase) {
    // CDN baseUrl + key
    return `${publicBase.replace(/\/+$/, '')}/${key}`;
  }
  if (endpoint) {
    // endpoint/bucket/key
    return `${endpoint.replace(/\/+$/, '')}/${bucket}/${key}`;
  }
  // AWS S3 region-based
  return `https://${bucket}.s3.${region}.amazonaws.com/${key}`;
}