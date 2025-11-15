export interface MongoDoc {
  _id: string;
  locale: 'fa' | 'en';
  title?: string;
  excerpt?: string;
  content?: string;
  slug?: string | null;
  content_group_id?: string;
  tags?: string[];
  category?: string;
  publishedAt?: Date | null;
  author?: {
    id: string;
    name: string;
  };
}

export interface EsDoc {
  id: string;
  locale: string;
  slug?: string;
  title: string;
  excerpt?: string;
  content?: string;
  content_group_id?: string;
  tags?: string[];
  category?: string;
  publishedAt?: string;
  author?: {
    id: string;
    name: string;
  };
}