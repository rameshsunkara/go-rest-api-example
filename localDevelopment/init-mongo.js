// Production-ready MongoDB initialization script
// This runs automatically when MongoDB container starts for the first time

print('🚀 Initializing ecommerce database...');

// Explicitly switch to ecommerce database
db = db.getSiblingDB('ecommerce');

print('� Current database: ' + db.getName());

// Create application user with minimal required permissions
db.createUser({
  user: 'ecommerce_service',
  pwd: 'ecommerce_secure_password',
  roles: [
    {
      role: 'readWrite',
      db: 'ecommerce'  // Only access to ecommerce database
    }
  ]
});

print('✅ Created application user: ecommerce_service');

// Create useful indexes for the application
db.purchaseOrders.createIndex({ "user": 1 }, { background: true });
db.purchaseOrders.createIndex({ "createdAt": -1 }, { background: true });
db.purchaseOrders.createIndex({ "status": 1 }, { background: true });

print('✅ Created performance indexes');

// Insert a test document to verify everything works
db.purchaseOrders.insertOne({
  _id: ObjectId(),
  user: "system@test.com",
  status: "SystemTest",
  totalAmount: 0,
  products: [],
  createdAt: new Date(),
  updatedAt: new Date(),
  version: 1
});

print('✅ Inserted test document in database: ' + db.getName());
print('🎉 Database initialization completed successfully!');