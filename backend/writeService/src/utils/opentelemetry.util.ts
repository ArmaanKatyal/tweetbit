import opentelemetry from '@opentelemetry/api';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { JaegerExporter } from '@opentelemetry/exporter-jaeger';
import { GrpcInstrumentation } from '@opentelemetry/instrumentation-grpc';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';

export function initTracer() {
    const provider = new NodeTracerProvider({
        resource: new Resource({
            [SemanticResourceAttributes.SERVICE_NAME]: 'writeService',
        }),
    });

    let exporter = new JaegerExporter();

    provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

    provider.register();

    registerInstrumentations({
        instrumentations: [new GrpcInstrumentation(), new HttpInstrumentation()],
    });

    return opentelemetry.trace.getTracer('writeService');
}
