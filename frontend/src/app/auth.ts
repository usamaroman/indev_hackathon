import axios from 'axios';

interface AuthResponse {
  access_token: string;
}

export const login = async (): Promise<AuthResponse> => {
  try {
    const response = await axios.post<AuthResponse>('http://localhost:8080/auth/login', {
      login: 'admin',
      password: 'admin'
    });
    return response.data;
  } catch (error) {
    throw new Error('Login failed. Please check your credentials.');
  }
};
