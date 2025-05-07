import axios from 'axios';

const api = axios.create({
  baseURL: 'https://urbannest-backend.onrender.com',
});

export default api;