import { Client } from '@elastic/elasticsearch';
import { config } from './config';

export const es = new Client({
  node: config.elastic.url,
  auth: (config.elastic.username && config.elastic.password)
    ? { username: config.elastic.username, password: config.elastic.password }
    : undefined,
});