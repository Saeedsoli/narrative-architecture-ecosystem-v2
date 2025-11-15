// scripts/db/mongo/mongo-init.js
// Narrative Architecture - MongoDB Collections, Validators, Indexes, TTL

(function () {
  const DB_NAME = 'narrative_arch_content';
  // در mongosh داخل کانتینر ممکن است db از قبل ست باشد:
  // از getSiblingDB برای سازگاری استفاده می‌کنیم.
  const target = (typeof db !== 'undefined' && db.getSiblingDB)
    ? db.getSiblingDB(DB_NAME)
    : connect(`127.0.0.1/${DB_NAME}`);

  // Helper: create collection with schema + (fa) collation, or update validator via collMod
  function ensureCollection(name, validator, useFaCollation = true) {
    const exists = target.getCollectionNames().includes(name);
    if (!exists) {
      const options = { validator };
      if (useFaCollation) options.collation = { locale: 'fa', strength: 1 };
      target.createCollection(name, options);
      print(`✓ Created collection ${name}`);
    } else {
      target.runCommand({ collMod: name, validator });
      print(`↻ Updated validator for ${name}`);
    }
  }

  // ULID regex (Crockford Base32 uppercase)
  const ULID_REGEX = '^[0-9A-HJKMNP-TV-Z]{26}$';

  // ========== 1) articles ==========
  const articlesSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id', 'content_group_id', 'locale', 'title', 'content', 'author', 'status', 'createdAt', 'updatedAt'],
      properties: {
        _id:               { bsonType: 'string', pattern: ULID_REGEX },
        content_group_id:  { bsonType: 'string', pattern: ULID_REGEX },
        locale:            { enum: ['fa', 'en'] },
        type:              { enum: ['article', 'podcast', 'video', null] },
        slug:              { bsonType: ['string','null'] }, // فقط en
        title:             { bsonType: 'string', minLength: 3 },
        excerpt:           { bsonType: ['string','null'] },
        content:           { bsonType: 'string' }, // Markdown یا AST stringify
        coverImage: {
          bsonType: ['object','null'],
          properties: {
            url:      { bsonType: 'string', pattern: '^https?://' },
            alt:      { bsonType: ['string','null'] },
            blurhash: { bsonType: ['string','null'] }
          }
        },
        author: {
          bsonType: 'object',
          required: ['id','name'],
          properties: {
            id:     { bsonType: 'string', pattern: ULID_REGEX },
            name:   { bsonType: 'string' },
            avatar: { bsonType: ['string','null'] }
          }
        },
        metadata: {
          bsonType: 'object',
          properties: {
            tags:      { bsonType: ['array'], items: { bsonType: 'string' } },
            category:  { bsonType: ['string','null'] },
            readTime:  { bsonType: ['int','long','null'], minimum: 0 },
            difficulty:{ enum: ['beginner','intermediate','advanced', null] }
          }
        },
        status:      { enum: ['draft','review','published','scheduled'] },
        publishedAt: { bsonType: ['date','null'] },
        createdAt:   { bsonType: 'date' },
        updatedAt:   { bsonType: 'date' },
        deletedAt:   { bsonType: ['date','null'] }
      }
    }
  };
  ensureCollection('articles', articlesSchema, true);
  target.articles.createIndex({ content_group_id: 1, locale: 1 }, { unique: true });
  target.articles.createIndex({ locale: 1, slug: 1 }, { unique: true, partialFilterExpression: { slug: { $type: 'string' } } });
  target.articles.createIndex({ locale: 1, status: 1, publishedAt: -1 });
  target.articles.createIndex({ 'metadata.tags': 1 });
  target.articles.createIndex({ 'author.id': 1 });
  target.articles.createIndex({ title: 'text', excerpt: 'text' }, { default_language: 'persian' });

  // ========== 2) article_comments ==========
  const articleCommentsSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','article_id','user','body','createdAt'],
      properties: {
        _id:        { bsonType: 'string', pattern: ULID_REGEX },
        article_id: { bsonType: 'string', pattern: ULID_REGEX },
        user: {
          bsonType: 'object',
          required: ['id','username'],
          properties: {
            id:       { bsonType: 'string', pattern: ULID_REGEX },
            username: { bsonType: 'string' },
            avatar:   { bsonType: ['string','null'] }
          }
        },
        parent_id:  { bsonType: ['string','null'], pattern: ULID_REGEX },
        thread_id:  { bsonType: ['string','null'], pattern: ULID_REGEX },
        body:       { bsonType: 'string', minLength: 1 },
        likes_count:    { bsonType: ['int','long'], minimum: 0 },
        dislikes_count: { bsonType: ['int','long'], minimum: 0 },
        moderation_flags: {
          bsonType: ['object','null'],
          properties: {
            isFlagged:   { bsonType: 'bool' },
            reason:      { bsonType: 'string' },
            moderatedAt: { bsonType: 'date' }
          }
        },
        createdAt:  { bsonType: 'date' },
        deletedAt:  { bsonType: ['date','null'] }
      }
    }
  };
  ensureCollection('article_comments', articleCommentsSchema, true);
  target.article_comments.createIndex({ article_id: 1, createdAt: -1 });
  target.article_comments.createIndex({ thread_id: 1, createdAt: 1 });
  target.article_comments.createIndex({ parent_id: 1 });
  target.article_comments.createIndex({ 'user.id': 1 });

  // ========== 3) comment_votes ==========
  const commentVotesSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','target_type','target_id','user_id','value','createdAt'],
      properties: {
        _id:         { bsonType: 'string', pattern: ULID_REGEX },
        target_type: { enum: ['article_comment','forum_post'] },
        target_id:   { bsonType: 'string', pattern: ULID_REGEX },
        user_id:     { bsonType: 'string', pattern: ULID_REGEX },
        value:       { enum: [1,-1] },
        createdAt:   { bsonType: 'date' }
      }
    }
  };
  ensureCollection('comment_votes', commentVotesSchema, true);
  target.comment_votes.createIndex({ target_id: 1, user_id: 1 }, { unique: true });
  target.comment_votes.createIndex({ createdAt: -1 });

  // ========== 4) forum_topics ==========
  const forumTopicsSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','locale','title','body','author','createdAt'],
      properties: {
        _id:      { bsonType: 'string', pattern: ULID_REGEX },
        locale:   { enum: ['fa','en'] },
        title:    { bsonType: 'string', minLength: 3 },
        body:     { bsonType: 'string' },
        tags:     { bsonType: 'array', items: { bsonType: 'string' } },
        author:   {
          bsonType: 'object',
          required: ['id','username'],
          properties: {
            id:       { bsonType: 'string', pattern: ULID_REGEX },
            username: { bsonType: 'string' }
          }
        },
        status:    { enum: ['open','locked', null] },
        createdAt: { bsonType: 'date' },
        updatedAt: { bsonType: ['date','null'] },
        deletedAt: { bsonType: ['date','null'] }
      }
    }
  };
  ensureCollection('forum_topics', forumTopicsSchema, true);
  target.forum_topics.createIndex({ locale: 1, updatedAt: -1 });
  target.forum_topics.createIndex({ tags: 1 });
  target.forum_topics.createIndex({ title: 'text', body: 'text' }, { default_language: 'persian' });

  // ========== 5) forum_posts ==========
  const forumPostsSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','topic_id','user','body','createdAt'],
      properties: {
        _id:        { bsonType: 'string', pattern: ULID_REGEX },
        topic_id:   { bsonType: 'string', pattern: ULID_REGEX },
        parent_id:  { bsonType: ['string','null'], pattern: ULID_REGEX },
        user: {
          bsonType: 'object',
          required: ['id','username'],
          properties: {
            id:       { bsonType: 'string', pattern: ULID_REGEX },
            username: { bsonType: 'string' }
          }
        },
        body:       { bsonType: 'string' },
        likes_count:    { bsonType: ['int','long'], minimum: 0 },
        dislikes_count: { bsonType: ['int','long'], minimum: 0 },
        flags:      { bsonType: ['object','null'] },
        createdAt:  { bsonType: 'date' },
        deletedAt:  { bsonType: ['date','null'] }
      }
    }
  };
  ensureCollection('forum_posts', forumPostsSchema, true);
  target.forum_posts.createIndex({ topic_id: 1, createdAt: -1 });
  target.forum_posts.createIndex({ parent_id: 1 });
  target.forum_posts.createIndex({ 'user.id': 1 });

  // ========== 6) chat_messages (TTL 90d) ==========
  const chatMessagesSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','room_id','user','type','createdAt'],
      properties: {
        _id:      { bsonType: 'string', pattern: ULID_REGEX },
        room_id:  { bsonType: 'string', pattern: ULID_REGEX },
        user: {
          bsonType: 'object',
          required: ['id','username'],
          properties: {
            id:       { bsonType: 'string', pattern: ULID_REGEX },
            username: { bsonType: 'string' }
          }
        },
        type:  { enum: ['text','voice'] },
        text:  { bsonType: ['string','null'] },
        voice: {
          bsonType: ['object','null'],
          properties: {
            url:      { bsonType: 'string', pattern: '^https?://' },
            duration: { bsonType: 'int', minimum: 1, maximum: 60 } // 1 min max
          }
        },
        createdAt: { bsonType: 'date' }
      }
    }
  };
  ensureCollection('chat_messages', chatMessagesSchema, false);
  target.chat_messages.createIndex({ createdAt: 1 }, { expireAfterSeconds: 7776000 });
  target.chat_messages.createIndex({ room_id: 1, createdAt: -1 });
  target.chat_messages.createIndex({ 'user.id': 1 });

  // ========== 7) media_uploads ==========
  const mediaUploadsSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','owner_id','type','url','size','sha256','createdAt'],
      properties: {
        _id:       { bsonType: 'string', pattern: ULID_REGEX },
        owner_id:  { bsonType: 'string', pattern: ULID_REGEX },
        type:      { enum: ['image','audio','pdf','video'] },
        url:       { bsonType: 'string', pattern: '^https?://' },
        size:      { bsonType: 'long', minimum: 0, maximum: 10485760 }, // <= 10MB
        sha256:    { bsonType: 'string', minLength: 64, maxLength: 64 },
        createdAt: { bsonType: 'date' }
      }
    }
  };
  ensureCollection('media_uploads', mediaUploadsSchema, false);
  target.media_uploads.createIndex({ owner_id: 1, createdAt: -1 });
  target.media_uploads.createIndex({ sha256: 1 }, { unique: true });

  // ========== 8) content_links ==========
  const contentLinksSchema = {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id','source','targets','createdAt'],
      properties: {
        _id:    { bsonType: 'string', pattern: ULID_REGEX },
        source: {
          bsonType: 'object',
          required: ['type','id'],
          properties: {
            type: { enum: ['article','video','podcast','exercise'] },
            id:   { bsonType: 'string' }
          }
        },
        targets: {
          bsonType: 'array',
          minItems: 1,
          items: {
            bsonType: 'object',
            required: ['type'],
            properties: {
              type:  { enum: ['article','video','podcast','exercise','resource'] },
              id:    { bsonType: ['string','null'] },
              url:   { bsonType: ['string','null'] },
              title: { bsonType: ['string','null'] }
            }
          }
        },
        createdAt: { bsonType: 'date' }
      }
    }
  };
  ensureCollection('content_links', contentLinksSchema, true);
  target.content_links.createIndex({ 'source.type': 1, 'source.id': 1, createdAt: -1 });

  print('✔ MongoDB: validators & indexes applied.');
})();