import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_NODE_ENV === 'production'
      ? 'https://urbannest-backend.onrender.com'
      : '/api',
});

export default api;