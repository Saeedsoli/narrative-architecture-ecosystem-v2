export const config = {
  mongo: {
    uri: process.env.MONGODB_URI!,
    db: process.env.MONGODB_DB || 'narrative_arch_content',
  },
  elastic: {
    url: process.env.ELASTIC_URL!,
    username: process.env.ELASTIC_USERNAME,
    password: process.env.ELASTIC_PASSWORD,
  },
  redis: {
    host: process.env.REDIS_HOST || 'localhost',
    port: parseInt(process.env.REDIS_PORT || '6379', 10),
  },
  sync: {
    interval: (parseInt(process.env.SYNC_INTERVAL_MINUTES || '5', 10)) * 60 * 1000,
    batchSize: parseInt(process.env.SYNC_BATCH_SIZE || '100', 10),
    initial: process.env.SYNC_INITIAL === 'true',
  },
  logLevel: process.env.LOG_LEVEL || 'info',
};