-- AddForeignKey
ALTER TABLE "Tweet_Likes" ADD CONSTRAINT "Tweet_Likes_tweet_id_fkey" FOREIGN KEY ("tweet_id") REFERENCES "Tweet"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Tweet_Comments" ADD CONSTRAINT "Tweet_Comments_tweet_id_fkey" FOREIGN KEY ("tweet_id") REFERENCES "Tweet"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
