import pino from 'pino';
import pinoPretty from 'pino-pretty';
import fs from 'fs';

const levels = {
    http: 10,
    debug: 20,
    info: 30,
    warn: 40,
    error: 50,
    fatal: 60,
};

/**
 * Check if log file exists, if not create it
 * @param destination string
 * @returns string
 */
const checkLogExist = (destination: string): string => {
    if (fs.existsSync(destination)) {
        return destination;
    } else {
        fs.closeSync(fs.openSync(destination, 'a'));
        return destination;
    }
};

const logger = pino(
    {
        prettifier: pinoPretty({
            colorize: true,
            levelFirst: true,
            translateTime: 'yyyy-mm-dd h:MM:ss TT',
        }),
        customLevels: levels,
        useOnlyCustomLevels: true,
        level: 'http',
    },
    pino.destination(checkLogExist('combined.log'))
);

export default logger;
