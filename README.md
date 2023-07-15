# TweetBit

A microservices based clone of Twitter.

## Use cases and constraints

### Use cases
- User posts a tweet
    - Service pushes tweets to followers, sending push notifications
- User views the user timeline (activity from the user)
- User views the home timeline (activity from people the user is following)
- User views the home timeline (activity from people the user is following)
- User searches a keyword
- Service has high availability

### Constraints and assumptions

#### State Assumptions

#### General
- Traffic is not evenly distributed
- Posting a tweet should be fast
    - Fanning out a tweet to all of your folowers should be fast, unless you have _millions_ of followers
- Each tweet averages a fanout of 10 deliveries

#### Timeline
- Viewing the timeline should be fast
- Tweetbit is more read heavy than write heavy
    - Optimize for fast reads of tweets
- Timeline is not ordered   
    - Tweets are not ordered by time, but by relevance
    - Tweets are ordered by relevance, which is determined by the number of likes, retweets, and comments

#### Search
- Searching should be fast and accurate
- Search is read heavy

## High level design
![High level design](./frontend/highleveldesign.png)

## Detailed design

### User posts a tweet
