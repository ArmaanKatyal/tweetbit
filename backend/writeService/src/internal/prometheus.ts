import client from 'prom-client';

const collectDefaultMetrics = client.collectDefaultMetrics;
const prefix = 'write_service_';

collectDefaultMetrics({ prefix });

export const register = client.register;

export enum MetricsCode {
    Ok = '200',
    BadRequest = '400',
    InternalServerError = '500',
}

export enum MetricsMethod {
    Get = 'GET',
    Post = 'POST',
    Success = 'SUCCESS',
    Error = 'ERROR',
}

export const httpTransactionTotal = new client.Counter({
    name: `${prefix}http_transaction_total`,
    help: 'Total number of HTTP transactions',
    labelNames: ['code', 'method'],
});

export const httpResponseTimeHistogram = new client.Histogram({
    name: `${prefix}http_response_time_histogram`,
    help: 'HTTP response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code', 'method'],
});

export const createTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}create_tweet_response_time_histogram`,
    help: 'Create tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const deleteTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}delete_tweet_response_time_histogram`,
    help: 'Delete tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const likeTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}like_tweet_response_time_histogram`,
    help: 'Like tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const unlikeTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}unlike_tweet_response_time_histogram`,
    help: 'Unlike tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const retweetTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}retweet_tweet_response_time_histogram`,
    help: 'Retweet tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const commentTweetResponseTimeHistogram = new client.Histogram({
    name: `${prefix}comment_tweet_response_time_histogram`,
    help: 'Comment tweet response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const followUserResponseTimeHistogram = new client.Histogram({
    name: `${prefix}follow_user_response_time_histogram`,
    help: 'Follow user response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const unfollowUserResponseTimeHistogram = new client.Histogram({
    name: `${prefix}unfollow_user_response_time_histogram`,
    help: 'Unfollow user response time in milliseconds',
    buckets: [0, 5, 20],
    labelNames: ['code'],
});

export const IncHttpTransaction = (code: MetricsCode, method: MetricsMethod) => {
    httpTransactionTotal.labels(code, method).inc();
};

export const ObserveHttpResponseTime = (code: MetricsCode, method: MetricsMethod, time: number) => {
    httpResponseTimeHistogram.labels(code, method).observe(time);
};

export const collectMetrics = (
    functionName: string,
    code: MetricsCode,
    method: MetricsMethod,
    time: number
) => {
    IncHttpTransaction(code, method);
    ObserveHttpResponseTime(code, method, Date.now() - time);
    ObserveMethodResponseTime(functionName, code, time);
};

const ObserveMethodResponseTime = (method: string, code: string, time: number) => {
    switch (method) {
        case 'createTweet': {
            createTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'deleteTweet': {
            deleteTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'likeTweet': {
            likeTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'unlikeTweet': {
            unlikeTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'retweetTweet': {
            retweetTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'commentTweet': {
            commentTweetResponseTimeHistogram.labels(code).observe(time);
        }
        case 'followUser': {
            followUserResponseTimeHistogram.labels(code).observe(time);
        }
        case 'unfollowUser': {
            unfollowUserResponseTimeHistogram.labels(code).observe(time);
        }
    }
};
