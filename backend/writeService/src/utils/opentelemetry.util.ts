import opentelemetry from '@opentelemetry/api';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { GrpcInstrumentation } from '@opentelemetry/instrumentation-grpc';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { ExpressInstrumentation } from '@opentelemetry/instrumentation-express';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';

const collectorOptions = {
    url: 'http://jaeger:4318/v1/traces',
    headers: {},
    concurrencyLimit: 10,
};

export function initTracer() {
    const provider = new NodeTracerProvider({
        resource: new Resource({
            [SemanticResourceAttributes.SERVICE_NAME]: 'writeService',
        }),
    });

    let exporter = new OTLPTraceExporter(collectorOptions);

    provider.addSpanProcessor(
        new BatchSpanProcessor(exporter, {
            maxQueueSize: 100,
            maxExportBatchSize: 10,
            scheduledDelayMillis: 500,
            exportTimeoutMillis: 30000,
        })
    );

    provider.register();

    registerInstrumentations({
        instrumentations: [
            new GrpcInstrumentation(),
            new HttpInstrumentation(),
            new ExpressInstrumentation(),
        ],
    });

    return opentelemetry.trace.getTracer('writeService');
}
