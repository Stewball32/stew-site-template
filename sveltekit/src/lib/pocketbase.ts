import PocketBase from 'pocketbase';
import { PUBLIC_PB_PORT } from '$env/static/public';
import { dev } from '$app/environment';

const pb = new PocketBase(dev ? `http://localhost:${PUBLIC_PB_PORT}` : undefined);

export default pb;
