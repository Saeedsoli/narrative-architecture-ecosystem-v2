// apps/platform/lib/server/mongo.ts
import { MongoClient, Db } from 'mongodb';

let client: MongoClient | null = null;
let db: Db | null = null;

export async function getMongoDb(): Promise<Db> {
  if (db) return db;

  const uri = process.env.MONGODB_URI!;
  const dbName = process.env.MONGODB_DB || 'narrative_arch_content';
  if (!uri) throw new Error('MONGODB_URI is not set');

  client = new MongoClient(uri, {
    maxPoolSize: 10,
    retryWrites: true,
  });
  await client.connect();
  db = client.db(dbName);
  return db!;
}