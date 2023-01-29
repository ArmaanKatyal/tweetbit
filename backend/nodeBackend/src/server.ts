import {app} from './app';
import nodeConfig from 'config';

const port = nodeConfig.get('server.port') || 8080;

app.listen(port, () => {
    console.log(`Server is listening on port ${port}`);
}) // start the server
