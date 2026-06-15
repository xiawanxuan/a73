import axios from 'axios'
import type { ApiResponse } from '../types'

const http = axios.create({
  baseURL: '/api',
  timeout: 30000
})

http.interceptors.request.use(
  (config) => {
    const userId = localStorage.getItem('user_id')
    if (userId) {
      config.headers['X-User-ID'] = userId
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

http.interceptors.response.use(
  (response) => {
    const data: ApiResponse<any> = response.data
    if (data.code === 0) {
      return data.data
    }
    return Promise.reject(new Error(data.msg || '请求失败'))
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default http
