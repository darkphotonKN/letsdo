export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: "/api/auth/login",
    REGISTER: "/api/auth/register",
    REFRESH: "/api/auth/refresh",
    ME: "/api/auth/me",
  },

  TODO: {
    CREATE: "/api/todos",
  },
} as const;
