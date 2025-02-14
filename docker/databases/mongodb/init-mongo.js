// /mongodb/init-mongo.js
db = db.getSiblingDB('admin');

// Create root user if it doesn't exist
if (!db.getUser("root")) {
  db.createUser({
    user: "root",
    pwd: "example",  // Should be replaced with environment variable
    roles: [{ role: "root", db: "admin" }]
  });
}

// Initialize conversations database
db = db.getSiblingDB('my_database');

// Create collections with validation
db.createCollection("conversations", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["user_id", "start_time", "messages"],
      properties: {
        user_id: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        start_time: {
          bsonType: "date",
          description: "must be a date and is required"
        },
        end_time: {
          bsonType: ["date", "null"],
          description: "must be a date or null"
        },
        messages: {
          bsonType: "array",
          description: "must be an array and is required",
          items: {
            bsonType: "object",
            required: ["user_id", "sender", "text", "timestamp"],
            properties: {
              user_id: {
                bsonType: "string"
              },
              sender: {
                bsonType: "string"
              },
              text: {
                bsonType: "string"
              },
              timestamp: {
                bsonType: "date"
              }
            }
          }
        }
      }
    }
  }
});

// Create indexes
db.conversations.createIndex({ "user_id": 1 });
db.conversations.createIndex({ "start_time": -1 });
