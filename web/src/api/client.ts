import axios from 'axios';

const client = axios.create({
  baseURL: '/api',
  timeout: 300000  // 300 秒（5分钟），Graph 等复杂任务可能需要更长时间
});

export default client;
