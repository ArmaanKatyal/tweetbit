db.createUser(
  {
    user: "admin",
    pwd: "root",
    roles: [
      {
        role: "readWrite",
        db: "prod"
      }
    ]
  }
);
db.createCollection("users");
db.createCollection("auth");
