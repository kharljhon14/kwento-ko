import type { User } from '@/types/user';
import axios, { type AxiosResponse } from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';
axios.defaults.withCredentials = true;

axios.interceptors.request.use((config) => {
  return config;
});

axios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error.response);
  }
);

type Envelope<T> = {
  data: T;
};

const responseBody = <T>(res: AxiosResponse<Envelope<T>>) => res.data;

const requests = {
  get: <T>(url: string) => axios.get(url).then(responseBody<T>),
  post: <T, D>(url: string, body: D) => axios.post(url, body).then(responseBody<T>),
  patch: <T, D>(url: string, body: D) => axios.patch(url, body).then(responseBody<T>),
  delete: <T>(url: string) => axios.delete(url).then(responseBody<T>)
};

const user = {
  getUser: () => requests.get<User>('/api/v1/users')
};

const agent = {
  user
};

export default agent;
