/*
  Warnings:

  - A unique constraint covering the columns `[tweet_id]` on the table `Tweet_Likes` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `Tweet_Likes` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Likes_tweet_id_key" ON "Tweet_Likes"("tweet_id");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Likes_user_id_key" ON "Tweet_Likes"("user_id");
