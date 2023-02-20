-- CreateTable
CREATE TABLE "User" (
    "id" INT8 NOT NULL DEFAULT unique_rowid(),
    "uuid" STRING NOT NULL,
    "email" STRING NOT NULL,
    "first_name" STRING NOT NULL,
    "last_name" STRING NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "favorites_count" INT4 NOT NULL DEFAULT 0,
    "followers_count" INT4 NOT NULL DEFAULT 0,
    "following_count" INT4 NOT NULL DEFAULT 0,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Tweet" (
    "id" INT8 NOT NULL DEFAULT unique_rowid(),
    "uuid" STRING NOT NULL,
    "user_id" INT8 NOT NULL,
    "content" STRING NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "likes_count" INT4 NOT NULL DEFAULT 0,
    "retweets_count" INT4 NOT NULL DEFAULT 0,

    CONSTRAINT "Tweet_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "User_Followers" (
    "id" INT8 NOT NULL DEFAULT unique_rowid(),
    "user_id" INT8 NOT NULL,
    "follower_id" INT8 NOT NULL,

    CONSTRAINT "User_Followers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Tweet_Likes" (
    "id" INT8 NOT NULL DEFAULT unique_rowid(),
    "tweet_id" INT8 NOT NULL,
    "user_id" INT8 NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Tweet_Likes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Tweet_Comments" (
    "id" INT8 NOT NULL DEFAULT unique_rowid(),
    "tweet_id" INT8 NOT NULL,
    "user_id" INT8 NOT NULL,
    "content" STRING NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Tweet_Comments_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_uuid_key" ON "User"("uuid");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_uuid_key" ON "Tweet"("uuid");

-- AddForeignKey
ALTER TABLE "Tweet" ADD CONSTRAINT "Tweet_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "User_Followers" ADD CONSTRAINT "User_Followers_follower_id_fkey" FOREIGN KEY ("follower_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tweet_Likes" ADD CONSTRAINT "Tweet_Likes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tweet_Comments" ADD CONSTRAINT "Tweet_Comments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
