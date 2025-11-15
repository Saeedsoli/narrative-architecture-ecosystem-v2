// apps/platform/app/api/voice-messages/upload/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { getMongoDb } from '@/lib/server/mongo';
import { putObject, buildPublicUrl } from '@/lib/server/s3';
import { sha256Hex, isValidUlid } from '@/lib/server/crypto';
import { ulid } from 'ulid';

export const runtime = 'nodejs'; // استفاده از Node APIs

const MAX_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_TYPES = new Set(['audio/webm', 'audio/webm;codecs=opus', 'audio/ogg']);

export async function POST(req: NextRequest) {
  try {
    // اختیاری: امنیت ساده برای internal upload (می‌توانید حذف کنید)
    const token = process.env.INTERNAL_UPLOAD_TOKEN;
    if (token) {
      const hdr = req.headers.get('x-internal-upload-token');
      if (hdr !== token) {
        return NextResponse.json({ error: 'unauthorized' }, { status: 401 });
      }
    }

    // دریافت FormData
    const form = await req.formData();
    const file = form.get('audio') as File | null;
    const durationStr = form.get('duration')?.toString();
    // مالک: اگر Auth هنوز وصل نیست، از ULID صفر به‌عنوان ناشناس استفاده می‌کنیم
    let ownerId = (form.get('ownerId')?.toString() || '').trim();
    if (!ownerId || !isValidUlid(ownerId)) {
      ownerId = '00000000000000000000000000'; // ULID 26 کاراکتر صفر (مجاز توسط regex)
    }

    if (!file) {
      return NextResponse.json({ error: 'audio file required' }, { status: 400 });
    }

    const contentType = file.type || 'application/octet-stream';
    if (!ALLOWED_TYPES.has(contentType)) {
      return NextResponse.json({ error: 'unsupported content-type' }, { status: 415 });
    }

    const arrayBuf = await file.arrayBuffer();
    const buf = Buffer.from(arrayBuf);
    if (buf.byteLength === 0 || buf.byteLength > MAX_SIZE) {
      return NextResponse.json({ error: 'file too large (max 10MB)' }, { status: 413 });
    }

    // مدت (ثانیه) - از کلاینت دریافت می‌شود
    const duration = durationStr ? parseInt(durationStr, 10) : NaN;
    if (!Number.isFinite(duration) || duration < 1 || duration > 60) {
      return NextResponse.json({ error: 'invalid duration (1..60 seconds)' }, { status: 400 });
    }

    // sha256 برای Dedup
    const sha256 = sha256Hex(buf);

    const db = await getMongoDb();
    const coll = db.collection('media_uploads');

    // اگر قبلاً همین فایل آپلود شده، همان URL را برگردان
    const existing = await coll.findOne({ sha256 });
    if (existing) {
      return NextResponse.json({
        id: existing._id,
        url: existing.url,
        duration,
        sha256,
        dedup: true,
      });
    }

    // تولید id و کلید S3
    const id = ulid().toUpperCase();
    const key = `voice/${id}.webm`;

    // آپلود به S3/R2
    await putObject(key, buf, 'audio/webm');

    const url = buildPublicUrl(key);

    // ذخیره متادیتا در media_uploads (مطابق schema)
    await coll.insertOne({
      _id: id,
      owner_id: ownerId,
      type: 'audio',
      url,
      size: buf.byteLength,
      sha256,
      createdAt: new Date(),
    });

    return NextResponse.json({
      id,
      url,
      duration,
      sha256,
      dedup: false,
    });
  } catch (err: any) {
    console.error('voice upload error:', err);
    return NextResponse.json({ error: 'internal_error' }, { status: 500 });
  }
}