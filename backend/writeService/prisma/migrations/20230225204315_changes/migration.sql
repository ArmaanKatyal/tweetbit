/*
  Warnings:

  - A unique constraint covering the columns `[uuid]` on the table `Tweet` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `Tweet` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[tweet_id]` on the table `Tweet_Comments` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `Tweet_Comments` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `User_Followers` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[follower_id]` on the table `User_Followers` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Tweet_uuid_key" ON "Tweet"("uuid");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_user_id_key" ON "Tweet"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Comments_tweet_id_key" ON "Tweet_Comments"("tweet_id");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Comments_user_id_key" ON "Tweet_Comments"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "User_Followers_user_id_key" ON "User_Followers"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "User_Followers_follower_id_key" ON "User_Followers"("follower_id");
