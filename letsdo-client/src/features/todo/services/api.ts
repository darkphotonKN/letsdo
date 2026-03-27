import { apiClient } from "@/lib/api/client";
import { API_ENDPOINTS } from "@/lib/api/endpoints";
import type { Todo, CreateTodoRequest } from "../types";

export const todoService = {
  create: async (payload: CreateTodoRequest): Promise<Todo> => {
    const { data } = await apiClient.post(API_ENDPOINTS.TODO.CREATE, payload);
    return data;
  },
};
