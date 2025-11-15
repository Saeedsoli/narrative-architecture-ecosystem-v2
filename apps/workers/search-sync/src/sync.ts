import { getMongoDb } from './mongo';
import { es } from './elasticsearch';
import { config } from './config';
import type { MongoDoc, EsDoc } from './types';
import pino from 'pino';

const logger = pino({ level: config.logLevel });

const collectionsToSync = [
  { mongo: 'articles', es_fa: 'articles_fa', es_en: 'articles_en' },
  { mongo: 'forum_topics', es_fa: 'forum_fa', es_en: 'forum_en' },
];

function transformDoc(doc: MongoDoc): EsDoc {
  return {
    id: doc._id,
    locale: doc.locale,
    slug: doc.slug,
    title: doc.title || '',
    excerpt: doc.excerpt || '',
    content: doc.content || '',
    content_group_id: doc.content_group_id,
    tags: doc.metadata?.tags || doc.tags,
    category: doc.metadata?.category,
    publishedAt: doc.publishedAt?.toISOString(),
    author: doc.author,
  };
}

export async function sync(lastSyncTime: Date) {
  const db = await getMongoDb();
  let totalDocsSynced = 0;

  for (const collection of collectionsToSync) {
    logger.info(`🔄 Syncing collection: ${collection.mongo}`);

    const cursor = db.collection(collection.mongo).find({
      updatedAt: { $gt: lastSyncTime },
    });

    let bulkOps: any[] = [];
    let count = 0;

    for await (const doc of cursor) {
      count++;
      const esDoc = transformDoc(doc as MongoDoc);
      const esIndex = doc.locale === 'fa' ? collection.es_fa : collection.es_en;

      if (!esIndex) {
        logger.warn(`No ES index for locale ${doc.locale} in ${collection.mongo}`);
        continue;
      }

      bulkOps.push({ index: { _index: esIndex, _id: doc._id } });
      bulkOps.push(esDoc);

      if (bulkOps.length >= config.sync.batchSize * 2) {
        await flushBulk(bulkOps);
        totalDocsSynced += bulkOps.length / 2;
        bulkOps = [];
      }
    }

    if (bulkOps.length > 0) {
      await flushBulk(bulkOps);
      totalDocsSynced += bulkOps.length / 2;
    }
    logger.info(`✅ Synced ${count} docs from ${collection.mongo}`);
  }

  return totalDocsSynced;
}

async function flushBulk(bulkOps: any[]) {
  if (bulkOps.length === 0) return;
  
  logger.info(`Flushing ${bulkOps.length / 2} docs to Elasticsearch...`);
  
  const { body: bulkResponse } = await es.bulk({
    refresh: true,
    body: bulkOps,
  });

  if (bulkResponse.errors) {
    const erroredDocs = bulkResponse.items.filter((item: any) => item.index?.error);
    logger.error('Bulk errors:', JSON.stringify(erroredDocs, null, 2));
  }
}