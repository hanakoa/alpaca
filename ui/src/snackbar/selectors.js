import { NAME } from './constants';

const messageSelector = state => state[NAME].message;
const openSelector = state => state[NAME].open;

export { messageSelector, openSelector };
