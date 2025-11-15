// apps/platform/lib/server/crypto.ts
import crypto from 'crypto';

export function sha256Hex(buf: Buffer): string {
  return crypto.createHash('sha256').update(buf).digest('hex');
}

export function isValidUlid(id: string): boolean {
  return /^[0-9A-HJKMNP-TV-Z]{26}$/.test(id);
}