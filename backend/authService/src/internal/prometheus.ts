import client from 'prom-client';

const collectDefaultMetrics = client.collectDefaultMetrics;
const prefix = 'auth_service_';

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

export const IncHttpTransaction = (code: MetricsCode, method: MetricsMethod) => {
    httpTransactionTotal.labels(code, method).inc();
};

export const ObserveHttpResponseTime = (code: MetricsCode, method: MetricsMethod, time: number) => {
    httpResponseTimeHistogram.labels(code, method).observe(time);
};
