export interface Todo {
  id: string;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTodoRequest {
  name: string;
  description?: string;
}
