import { MongoClient, Db } from 'mongodb';
import { config } from './config';

let client: MongoClient | null = null;
let db: Db | null = null;

export async function getMongoDb(): Promise<Db> {
  if (db) return db;
  client = new MongoClient(config.mongo.uri);
  await client.connect();
  db = client.db(config.mongo.db);
  return db!;
}