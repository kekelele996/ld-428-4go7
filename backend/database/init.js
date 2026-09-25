// artvault MongoDB 初始化脚本（参考用）
// 实际集合与种子数据由后端启动时自动写入（database/seeds/seed.go）。
db = db.getSiblingDB('artvault');
db.createCollection('users');
db.createCollection('artists');
db.createCollection('artworks');
db.createCollection('exhibitions');
db.createCollection('interactions');
db.createCollection('review_logs');
db.createCollection('audit_logs');
