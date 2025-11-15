import { Worker } from 'bullmq';
import { sync } from './sync';
import { config } from './config';
import pino from 'pino';

const logger = pino({ level: config.logLevel });

// آخرین زمان سینک را می‌توان در Redis یا یک فایل ذخیره کرد.
let lastSyncTime = new Date(Date.now() - config.sync.interval); // default
if (config.sync.initial) {
  lastSyncTime = new Date(0); // Sync from beginning
}

async function runSync() {
  logger.info('🚀 Starting search sync worker...');
  const startTime = Date.now();
  
  try {
    const count = await sync(lastSyncTime);
    logger.info(`✅ Sync finished. ${count} documents processed.`);
    lastSyncTime = new Date(startTime);
  } catch (err) {
    logger.error(err, '❌ Sync failed');
  }
}

// استفاده از BullMQ (Job Queue) برای مدیریت Cron
const worker = new Worker('sync-queue', async job => {
  if (job.name === 'sync-search') {
    await runSync();
  }
}, { connection: config.redis });

logger.info('Worker is listening for jobs...');