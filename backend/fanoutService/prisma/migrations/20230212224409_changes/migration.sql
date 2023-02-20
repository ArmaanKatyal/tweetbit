/*
  Warnings:

  - A unique constraint covering the columns `[user_id]` on the table `Tweet` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `Tweet_Comments` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id]` on the table `Tweet_Likes` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Tweet_user_id_key" ON "Tweet"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Comments_user_id_key" ON "Tweet_Comments"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "Tweet_Likes_user_id_key" ON "Tweet_Likes"("user_id");
